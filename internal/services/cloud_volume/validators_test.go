package cloud_volume

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/G-Core/terraform-provider-gcore/internal/importadoption/adoptiontest"
)

// TestCheckCreatePrerequisites pins the plan-time mirror of the API's create
// rules, driven through the resource's ModifyPlan so the sealed delegation in
// resource.go is exercised, not just the helper. The rules exist because
// removal of a create-only attribute plans a replacement and Terraform
// destroys before it creates: without them a config that drops image_id from
// a source="image" volume would destroy the volume and then 400 on create.
// They are state-aware: a 2.0.0-rc.1 state restated by the config (no
// replacement planned) warns instead of erroring, so the upgrade keeps
// planning no changes.
func TestCheckCreatePrerequisites(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	s := ResourceSchema(ctx)
	objType := s.Type().TerraformType(ctx).(tftypes.Object)
	unknown := tftypes.UnknownValue
	raw := func(set map[string]any) tftypes.Value {
		vals := map[string]tftypes.Value{}
		for name, ty := range objType.AttributeTypes {
			vals[name] = tftypes.NewValue(ty, nil)
		}
		for name, v := range set {
			vals[name] = tftypes.NewValue(objType.AttributeTypes[name], v)
		}
		return tftypes.NewValue(objType, vals)
	}
	nullRaw := tftypes.NewValue(objType, nil)

	type shape = map[string]any
	tests := []struct {
		name     string
		plan     shape
		config   shape  // nil = same as plan
		state    shape  // nil = creating
		wantErr  string // attribute the error is attached to; "" = no error
		wantWarn string // attribute the warning is attached to; "" = no warning
	}{
		// Create.
		{name: "create: image without image_id", plan: shape{"source": "image"}, wantErr: "image_id"},
		{name: "create: IMAGE case-insensitive without image_id", plan: shape{"source": "Image"}, wantErr: "image_id"},
		{name: "create: image with image_id", plan: shape{"source": "image", "image_id": "img"}},
		{name: "create: snapshot without snapshot_id", plan: shape{"source": "snapshot"}, wantErr: "snapshot_id"},
		{name: "create: snapshot with snapshot_id", plan: shape{"source": "snapshot", "snapshot_id": "snap"}},
		{name: "create: attachment_tag without instance", plan: shape{"source": "new-volume", "attachment_tag": "data"}, wantErr: "instance_id_to_attach_to"},
		{name: "create: attachment_tag with instance", plan: shape{"source": "new-volume", "attachment_tag": "data", "instance_id_to_attach_to": "srv"}},
		{name: "create: new-volume with image_id is tolerated", plan: shape{"source": "new-volume", "image_id": "img"}},
		{name: "create: unknown source is left to the API", plan: shape{"source": unknown}},
		{name: "create: unknown image_id is left to the API", plan: shape{"source": "image", "image_id": unknown}},
		{name: "create: unknown attachment_tag is left to the API", plan: shape{"source": "new-volume", "attachment_tag": unknown}},
		{name: "create: null source", plan: shape{}},
		// rc.1-poisoned state restated by the config: warn, never error.
		{
			name:     "poisoned image without image_id, unchanged -> warning",
			plan:     shape{"source": "image", "name": "v"},
			state:    shape{"source": "image", "name": "v"},
			wantWarn: "image_id",
		},
		{
			name:     "poisoned attachment_tag without instance, unchanged -> warning",
			plan:     shape{"source": "new-volume", "attachment_tag": "data", "name": "v"},
			state:    shape{"source": "new-volume", "attachment_tag": "data", "name": "v"},
			wantWarn: "instance_id_to_attach_to",
		},
		{
			name:     "poisoned shape with a name/tags-only change -> warning, in-place",
			plan:     shape{"source": "image", "name": "renamed"},
			state:    shape{"source": "image", "name": "v"},
			wantWarn: "image_id",
		},
		{
			// type_name is Computed: the framework plans it unknown when it is
			// unconfigured and prior state is null; RequiresReplaceIfConfigured
			// ignores an unconfigured attribute.
			name:     "poisoned shape with an unknown UNconfigured type_name (refresh=false) is not a replacement",
			plan:     shape{"source": "image", "name": "renamed", "type_name": unknown},
			config:   shape{"source": "image", "name": "renamed"},
			state:    shape{"source": "image", "name": "v"},
			wantWarn: "image_id",
		},
		{
			// ImportState writes project_id/region_id from the import id; a
			// configuration that omits them plans null for both, and
			// RequiresReplaceIfConfigured does not replace on an unconfigured
			// attribute. This is the unchanged rc.1 upgrade shape.
			name:     "poisoned shape with env-default project_id/region_id going null after import -> warning",
			plan:     shape{"source": "image", "name": "v"},
			state:    shape{"source": "image", "name": "v", "project_id": int64(379987), "region_id": int64(76)},
			wantWarn: "image_id",
		},
		{
			name:     "poisoned shape with project_id dropped from the configuration -> warning",
			plan:     shape{"source": "image", "name": "renamed"},
			state:    shape{"source": "image", "name": "v", "project_id": int64(379987)},
			wantWarn: "image_id",
		},
		{
			name:     "poisoned shape with type_name configured in another case -> warning",
			plan:     shape{"source": "image", "name": "v", "type_name": "standard"}, // the case-fold modifier planned the prior value
			config:   shape{"source": "image", "name": "v", "type_name": "Standard"},
			state:    shape{"source": "image", "name": "v", "type_name": "standard"},
			wantWarn: "image_id",
		},
		{
			// ListRequiresReplaceIfNotNull adopts a value over a null prior in place.
			name:     "poisoned shape with lifecycle_policy_ids set over a null record -> warning",
			plan:     shape{"source": "image", "name": "v", "lifecycle_policy_ids": []tftypes.Value{tftypes.NewValue(tftypes.Number, 1)}},
			state:    shape{"source": "image", "name": "v"},
			wantWarn: "image_id",
		},
		// A replacement is coming: error before the destroy.
		{
			name:    "poisoned shape plus a type_name change -> error",
			plan:    shape{"source": "image", "type_name": "ssd_hiiops"},
			state:   shape{"source": "image", "type_name": "standard"},
			wantErr: "image_id",
		},
		{
			name:    "poisoned shape plus a project_id change -> error",
			plan:    shape{"source": "image", "project_id": int64(1)},
			state:   shape{"source": "image", "project_id": int64(379987)},
			wantErr: "image_id",
		},
		{
			name:    "poisoned shape plus a lifecycle_policy_ids change -> error",
			plan:    shape{"source": "image", "lifecycle_policy_ids": []tftypes.Value{tftypes.NewValue(tftypes.Number, 2)}},
			state:   shape{"source": "image", "lifecycle_policy_ids": []tftypes.Value{tftypes.NewValue(tftypes.Number, 1)}},
			wantErr: "image_id",
		},
		{
			// A configured unknown (a same-apply resource id) on a create-only
			// attribute replaces without a marker; Terraform would destroy the
			// volume and the create would 400.
			name:    "poisoned shape plus a configured unknown instance_id_to_attach_to -> error",
			plan:    shape{"source": "image", "name": "v", "instance_id_to_attach_to": unknown},
			state:   shape{"source": "image", "name": "v"},
			wantErr: "image_id",
		},
		{
			name:    "poisoned shape plus a configured unknown snapshot_id -> error",
			plan:    shape{"source": "image", "name": "v", "snapshot_id": unknown},
			state:   shape{"source": "image", "name": "v"},
			wantErr: "image_id",
		},
		{
			name:    "poisoned shape plus a configured unknown type_name over a known record -> error",
			plan:    shape{"source": "image", "name": "v", "type_name": unknown},
			state:   shape{"source": "image", "name": "v", "type_name": "standard"},
			wantErr: "image_id",
		},
		{
			name:    "image_id removed from an image volume -> error",
			plan:    shape{"source": "image"},
			state:   shape{"source": "image", "image_id": "img"},
			wantErr: "image_id",
		},
		{
			name:    "source changed to image without image_id -> error",
			plan:    shape{"source": "image"},
			state:   shape{"source": "new-volume"},
			wantErr: "image_id",
		},
		{
			name:    "adopting source=image after import without image_id -> error",
			plan:    shape{"source": "image"},
			state:   shape{},
			wantErr: "image_id",
		},
		{
			name:    "attachment_tag set on day two without instance -> error",
			plan:    shape{"source": "new-volume", "attachment_tag": "data"},
			state:   shape{"source": "new-volume"},
			wantErr: "instance_id_to_attach_to",
		},
		{
			name:  "valid update is silent",
			plan:  shape{"source": "image", "image_id": "img", "name": "renamed"},
			state: shape{"source": "image", "image_id": "img", "name": "v"},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			stateRaw := nullRaw
			if tc.state != nil {
				stateRaw = raw(tc.state)
			}
			config := tc.plan
			if tc.config != nil {
				config = tc.config
			}
			req := resource.ModifyPlanRequest{
				Config: tfsdk.Config{Raw: raw(config), Schema: s},
				Plan:   tfsdk.Plan{Raw: raw(tc.plan), Schema: s},
				State:  tfsdk.State{Raw: stateRaw, Schema: s},
			}
			resp := &resource.ModifyPlanResponse{Plan: req.Plan}
			(&CloudVolumeResource{}).ModifyPlan(ctx, req, resp)

			var errs, warns []string
			for _, d := range resp.Diagnostics.Errors() {
				errs = append(errs, d.Summary()+": "+d.Detail())
			}
			for _, d := range resp.Diagnostics.Warnings() {
				if strings.Contains(d.Summary(), "cannot be re-created") {
					warns = append(warns, d.Summary()+": "+d.Detail())
				}
			}
			if tc.wantErr == "" && len(errs) > 0 {
				t.Fatalf("unexpected errors:\n%s", strings.Join(errs, "\n"))
			}
			if tc.wantErr != "" {
				if len(errs) != 1 {
					t.Fatalf("want exactly one error on %q, got %d:\n%s", tc.wantErr, len(errs), strings.Join(errs, "\n"))
				}
				if !strings.Contains(errs[0], "Missing required argument") || !strings.Contains(errs[0], `"`+tc.wantErr+`"`) {
					t.Errorf("error does not name %q as missing:\n%s", tc.wantErr, errs[0])
				}
			}
			if tc.wantWarn == "" && len(warns) > 0 {
				t.Fatalf("unexpected warnings:\n%s", strings.Join(warns, "\n"))
			}
			if tc.wantWarn != "" {
				if len(warns) != 1 {
					t.Fatalf("want exactly one warning on %q, got %d:\n%s", tc.wantWarn, len(warns), strings.Join(warns, "\n"))
				}
				for _, want := range []string{`"` + tc.wantWarn + `"`, "-replace=", "remove the resource from state"} {
					if !strings.Contains(warns[0], want) {
						t.Errorf("warning does not contain %q:\n%s", want, warns[0])
					}
				}
			}
		})
	}

	// Destroy plans are silent.
	req := resource.ModifyPlanRequest{
		Plan:  tfsdk.Plan{Raw: nullRaw, Schema: s},
		State: tfsdk.State{Raw: raw(shape{"source": "image"}), Schema: s},
	}
	resp := &resource.ModifyPlanResponse{Plan: req.Plan}
	(&CloudVolumeResource{}).ModifyPlan(ctx, req, resp)
	if len(resp.Diagnostics) != 0 {
		t.Errorf("destroy: %v", resp.Diagnostics)
	}
}

// TestReplacementDerivedWhileAdopting covers the row ModifyPlan cannot be driven
// through in a unit test - the import-adoption marker lives in a framework-
// internal private-state type - by calling the derivation with the marker's
// presence as a value: a null prior value adopted under a live marker is not a
// replacement (the adopting family's modifier skips it), a null prior value
// without the marker is.
func TestReplacementDerivedWhileAdopting(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	s := ResourceSchema(ctx)
	objType := s.Type().TerraformType(ctx).(tftypes.Object)
	raw := func(set map[string]any) tftypes.Value {
		vals := map[string]tftypes.Value{}
		for name, ty := range objType.AttributeTypes {
			vals[name] = tftypes.NewValue(ty, nil)
		}
		for name, v := range set {
			vals[name] = tftypes.NewValue(objType.AttributeTypes[name], v)
		}
		return tftypes.NewValue(objType, vals)
	}
	type shape = map[string]any
	derived := func(config, plan, state shape, adopting bool) bool {
		return replacementDerived(ctx,
			tfsdk.Config{Raw: raw(config), Schema: s},
			tfsdk.Plan{Raw: raw(plan), Schema: s},
			tfsdk.State{Raw: raw(state), Schema: s}, adopting)
	}

	imported := shape{"name": "v", "project_id": int64(379987), "region_id": int64(76)}
	adoptingPlan := shape{"source": "image", "image_id": "img", "instance_id_to_attach_to": "srv", "name": "v"}
	if derived(adoptingPlan, adoptingPlan, imported, true) {
		t.Error("adopting create-only attributes over a null import record under a live marker was derived as a replacement")
	}
	if !derived(adoptingPlan, adoptingPlan, imported, false) {
		t.Error("setting create-only attributes over a null record WITHOUT a marker was not derived as a replacement")
	}
	// A configured unknown over a null record is adopted too (the modifier
	// only looks at prior null + marker), a known-to-known change is not.
	unknownPlan := shape{"source": "image", "image_id": tftypes.UnknownValue, "name": "v"}
	if derived(unknownPlan, unknownPlan, imported, true) {
		t.Error("a configured unknown image_id over a null record under a live marker was derived as a replacement")
	}
	if !derived(unknownPlan, unknownPlan, imported, false) {
		t.Error("a configured unknown image_id over a null record without a marker was not derived as a replacement")
	}
	changed := shape{"source": "image", "image_id": "other", "name": "v"}
	if !derived(changed, changed, shape{"source": "image", "image_id": "img", "name": "v"}, true) {
		t.Error("a known image_id change under a live marker was not derived as a replacement")
	}
}

// TestReplaceBearingMatchesSchema pins replaceBearing against the schema: every
// attribute carrying a replace-capable modifier is listed with the rule that
// modifier applies, and nothing else is listed.
func TestReplaceBearingMatchesSchema(t *testing.T) {
	t.Parallel()

	ruleByModifier := map[string]replaceRule{
		"planmodifiers.stringRequiresReplaceUnlessAdoptingModifier": replaceUnlessAdopting,
		"int64planmodifier.requiresReplaceIfModifier":               replaceIfConfigured,
		"stringplanmodifier.requiresReplaceIfModifier":              replaceIfConfigured,
		"planmodifiers.listRequiresReplaceIfNotNullModifier":        replaceIfStateNotNull,
	}
	attrs := ResourceSchema(context.Background()).Attributes
	want := map[string]replaceRule{}
	for name, a := range attrs {
		switch a.(type) {
		case schema.StringAttribute, schema.Int64Attribute, schema.ListAttribute:
		default:
			continue // nested attributes carry no replace modifier here
		}
		mods := adoptiontest.ModifierTypeNames(t, attrs, name)
		caseFolded := false
		for _, m := range mods {
			if m == "planmodifiers.stringUseStateForCaseInsensitiveEqual" {
				caseFolded = true
			}
		}
		for _, m := range mods {
			rule, ok := ruleByModifier[m]
			if !ok {
				continue
			}
			if rule == replaceIfConfigured && caseFolded {
				rule = replaceIfConfiguredCaseInsensitive
			}
			want[name] = rule
		}
	}
	for name, rule := range want {
		if got, ok := replaceBearing[name]; !ok || got != rule {
			t.Errorf("replaceBearing[%q] = %v (listed=%v), schema modifiers imply %v", name, got, ok, rule)
		}
	}
	for name := range replaceBearing {
		if _, ok := want[name]; !ok {
			t.Errorf("replaceBearing lists %q but no schema modifier on it plans a replacement", name)
		}
	}
}

// TestConfigValidatorsEmpty pins that the create rules no longer run as
// config validators: a ConfigValidator cannot see prior state, so it would
// reject an unchanged 2.0.0-rc.1 configuration outright.
func TestConfigValidatorsEmpty(t *testing.T) {
	t.Parallel()

	if got := (&CloudVolumeResource{}).ConfigValidators(context.Background()); len(got) != 0 {
		t.Errorf("ConfigValidators = %v, want none (the rules are state-aware and live in ModifyPlan)", got)
	}
}
