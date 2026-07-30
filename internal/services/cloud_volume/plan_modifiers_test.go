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
