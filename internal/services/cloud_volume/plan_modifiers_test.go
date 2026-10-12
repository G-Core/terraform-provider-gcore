package cloud_volume_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cloud_volume"
)

// sizeModifierSchema is a cut-down stand-in for the volume schema holding only
// the two attributes RequireSizeUnlessDerived reads.
var sizeModifierSchema = schema.Schema{
	Attributes: map[string]schema.Attribute{
		"source": schema.StringAttribute{Required: true},
		"size": schema.Int64Attribute{
			Optional: true,
			Computed: true,
		},
	},
}

var sizeModifierObjectType = tftypes.Object{
	AttributeTypes: map[string]tftypes.Type{
		"source": tftypes.String,
		"size":   tftypes.Number,
	},
}

// sizeModifierValue builds a config/plan/state value for the trimmed schema.
// A nil size means the attribute is absent.
func sizeModifierValue(source string, size *int64) tftypes.Value {
	sizeValue := tftypes.NewValue(tftypes.Number, nil)
	if size != nil {
		sizeValue = tftypes.NewValue(tftypes.Number, *size)
	}
	return tftypes.NewValue(sizeModifierObjectType, map[string]tftypes.Value{
		"source": tftypes.NewValue(tftypes.String, source),
		"size":   sizeValue,
	})
}

func sizeModifierRequest(source string, configSize *int64, stateSize *int64) planmodifier.Int64Request {
	config := tfsdk.Config{Schema: sizeModifierSchema, Raw: sizeModifierValue(source, configSize)}
	plan := tfsdk.Plan{Schema: sizeModifierSchema, Raw: sizeModifierValue(source, configSize)}

	state := tfsdk.State{Schema: sizeModifierSchema, Raw: tftypes.NewValue(sizeModifierObjectType, nil)}
	if stateSize != nil {
		state.Raw = sizeModifierValue(source, stateSize)
	}

	configValue := types.Int64Null()
	if configSize != nil {
		configValue = types.Int64Value(*configSize)
	}

	return planmodifier.Int64Request{
		Path:        path.Root("size"),
		Config:      config,
		Plan:        plan,
		State:       state,
		ConfigValue: configValue,
	}
}

func int64Ptr(v int64) *int64 { return &v }

func TestRequireSizeUnlessDerived(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		source      string
		configSize  *int64
		stateSize   *int64
		expectError bool
	}{
		"create new-volume without size": {
			source:      "new-volume",
			expectError: true,
		},
		"create image volume without size": {
			source:      "image",
			expectError: true,
		},
		"create image volume with an uppercase source": {
			source:      "IMAGE",
			expectError: true,
		},
		"create snapshot volume without size": {
			source: "snapshot",
		},
		"create snapshot volume with a case-insensitive source": {
			source: "Snapshot",
		},
		"create new-volume with size": {
			source:     "new-volume",
			configSize: int64Ptr(10),
		},
		// An existing volume keeps the size in state, so an absent size in the
		// configuration is a no-op, not an error.
		"update new-volume after size was dropped from the config": {
			source:    "new-volume",
			stateSize: int64Ptr(10),
		},
		"update image volume after size was dropped from the config": {
			source:    "image",
			stateSize: int64Ptr(10),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := &planmodifier.Int64Response{}
			cloud_volume.RequireSizeUnlessDerived().PlanModifyInt64(
				context.Background(),
				sizeModifierRequest(test.source, test.configSize, test.stateSize),
				resp,
			)

			if test.expectError && !resp.Diagnostics.HasError() {
				t.Fatal("expected an error, got none")
			}
			if !test.expectError && resp.Diagnostics.HasError() {
				t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
			}
		})
	}
}

// TestRequireSizeUnlessDerived_Destroy checks that a destroy plan, where the
// planned object is null, is left alone.
func TestRequireSizeUnlessDerived_Destroy(t *testing.T) {
	t.Parallel()

	req := sizeModifierRequest("new-volume", nil, int64Ptr(10))
	req.Plan = tfsdk.Plan{Schema: sizeModifierSchema, Raw: tftypes.NewValue(sizeModifierObjectType, nil)}

	resp := &planmodifier.Int64Response{}
	cloud_volume.RequireSizeUnlessDerived().PlanModifyInt64(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
	}
}

// TestRequireSizeUnlessDerived_UnknownSource checks that a source Terraform
// cannot resolve yet is deferred to the API instead of failing the plan.
func TestRequireSizeUnlessDerived_UnknownSource(t *testing.T) {
	t.Parallel()

	unknownSource := tftypes.NewValue(sizeModifierObjectType, map[string]tftypes.Value{
		"source": tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
		"size":   tftypes.NewValue(tftypes.Number, nil),
	})

	req := planmodifier.Int64Request{
		Path:        path.Root("size"),
		Config:      tfsdk.Config{Schema: sizeModifierSchema, Raw: unknownSource},
		Plan:        tfsdk.Plan{Schema: sizeModifierSchema, Raw: unknownSource},
		State:       tfsdk.State{Schema: sizeModifierSchema, Raw: tftypes.NewValue(sizeModifierObjectType, nil)},
		ConfigValue: types.Int64Null(),
	}

	resp := &planmodifier.Int64Response{}
	cloud_volume.RequireSizeUnlessDerived().PlanModifyInt64(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
	}
}

// TestTypeNameRetypeRefusedWithSnapshots: an in-place retype of a volume that
// has snapshots is an error at plan time (the API refuses it at apply); a
// replacement, an unchanged type or an unchanged record is not.
func TestTypeNameRetypeRefusedWithSnapshots(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	s := cloud_volume.ResourceSchema(ctx)
	objType := s.Type().TerraformType(ctx).(tftypes.Object)
	raw := func(typeName, volumeType string, snapshots []string) tftypes.Value {
		vals := map[string]tftypes.Value{}
		for name, ty := range objType.AttributeTypes {
			vals[name] = tftypes.NewValue(ty, nil)
		}
		vals["type_name"] = tftypes.NewValue(tftypes.String, typeName)
		vals["volume_type"] = tftypes.NewValue(tftypes.String, volumeType)
		ids := []tftypes.Value{}
		for _, id := range snapshots {
			ids = append(ids, tftypes.NewValue(tftypes.String, id))
		}
		vals["snapshot_ids"] = tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, ids)
		return tftypes.NewValue(objType, vals)
	}
	// recorded is state type_name, real is volume_type.
	run := func(recorded, real, to string, snapshots []string) *planmodifier.StringResponse {
		req := planmodifier.StringRequest{
			Path:        path.Root("type_name"),
			StateValue:  types.StringValue(recorded),
			ConfigValue: types.StringValue(to),
			PlanValue:   types.StringValue(to),
			State:       tfsdk.State{Raw: raw(recorded, real, snapshots), Schema: s},
			Plan:        tfsdk.Plan{Raw: raw(to, real, snapshots), Schema: s},
		}
		resp := &planmodifier.StringResponse{PlanValue: req.PlanValue}
		cloud_volume.TypeNameReplaceUnlessRetypable().PlanModifyString(ctx, req, resp)
		return resp
	}
	snaps := []string{"snap-1"}

	if resp := run("standard", "standard", "ssd_hiiops", snaps); !resp.Diagnostics.HasError() || resp.RequiresReplace {
		t.Errorf("retype with snapshots: error = %v, replace = %v; want an error and no replacement",
			resp.Diagnostics.HasError(), resp.RequiresReplace)
	}
	if resp := run("standard", "standard", "ssd_hiiops", nil); resp.Diagnostics.HasError() || resp.RequiresReplace {
		t.Errorf("retype without snapshots: diags = %v, replace = %v; want in-place", resp.Diagnostics, resp.RequiresReplace)
	}
	if resp := run("standard", "standard", "cold", snaps); resp.Diagnostics.HasError() || !resp.RequiresReplace {
		t.Errorf("non-retypable change with snapshots: diags = %v, replace = %v; want a replacement", resp.Diagnostics, resp.RequiresReplace)
	}
	if resp := run("standard", "standard", "standard", snaps); resp.Diagnostics.HasError() || resp.RequiresReplace {
		t.Errorf("unchanged type with snapshots: diags = %v, replace = %v; want nothing", resp.Diagnostics, resp.RequiresReplace)
	}
	// rc.1 drift: restating the wrong record is no retype, nor is correcting
	// it to the real type.
	if resp := run("ssd_hiiops", "standard", "ssd_hiiops", snaps); resp.Diagnostics.HasError() || resp.RequiresReplace {
		t.Errorf("unchanged wrong record with snapshots: diags = %v, replace = %v; want nothing", resp.Diagnostics, resp.RequiresReplace)
	}
	if resp := run("ssd_hiiops", "standard", "standard", snaps); resp.Diagnostics.HasError() || resp.RequiresReplace {
		t.Errorf("record corrected to the real type with snapshots: diags = %v, replace = %v; want nothing", resp.Diagnostics, resp.RequiresReplace)
	}
}
