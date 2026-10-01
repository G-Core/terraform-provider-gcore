package cloud_volume

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/G-Core/terraform-provider-gcore/internal/importadoption/adoptiontest"
)

// TestImportAdoptionWiringPinned asserts every sealed hook is literally present
// in resource.go, read off disk. resource.go is generated; these hooks are
// custom hunks, the class of change that silently evaporates on a bad reseal.
func TestImportAdoptionWiringPinned(t *testing.T) {
	t.Parallel()

	adoptiontest.AssertWired(t, "resource.go",
		adoptiontest.StandardWiring("CloudVolumeResource", "importAdoption", "importAdoptionPlan").
			With("Update", "verifyAdoption(ctx,*data,*state)", `importadoption.KeepKeys(dataBytes,"name","tags")`, "backfillTypeName(data)").
			With("Create", "backfillTypeName(data)").
			With("Read", "backfillTypeName(data)").
			With("ImportState", "backfillTypeName(data)").
			With("ModifyPlan", "refuseAdoptionWithoutSource(ctx,req,resp)", "warnRecordedTypeNameDrift(ctx,req,resp)", "checkCreatePrerequisites(ctx,req,resp)"))
	// The guard runs before the warning hook: the hook returns on an existing
	// error, so a refused adoption is one error, not an error plus a warning
	// that says the values are being recorded.
	adoptiontest.AssertOrderInMethod(t, "resource.go", "CloudVolumeResource", "ModifyPlan",
		"refuseAdoptionWithoutSource(ctx,req,resp)", "importAdoptionPlan.ModifyPlan(ctx,req,resp)")
	// BeginUpdate must precede the plan decode so every error raised in Update
	// re-arms the marker (the documented retry guarantee).
	adoptiontest.AssertOrderInMethod(t, "resource.go", "CloudVolumeResource", "Update",
		"importAdoption.BeginUpdate(", "req.Plan.Get(ctx,&data)")
	// verifyAdoption must run before the first mutation (the resize): moved
	// into the PATCH branch, an adopting apply with no name/tag change would
	// record image_id / instance_id_to_attach_to unverified.
	adoptiontest.AssertOrderInMethod(t, "resource.go", "CloudVolumeResource", "Update",
		"verifyAdoption(ctx,*data,*state)", "sizeChanged:=")
	// Advance is one phase per Read; a second call retires the marker in one
	// refresh and the apply after `terraform import` plans a replacement.
	if n := adoptiontest.CountInMethod(t, "resource.go", "CloudVolumeResource", "Read", "importAdoption.Advance(ctx,resp.Private)"); n != 1 {
		t.Errorf("Read calls importAdoption.Advance %d times, want exactly 1", n)
	}
}

// adoptingAttributes are the five no_refresh attributes that carry the
// marker-keyed modifier in its plain (removal replaces) form.
var adoptingAttributes = []string{"source", "attachment_tag", "image_id", "instance_id_to_attach_to", "snapshot_id"}

// TestSchemaModifiersPinned pins the schema half: the five no_refresh
// attributes carry the marker-keyed modifier, and type_name is Computed+Optional
// with UseNonNullStateForUnknown FIRST (a null rc.1 state under -refresh=false
// plus a coincident change must plan unknown, which the apply-time back-fill
// resolves, not null, which it would contradict), the case-folding modifier
// before the framework's RequiresReplaceIfConfigured. Every constructor is
// pinned by its exact callee as well as by %T: the framework's
// RequiresReplace() and RequiresReplaceIfConfigured() share a type, as do
// our plain and ...AllowingRemoval adopting constructors, so %T alone lets a
// reseal swap the rule while keeping the pin green.
func TestSchemaModifiersPinned(t *testing.T) {
	t.Parallel()

	attrs := ResourceSchema(context.Background()).Attributes
	adoptiontest.AssertAdoptingModifiers(t, attrs, adoptingAttributes...)
	// %T cannot see the constructor argument: each call must carry this
	// resource's marker key, or the modifier reads a key nothing writes and
	// every import plans a replacement. Nor can it see the allowRemoval flag:
	// the exact callee (not ...AllowingRemoval) is pinned here and the removal
	// row itself in TestRemovalReplacesThroughChain.
	for _, name := range adoptingAttributes {
		adoptiontest.AssertAttributeCallArgs(t, "schema.go", name, "planmodifiers.StringRequiresReplaceUnlessAdopting", "importAdoption.Key")
	}
	// project_id / region_id: an unconfigured id (null config over the value
	// ImportState recorded) must NOT replace; the framework's RequiresReplace()
	// would, and replaceBearing / checkCreatePrerequisites derive "not a
	// replacement" for that shape, so the volume would be destroyed while the
	// create is refused.
	for _, name := range []string{"project_id", "region_id"} {
		adoptiontest.AssertAttributeCallArgs(t, "schema.go", name, "int64planmodifier.RequiresReplaceIfConfigured", "")
	}
	adoptiontest.AssertComputedOptional(t, attrs, "type_name")
	adoptiontest.AssertModifierTypes(t, attrs, "type_name",
		"stringplanmodifier.useNonNullStateForUnknown", "planmodifiers.stringUseStateForCaseInsensitiveEqual", "stringplanmodifier.requiresReplaceIfModifier")
	for _, fn := range []string{"stringplanmodifier.UseNonNullStateForUnknown", "planmodifiers.StringUseStateForCaseInsensitiveEqual", "stringplanmodifier.RequiresReplaceIfConfigured"} {
		adoptiontest.AssertAttributeCallArgs(t, "schema.go", "type_name", fn, "")
	}
	// lifecycle_policy_ids keeps the null-prior heuristic on purpose (scope-out,
	// see DESIGN.md 0.1); replaceBearing derives replaceIfStateNotNull for it.
	adoptiontest.AssertAttributeCallArgs(t, "schema.go", "lifecycle_policy_ids", "planmodifiers.ListRequiresReplaceIfNotNull", "")
}

// TestRemovalReplacesThroughChain runs removal - config and plan null over a
// known prior value, no marker - through each adopting attribute's real
// PlanModifiers chain and asserts a replacement. %T and the exact-callee pin
// both live in this package's tests; this is the behavioural mirror of cdn's
// handover-shape test, and it is what fails when a seal resolution copies
// cdn's ...AllowingRemoval constructor into the volume: there removal would
// silently drop the record instead of replacing.
func TestRemovalReplacesThroughChain(t *testing.T) {
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
	for _, name := range adoptingAttributes {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			prior := types.StringValue("recorded")
			req := planmodifier.StringRequest{
				Path:        path.Root(name),
				StateValue:  prior,
				ConfigValue: types.StringNull(),
				PlanValue:   types.StringNull(),
				State:       tfsdk.State{Raw: raw(map[string]any{"id": "vol-1", name: "recorded"}), Schema: s},
				Plan:        tfsdk.Plan{Raw: raw(map[string]any{"id": "vol-1"}), Schema: s},
				Config:      tfsdk.Config{Raw: raw(map[string]any{"id": "vol-1"}), Schema: s},
			}
			resp := &planmodifier.StringResponse{PlanValue: req.PlanValue}
			for _, m := range s.Attributes[name].(schema.StringAttribute).PlanModifiers {
				req.PlanValue = resp.PlanValue
				m.PlanModifyString(ctx, req, resp)
			}
			if resp.Diagnostics.HasError() {
				t.Fatalf("unexpected error: %v", resp.Diagnostics)
			}
			if !resp.RequiresReplace {
				t.Errorf("%s: removal over a known prior value without a marker did not plan a replacement (plan = %v)", name, resp.PlanValue)
			}
		})
	}
}

// TestTypeNameChainDecisions runs type_name's real plan modifier chain on the
// shapes the drift warning and the import note describe, so the wording can be
// pinned against the decision rather than against itself.
func TestTypeNameChainDecisions(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	s := ResourceSchema(ctx)
	objType := s.Type().TerraformType(ctx).(tftypes.Object)
	rawWith := func(typeName, volumeType any) tftypes.Value {
		vals := map[string]tftypes.Value{}
		for name, ty := range objType.AttributeTypes {
			vals[name] = tftypes.NewValue(ty, nil)
		}
		vals["type_name"] = tftypes.NewValue(tftypes.String, typeName)
		vals["volume_type"] = tftypes.NewValue(tftypes.String, volumeType)
		return tftypes.NewValue(objType, vals)
	}
	run := func(stateTypeName, volumeType string, config types.String) (types.String, bool) {
		req := planmodifier.StringRequest{
			StateValue: types.StringValue(stateTypeName), ConfigValue: config, PlanValue: config,
			State: tfsdk.State{Raw: rawWith(stateTypeName, volumeType), Schema: s},
			Plan:  tfsdk.Plan{Raw: rawWith(stateTypeName, volumeType), Schema: s},
		}
		resp := &planmodifier.StringResponse{PlanValue: config}
		for _, m := range s.Attributes["type_name"].(schema.StringAttribute).PlanModifiers {
			req.PlanValue = resp.PlanValue
			m.PlanModifyString(ctx, req, resp)
		}
		return resp.PlanValue, resp.RequiresReplace
	}

	// Import with "Standard" against a back-filled "standard": no replacement.
	if plan, replace := run("standard", "standard", types.StringValue("Standard")); replace || !plan.Equal(types.StringValue("standard")) {
		t.Errorf("case-folded match: plan = %v, replace = %v; want standard and no replacement", plan, replace)
	}
	// A different type replaces (no retype in Update).
	if _, replace := run("standard", "standard", types.StringValue("ssd_hiiops")); !replace {
		t.Error("different type_name did not plan a replacement")
	}
	// The rc.1-poisoned record: restating it is a no-op, and setting the TRUE
	// value replaces - exactly what warnRecordedTypeNameDrift must say.
	if _, replace := run("ssd_hiiops", "standard", types.StringValue("ssd_hiiops")); replace {
		t.Error("restating the recorded (wrong) type_name planned a replacement")
	}
	if _, replace := run("ssd_hiiops", "standard", types.StringValue("standard")); !replace {
		t.Error("setting type_name to the reported volume_type over a wrong record did not plan a replacement; " +
			"warnRecordedTypeNameDrift promises it does")
	}
}

// TestCreateOnlyAttributesMatchSchema: every Table key carries the adopting
// modifier and every attribute carrying it is in the Table, so the warning
// names exactly the attributes the modifiers adopt.
func TestCreateOnlyAttributesMatchSchema(t *testing.T) {
	t.Parallel()

	attrs := ResourceSchema(context.Background()).Attributes
	adopting := map[string]bool{}
	for name, a := range attrs {
		sa, ok := a.(schema.StringAttribute)
		if !ok {
			continue
		}
		for _, m := range sa.PlanModifiers {
			if strings.Contains(fmt.Sprintf("%T", m), "RequiresReplaceUnlessAdopting") {
				adopting[name] = true
			}
		}
	}
	for name := range createOnlyAttributes {
		if !adopting[name] {
			t.Errorf("createOnlyAttributes names %q but the schema attribute has no adopting modifier", name)
		}
	}
	for name := range adopting {
		if _, ok := createOnlyAttributes[name]; !ok {
			t.Errorf("schema attribute %q carries the adopting modifier but is missing from createOnlyAttributes", name)
		}
	}
}

// TestAdoptionWithoutSource is the decision table of the ignore_changes guard.
// The marker-live half cannot be driven through ModifyPlan (the framework's
// private-state type is internal), so the pure function carries it.
func TestAdoptionWithoutSource(t *testing.T) {
	t.Parallel()

	null := CloudVolumeModel{}
	tests := []struct {
		name     string
		adopting bool
		creating bool
		plan     CloudVolumeModel
		state    CloudVolumeModel
		want     []string
	}{
		{name: "no marker: day-two set is the modifiers' replacement, not this guard",
			plan: CloudVolumeModel{SnapshotID: types.StringValue("snap")}, state: null},
		{name: "creating: every attribute is applied, nothing adopted", adopting: true, creating: true,
			plan: CloudVolumeModel{SnapshotID: types.StringValue("snap")}, state: null},
		{name: "source planned: the ordinary adopting apply records source too", adopting: true,
			plan: CloudVolumeModel{Source: types.StringValue("snapshot"), SnapshotID: types.StringValue("snap")}, state: null},
		{name: "source unknown: not the ignore_changes shape", adopting: true,
			plan: CloudVolumeModel{Source: types.StringUnknown(), SnapshotID: types.StringValue("snap")}, state: null},
		{name: "source ignored, nothing adopted: the no-op apply is the user's choice", adopting: true,
			plan: null, state: null},
		{name: "source ignored, snapshot_id adopted", adopting: true,
			plan: CloudVolumeModel{SnapshotID: types.StringValue("snap")}, state: null, want: []string{"snapshot_id"}},
		{name: "source ignored, several adopted (sorted)", adopting: true,
			plan:  CloudVolumeModel{SnapshotID: types.StringValue("snap"), AttachmentTag: types.StringValue("data"), ImageID: types.StringValue("img")},
			state: null, want: []string{"attachment_tag", "image_id", "snapshot_id"}},
		{name: "source ignored, value already in state: restated, not adopted", adopting: true,
			plan:  CloudVolumeModel{SnapshotID: types.StringValue("snap")},
			state: CloudVolumeModel{SnapshotID: types.StringValue("snap")}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := adoptionWithoutSource(tc.adopting, tc.creating, tc.plan, tc.state)
			if fmt.Sprint(got) != fmt.Sprint(tc.want) {
				t.Errorf("adoptionWithoutSource = %v, want %v", got, tc.want)
			}
		})
	}

	d := adoptionWithoutSourceError([]string{"snapshot_id"})
	if !d.Equal(diag.NewAttributeErrorDiagnostic(path.Root("source"), d.Summary(), d.Detail())) {
		t.Errorf("the error must be attributed to source (the ignored attribute), got %v", d)
	}
	for _, want := range []string{`"snapshot_id"`, "ignore_changes", "-replace="} {
		if !strings.Contains(d.Detail(), want) {
			t.Errorf("error detail lacks %q:\n%s", want, d.Detail())
		}
	}
}

// TestRefuseAdoptionWithoutSourceNoMarker drives the hook through ModifyPlan
// on the shapes it must leave alone without a marker: the guard never fires
// on a create, a destroy, or a day-two set (which the modifiers replace).
func TestRefuseAdoptionWithoutSourceNoMarker(t *testing.T) {
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
	nullRaw := tftypes.NewValue(objType, nil)
	imported := raw(map[string]any{"id": "vol-1", "size": int64(1)})
	dayTwo := raw(map[string]any{"id": "vol-1", "size": int64(1), "snapshot_id": "snap"})

	for name, req := range map[string]resource.ModifyPlanRequest{
		"create":                 {State: tfsdk.State{Raw: nullRaw, Schema: s}, Plan: tfsdk.Plan{Raw: dayTwo, Schema: s}},
		"destroy":                {State: tfsdk.State{Raw: imported, Schema: s}, Plan: tfsdk.Plan{Raw: nullRaw, Schema: s}},
		"day-two set, no marker": {State: tfsdk.State{Raw: imported, Schema: s}, Plan: tfsdk.Plan{Raw: dayTwo, Schema: s}},
	} {
		resp := &resource.ModifyPlanResponse{Plan: req.Plan}
		refuseAdoptionWithoutSource(ctx, req, resp)
		if len(resp.Diagnostics) != 0 {
			t.Errorf("%s: %v", name, resp.Diagnostics)
		}
	}
}

func TestBackfillTypeName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		typeName types.String
		volType  types.String
		want     types.String
	}{
		{"null + volume_type -> filled", types.StringNull(), types.StringValue("standard"), types.StringValue("standard")},
		{"unknown + volume_type -> filled", types.StringUnknown(), types.StringValue("ssd_hiiops"), types.StringValue("ssd_hiiops")},
		{"unknown + empty -> null", types.StringUnknown(), types.StringValue(""), types.StringNull()},
		{"unknown + null volume_type -> null", types.StringUnknown(), types.StringNull(), types.StringNull()},
		{"null + unknown volume_type -> null", types.StringNull(), types.StringUnknown(), types.StringNull()},
		// The rc.1-poisoned case: a recorded value is never overwritten, so
		// state == config keeps planning no changes (no repair step).
		{"recorded ssd_hiiops next to standard -> UNCHANGED", types.StringValue("ssd_hiiops"), types.StringValue("standard"), types.StringValue("ssd_hiiops")},
		{"equal -> unchanged", types.StringValue("standard"), types.StringValue("standard"), types.StringValue("standard")},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			data := &CloudVolumeModel{TypeName: tc.typeName, VolumeType: tc.volType}
			backfillTypeName(data)
			if !data.TypeName.Equal(tc.want) {
				t.Errorf("type_name = %v, want %v", data.TypeName, tc.want)
			}
		})
	}
	backfillTypeName(nil) // nil-safe
}

func TestVerifyAdoption(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	meta := func(imageID string) customfield.Map[types.String] {
		return customfield.NewMapMust[types.String](ctx, map[string]types.String{"image_id": types.StringValue(imageID)})
	}
	attachments := func(serverIDs ...string) customfield.NestedObjectList[CloudVolumeAttachmentsModel] {
		var atts []CloudVolumeAttachmentsModel
		for _, id := range serverIDs {
			atts = append(atts, CloudVolumeAttachmentsModel{ServerID: types.StringValue(id)})
		}
		return customfield.NewObjectListMust(ctx, atts)
	}

	tests := []struct {
		name      string
		plan      CloudVolumeModel
		state     CloudVolumeModel
		wantError string // substring of the summary; "" = no error
	}{
		{name: "nothing adopted", plan: CloudVolumeModel{}, state: CloudVolumeModel{}},
		{
			name:  "image_id: metadata absent -> adopt",
			plan:  CloudVolumeModel{ImageID: types.StringValue("img-1")},
			state: CloudVolumeModel{VolumeImageMetadata: customfield.NullMap[types.String](ctx)},
		},
		{
			name:  "image_id: metadata equal -> adopt",
			plan:  CloudVolumeModel{ImageID: types.StringValue("img-1")},
			state: CloudVolumeModel{VolumeImageMetadata: meta("img-1")},
		},
		{
			name:      "image_id: metadata different -> error",
			plan:      CloudVolumeModel{ImageID: types.StringValue("img-1")},
			state:     CloudVolumeModel{VolumeImageMetadata: meta("img-2")},
			wantError: "image_id contradicts",
		},
		{
			name:  "image_id already in state is not verified",
			plan:  CloudVolumeModel{ImageID: types.StringValue("img-1")},
			state: CloudVolumeModel{ImageID: types.StringValue("img-1"), VolumeImageMetadata: meta("img-2")},
		},
		{
			name:  "instance: attached -> adopt",
			plan:  CloudVolumeModel{InstanceIDToAttachTo: types.StringValue("srv-1")},
			state: CloudVolumeModel{Attachments: attachments("srv-0", "srv-1")},
		},
		{
			name:      "instance: not attached -> error",
			plan:      CloudVolumeModel{InstanceIDToAttachTo: types.StringValue("srv-1")},
			state:     CloudVolumeModel{Attachments: attachments("srv-0")},
			wantError: "instance_id_to_attach_to contradicts",
		},
		{
			name:      "instance: no attachments -> error",
			plan:      CloudVolumeModel{InstanceIDToAttachTo: types.StringValue("srv-1")},
			state:     CloudVolumeModel{Attachments: customfield.NullObjectList[CloudVolumeAttachmentsModel](ctx)},
			wantError: "instance_id_to_attach_to contradicts",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			diags := verifyAdoption(ctx, tc.plan, tc.state)
			if tc.wantError == "" {
				if diags.HasError() {
					t.Fatalf("unexpected error: %v", diags)
				}
				return
			}
			if !diags.HasError() {
				t.Fatalf("expected an error containing %q", tc.wantError)
			}
			found := false
			for _, d := range diags.Errors() {
				found = found || strings.Contains(d.Summary(), tc.wantError)
			}
			if !found {
				t.Errorf("no error summary contains %q: %v", tc.wantError, diags)
			}
		})
	}
}

// TestVolumeTransport pins the one promise in the warning that could be
// factually wrong: image_id and instance_id_to_attach_to are checked, so the
// text must not call them trusted as-is.
func TestVolumeTransport(t *testing.T) {
	t.Parallel()

	text := volumeTransport(nil)
	if !strings.HasPrefix(text, "Nothing is sent to the API for them.") {
		t.Errorf("transport must state nothing is sent (KeepKeys keeps the PATCH to name+tags): %q", text)
	}
	if strings.Contains(text, `"image_id" is trusted`) || strings.Contains(text, `"instance_id_to_attach_to" is trusted`) {
		t.Errorf("checked attributes described as trusted: %q", text)
	}
	if !strings.Contains(text, `"image_id" is checked`) || !strings.Contains(text, `"instance_id_to_attach_to" against its attachments`) {
		t.Errorf("transport does not name the checks: %q", text)
	}
	d := importAdoptionPlan.Warning.Diagnostic([]string{"source"})
	if !strings.Contains(d.Detail(), text) {
		t.Error("plan hook does not use volumeTransport")
	}
}

func TestWarnRecordedTypeNameDrift(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	s := ResourceSchema(ctx)
	objType := s.Type().TerraformType(ctx).(tftypes.Object)
	stateWith := func(typeName, volumeType any) tfsdk.State {
		vals := map[string]tftypes.Value{}
		for name, ty := range objType.AttributeTypes {
			vals[name] = tftypes.NewValue(ty, nil)
		}
		vals["type_name"] = tftypes.NewValue(tftypes.String, typeName)
		vals["volume_type"] = tftypes.NewValue(tftypes.String, volumeType)
		return tfsdk.State{Raw: tftypes.NewValue(objType, vals), Schema: s}
	}

	tests := []struct {
		name       string
		state      tfsdk.State
		wantDriftW bool
	}{
		{"matching", stateWith("standard", "standard"), false},
		{"case-folded match", stateWith("Standard", "standard"), false},
		{"null type_name", stateWith(nil, "standard"), false},
		{"null volume_type", stateWith("standard", nil), false},
		{"rc.1 lie", stateWith("ssd_hiiops", "standard"), true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			req := resource.ModifyPlanRequest{State: tc.state, Plan: tfsdk.Plan{Raw: tc.state.Raw, Schema: s}}
			resp := &resource.ModifyPlanResponse{Plan: req.Plan}
			warnRecordedTypeNameDrift(ctx, req, resp)
			got := false
			for _, d := range resp.Diagnostics.Warnings() {
				if !strings.Contains(d.Summary(), "Recorded type_name does not match") {
					continue
				}
				got = true
				// The wording must match TestTypeNameChainDecisions: any other
				// value replaces, omission keeps the record, state surgery corrects.
				for _, want := range []string{"plans a replacement", "omitting it keeps the recorded value", "terraform state rm", "terraform import"} {
					if !strings.Contains(d.Detail(), want) {
						t.Errorf("warning detail lacks %q:\n%s", want, d.Detail())
					}
				}
				if strings.Contains(d.Detail(), "no change is applied") {
					t.Errorf("warning promises a no-op repair the modifiers do not deliver:\n%s", d.Detail())
				}
			}
			if got != tc.wantDriftW {
				t.Errorf("drift warning = %v, want %v (%v)", got, tc.wantDriftW, resp.Diagnostics)
			}
			if resp.Diagnostics.HasError() {
				t.Errorf("unexpected error: %v", resp.Diagnostics)
			}
		})
	}

	// Create and destroy plans are silent.
	null := tfsdk.State{Raw: tftypes.NewValue(objType, nil), Schema: s}
	for name, req := range map[string]resource.ModifyPlanRequest{
		"create":  {State: null, Plan: tfsdk.Plan{Raw: stateWith("ssd_hiiops", "standard").Raw, Schema: s}},
		"destroy": {State: stateWith("ssd_hiiops", "standard"), Plan: tfsdk.Plan{Raw: null.Raw, Schema: s}},
	} {
		resp := &resource.ModifyPlanResponse{Plan: req.Plan}
		warnRecordedTypeNameDrift(ctx, req, resp)
		if len(resp.Diagnostics) != 0 {
			t.Errorf("%s: %v", name, resp.Diagnostics)
		}
	}
}
