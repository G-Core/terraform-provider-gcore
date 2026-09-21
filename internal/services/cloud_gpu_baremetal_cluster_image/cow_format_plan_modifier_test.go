package cloud_gpu_baremetal_cluster_image

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestCowFormatPlanModifier(t *testing.T) {
	t.Parallel()

	testSchema := schema.Schema{Attributes: map[string]schema.Attribute{
		"cow_format":  schema.BoolAttribute{Computed: true, Optional: true},
		"disk_format": schema.StringAttribute{Computed: true},
	}}
	objectType := tftypes.Object{AttributeTypes: map[string]tftypes.Type{
		"cow_format":  tftypes.Bool,
		"disk_format": tftypes.String,
	}}
	object := func(cow, disk interface{}) tftypes.Value {
		return tftypes.NewValue(objectType, map[string]tftypes.Value{
			"cow_format":  tftypes.NewValue(tftypes.Bool, cow),
			"disk_format": tftypes.NewValue(tftypes.String, disk),
		})
	}
	nullObject := tftypes.NewValue(objectType, nil)
	unknown := tftypes.UnknownValue

	cases := map[string]struct {
		state, plan     tftypes.Value
		config, planned types.Bool
		stateValue      types.Bool
		wantPlan        types.Bool
		wantReplace     bool
	}{
		// creation and destruction are untouched
		"create":  {nullObject, object(unknown, nil), types.BoolNull(), types.BoolUnknown(), types.BoolNull(), types.BoolUnknown(), false},
		"destroy": {object(true, "raw"), nullObject, types.BoolNull(), types.BoolNull(), types.BoolValue(true), types.BoolNull(), false},
		// not configured: plan what disk_format implies
		"omitted, state matches":        {object(true, "raw"), object(true, "raw"), types.BoolNull(), types.BoolValue(true), types.BoolValue(true), types.BoolValue(true), false},
		"omitted, stale false on raw":   {object(false, "raw"), object(false, "raw"), types.BoolNull(), types.BoolValue(false), types.BoolValue(false), types.BoolValue(true), false},
		"omitted, legacy null on raw":   {object(nil, "raw"), object(unknown, "raw"), types.BoolNull(), types.BoolUnknown(), types.BoolNull(), types.BoolValue(true), false},
		"omitted, unmapped, unknown":    {object(true, "vmdk"), object(unknown, "vmdk"), types.BoolNull(), types.BoolUnknown(), types.BoolValue(true), types.BoolValue(true), false},
		"omitted, unmapped, null state": {object(nil, "vmdk"), object(unknown, "vmdk"), types.BoolNull(), types.BoolUnknown(), types.BoolNull(), types.BoolUnknown(), false},
		"unknown config":                {object(true, "raw"), object(unknown, "raw"), types.BoolUnknown(), types.BoolUnknown(), types.BoolValue(true), types.BoolUnknown(), false},
		// configured: compare with what disk_format implies
		"config matches disk_format":      {object(true, "raw"), object(true, "raw"), types.BoolValue(true), types.BoolValue(true), types.BoolValue(true), types.BoolValue(true), false},
		"config differs from disk_format": {object(true, "raw"), object(false, "raw"), types.BoolValue(false), types.BoolValue(false), types.BoolValue(true), types.BoolValue(false), true},
		"legacy null state, raw, true":    {object(nil, "raw"), object(true, "raw"), types.BoolValue(true), types.BoolValue(true), types.BoolNull(), types.BoolValue(true), false},
		"legacy null state, raw, false":   {object(nil, "raw"), object(false, "raw"), types.BoolValue(false), types.BoolValue(false), types.BoolNull(), types.BoolValue(false), true},
		"stale false on raw, true":        {object(false, "raw"), object(true, "raw"), types.BoolValue(true), types.BoolValue(true), types.BoolValue(false), types.BoolValue(true), false},
		"qcow2, config true":              {object(false, "qcow2"), object(true, "qcow2"), types.BoolValue(true), types.BoolValue(true), types.BoolValue(false), types.BoolValue(true), true},
		// unmapped disk_format: RequiresReplaceIfConfigured semantics
		"unmapped, same as state": {object(true, "vmdk"), object(true, "vmdk"), types.BoolValue(true), types.BoolValue(true), types.BoolValue(true), types.BoolValue(true), false},
		"unmapped, differs":       {object(true, "vmdk"), object(false, "vmdk"), types.BoolValue(false), types.BoolValue(false), types.BoolValue(true), types.BoolValue(false), true},
		"unmapped, null state":    {object(nil, "vmdk"), object(true, "vmdk"), types.BoolValue(true), types.BoolValue(true), types.BoolNull(), types.BoolValue(true), true},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := planmodifier.BoolRequest{
				State:       tfsdk.State{Raw: tc.state, Schema: testSchema},
				Plan:        tfsdk.Plan{Raw: tc.plan, Schema: testSchema},
				ConfigValue: tc.config,
				PlanValue:   tc.planned,
				StateValue:  tc.stateValue,
			}
			resp := &planmodifier.BoolResponse{PlanValue: tc.planned}

			cowFormatPlanModifier().PlanModifyBool(context.Background(), req, resp)

			if resp.Diagnostics.HasError() {
				t.Fatalf("unexpected diagnostics: %v", resp.Diagnostics)
			}
			if !resp.PlanValue.Equal(tc.wantPlan) {
				t.Errorf("PlanValue = %v, want %v", resp.PlanValue, tc.wantPlan)
			}
			if resp.RequiresReplace != tc.wantReplace {
				t.Errorf("RequiresReplace = %v, want %v", resp.RequiresReplace, tc.wantReplace)
			}
		})
	}
}
