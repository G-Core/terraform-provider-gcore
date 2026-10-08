package cdn_resource_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/G-Core/terraform-provider-gcore/internal/services/cdn_resource"
)

// TestOriginGroupNamePlan runs origin_group_name's real plan modifier chain.
// The name is server-derived from origin_group, so it may only be carried over
// from state while origin_group stays the same. Pinning it across an
// origin_group switch makes the apply fail with "Provider produced
// inconsistent result" once the API returns the new group's name.
func TestOriginGroupNamePlan(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	s := cdn_resource.ResourceSchema(ctx)
	objType := s.Type().TerraformType(ctx).(tftypes.Object)
	raw := func(originGroup any, originGroupName any) tftypes.Value {
		vals := map[string]tftypes.Value{}
		for name, ty := range objType.AttributeTypes {
			vals[name] = tftypes.NewValue(ty, nil)
		}
		vals["id"] = tftypes.NewValue(tftypes.Number, big.NewFloat(1))
		vals["origin_group"] = tftypes.NewValue(tftypes.Number, originGroup)
		vals["origin_group_name"] = tftypes.NewValue(tftypes.String, originGroupName)
		return tftypes.NewValue(objType, vals)
	}
	run := func(configOriginGroup any) types.String {
		req := planmodifier.StringRequest{
			Path:        path.Root("origin_group_name"),
			StateValue:  types.StringValue("group-a"),
			ConfigValue: types.StringNull(),
			PlanValue:   types.StringUnknown(),
			State:       tfsdk.State{Raw: raw(big.NewFloat(1), "group-a"), Schema: s},
			Plan:        tfsdk.Plan{Raw: raw(configOriginGroup, tftypes.UnknownValue), Schema: s},
			Config:      tfsdk.Config{Raw: raw(configOriginGroup, nil), Schema: s},
		}
		resp := &planmodifier.StringResponse{PlanValue: req.PlanValue}
		for _, m := range s.Attributes["origin_group_name"].(schema.StringAttribute).PlanModifiers {
			req.PlanValue = resp.PlanValue
			m.PlanModifyString(ctx, req, resp)
		}
		if resp.Diagnostics.HasError() {
			t.Fatalf("unexpected error: %v", resp.Diagnostics)
		}
		return resp.PlanValue
	}

	if got := run(big.NewFloat(1)); !got.Equal(types.StringValue("group-a")) {
		t.Errorf("unchanged origin_group: origin_group_name plan = %v, want the state value", got)
	}
	if got := run(big.NewFloat(2)); !got.IsUnknown() {
		t.Errorf("switched origin_group: origin_group_name plan = %v, want unknown", got)
	}
	if got := run(tftypes.UnknownValue); !got.IsUnknown() {
		t.Errorf("origin_group known only after apply: origin_group_name plan = %v, want unknown", got)
	}
}
