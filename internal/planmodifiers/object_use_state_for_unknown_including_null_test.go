package planmodifiers_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/planmodifiers"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestUseStateForUnknownIncludingNullObject_UnknownPlan_NullState(t *testing.T) {
	t.Parallel()

	objType := types.ObjectType{AttrTypes: map[string]attr.Type{"field": types.StringType}}

	resp := &planmodifier.ObjectResponse{
		PlanValue: types.ObjectUnknown(objType.AttrTypes),
	}

	planmodifiers.UseStateForUnknownIncludingNullObject().PlanModifyObject(context.Background(), planmodifier.ObjectRequest{
		PlanValue:  types.ObjectUnknown(objType.AttrTypes),
		StateValue: types.ObjectNull(objType.AttrTypes),
		State:      tfsdk.State{Raw: tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"x": tftypes.String}}, map[string]tftypes.Value{"x": tftypes.NewValue(tftypes.String, nil)})},
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
	}
	if !resp.PlanValue.IsNull() {
		t.Fatalf("expected null (state value), got %s", resp.PlanValue)
	}
}

func TestUseStateForUnknownIncludingNullObject_UnknownPlan_WithState(t *testing.T) {
	t.Parallel()

	objType := types.ObjectType{AttrTypes: map[string]attr.Type{"field": types.StringType}}
	stateValue := types.ObjectValueMust(objType.AttrTypes, map[string]attr.Value{
		"field": types.StringValue("hello"),
	})

	resp := &planmodifier.ObjectResponse{
		PlanValue: types.ObjectUnknown(objType.AttrTypes),
	}

	planmodifiers.UseStateForUnknownIncludingNullObject().PlanModifyObject(context.Background(), planmodifier.ObjectRequest{
		PlanValue:  types.ObjectUnknown(objType.AttrTypes),
		StateValue: stateValue,
		State:      tfsdk.State{Raw: tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"x": tftypes.String}}, map[string]tftypes.Value{"x": tftypes.NewValue(tftypes.String, "hello")})},
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
	}
	if !resp.PlanValue.Equal(stateValue) {
		t.Fatalf("expected state value preserved, got %s", resp.PlanValue)
	}
}

func TestUseStateForUnknownIncludingNullObject_KnownPlan(t *testing.T) {
	t.Parallel()

	objType := types.ObjectType{AttrTypes: map[string]attr.Type{"field": types.StringType}}
	planValue := types.ObjectValueMust(objType.AttrTypes, map[string]attr.Value{
		"field": types.StringValue("new"),
	})

	resp := &planmodifier.ObjectResponse{
		PlanValue: planValue,
	}

	planmodifiers.UseStateForUnknownIncludingNullObject().PlanModifyObject(context.Background(), planmodifier.ObjectRequest{
		PlanValue:  planValue,
		StateValue: types.ObjectNull(objType.AttrTypes),
		State:      tfsdk.State{Raw: tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"x": tftypes.String}}, map[string]tftypes.Value{"x": tftypes.NewValue(tftypes.String, nil)})},
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
	}
	if !resp.PlanValue.Equal(planValue) {
		t.Fatalf("expected plan value unchanged, got %s", resp.PlanValue)
	}
}

func TestUseStateForUnknownIncludingNullObject_NewResource(t *testing.T) {
	t.Parallel()

	objType := types.ObjectType{AttrTypes: map[string]attr.Type{"field": types.StringType}}

	resp := &planmodifier.ObjectResponse{
		PlanValue: types.ObjectUnknown(objType.AttrTypes),
	}

	// New resource: state.Raw is null
	planmodifiers.UseStateForUnknownIncludingNullObject().PlanModifyObject(context.Background(), planmodifier.ObjectRequest{
		PlanValue:  types.ObjectUnknown(objType.AttrTypes),
		StateValue: types.ObjectNull(objType.AttrTypes),
		State:      tfsdk.State{Raw: tftypes.NewValue(tftypes.Object{}, nil)},
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
	}
	// For new resource, should remain unknown (let API compute)
	if !resp.PlanValue.IsUnknown() {
		t.Fatalf("expected unknown for new resource, got %s", resp.PlanValue)
	}
}
