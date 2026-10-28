package cdn_resource_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
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

// TestOptionsKeepStateWhenUnset guards against perpetual drift on options.
// Once the config sets any option, the framework plans every computed object
// or map the config leaves out as unknown, unless a plan modifier carries its
// value over from state. Attributes added to the spec later must get one too.
// Computed scalars are not checked: most carry a static default, and the API
// returns null for the rest when they are unset.
func TestOptionsKeepStateWhenUnset(t *testing.T) {
	t.Parallel()

	s := cdn_resource.ResourceSchema(context.Background())

	var check func(path string, attr schema.Attribute)
	check = func(path string, attr schema.Attribute) {
		switch a := attr.(type) {
		case schema.SingleNestedAttribute:
			if a.Computed && len(a.PlanModifiers) == 0 {
				t.Errorf("%s: computed object has no plan modifier", path)
			}
			for name, child := range a.Attributes {
				check(path+"."+name, child)
			}
		case schema.MapAttribute:
			if a.Computed && len(a.PlanModifiers) == 0 {
				t.Errorf("%s: computed map has no plan modifier", path)
			}
		}
	}
	check("options", s.Attributes["options"])
}

// TestOptionsKeepUnknownConfig runs the plan modifiers of options a config can
// set. An option the config leaves out keeps its state value, but one whose
// config is only known after apply must stay unknown: planning the state value
// instead makes the apply fail with "Provider produced inconsistent final plan".
func TestOptionsKeepUnknownConfig(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	s := cdn_resource.ResourceSchema(ctx)
	objType := s.Type().TerraformType(ctx).(tftypes.Object)
	vals := map[string]tftypes.Value{}
	for name, ty := range objType.AttributeTypes {
		vals[name] = tftypes.NewValue(ty, nil)
	}
	vals["id"] = tftypes.NewValue(tftypes.Number, big.NewFloat(1))
	state := tfsdk.State{Raw: tftypes.NewValue(objType, vals), Schema: s}

	options := s.Attributes["options"].(schema.SingleNestedAttribute)
	fastedge := options.Attributes["fastedge"].(schema.SingleNestedAttribute)
	for name, option := range map[string]schema.SingleNestedAttribute{
		"grpc_passthrough":                        options.Attributes["grpc_passthrough"].(schema.SingleNestedAttribute),
		"network_error_logging":                   options.Attributes["network_error_logging"].(schema.SingleNestedAttribute),
		"fastedge.on_request_headers_after_cache": fastedge.Attributes["on_request_headers_after_cache"].(schema.SingleNestedAttribute),
	} {
		attrTypes := option.GetType().(interface{ AttributeTypes() map[string]attr.Type }).AttributeTypes()
		run := func(config types.Object) types.Object {
			req := planmodifier.ObjectRequest{
				Path:        path.Root("options"),
				State:       state,
				StateValue:  types.ObjectNull(attrTypes),
				ConfigValue: config,
				PlanValue:   types.ObjectUnknown(attrTypes),
			}
			resp := &planmodifier.ObjectResponse{PlanValue: req.PlanValue}
			for _, m := range option.PlanModifiers {
				req.PlanValue = resp.PlanValue
				m.PlanModifyObject(ctx, req, resp)
			}
			return resp.PlanValue
		}

		if got := run(types.ObjectNull(attrTypes)); !got.IsNull() {
			t.Errorf("options.%s left out of the config: plan = %v, want the state value", name, got)
		}
		if got := run(types.ObjectUnknown(attrTypes)); !got.IsUnknown() {
			t.Errorf("options.%s known only after apply: plan = %v, want unknown", name, got)
		}
	}
}
