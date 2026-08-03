package planmodifiers_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/planmodifiers"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var hmAttrTypes = map[string]attr.Type{"delay": types.Int64Type}

func runConfigNullModifier(req planmodifier.ObjectRequest) *planmodifier.ObjectResponse {
	resp := &planmodifier.ObjectResponse{PlanValue: req.PlanValue}
	planmodifiers.ObjectUseStateForUnknownWhenConfigNull().PlanModifyObject(context.Background(), req, resp)
	return resp
}

// Create: block omitted, no prior state -> plan resolves to null.
func TestObjectUseStateForUnknownWhenConfigNull_CreateNulls(t *testing.T) {
	t.Parallel()
	resp := runConfigNullModifier(planmodifier.ObjectRequest{
		PlanValue:   types.ObjectUnknown(hmAttrTypes),
		ConfigValue: types.ObjectNull(hmAttrTypes),
		StateValue:  types.ObjectNull(hmAttrTypes),
	})
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
	}
	if !resp.PlanValue.IsNull() {
		t.Fatalf("expected null plan value, got %s", resp.PlanValue)
	}
}

// Update: block omitted, prior state has a value -> state value preserved.
func TestObjectUseStateForUnknownWhenConfigNull_PreservesState(t *testing.T) {
	t.Parallel()
	stateValue := types.ObjectValueMust(hmAttrTypes, map[string]attr.Value{
		"delay": types.Int64Value(10),
	})
	resp := runConfigNullModifier(planmodifier.ObjectRequest{
		PlanValue:   types.ObjectUnknown(hmAttrTypes),
		ConfigValue: types.ObjectNull(hmAttrTypes),
		StateValue:  stateValue,
	})
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
	}
	if !resp.PlanValue.Equal(stateValue) {
		t.Fatalf("expected state value preserved, got %s", resp.PlanValue)
	}
}

// Unknown config expression is not "omitted" -> plan stays unknown.
func TestObjectUseStateForUnknownWhenConfigNull_ConfigUnknownLeftAlone(t *testing.T) {
	t.Parallel()
	resp := runConfigNullModifier(planmodifier.ObjectRequest{
		PlanValue:   types.ObjectUnknown(hmAttrTypes),
		ConfigValue: types.ObjectUnknown(hmAttrTypes),
		StateValue:  types.ObjectNull(hmAttrTypes),
	})
	if !resp.PlanValue.IsUnknown() {
		t.Fatalf("expected plan value to remain unknown, got %s", resp.PlanValue)
	}
}

// Known plan value (block configured) -> untouched.
func TestObjectUseStateForUnknownWhenConfigNull_KnownPlanUntouched(t *testing.T) {
	t.Parallel()
	planValue := types.ObjectValueMust(hmAttrTypes, map[string]attr.Value{
		"delay": types.Int64Value(5),
	})
	resp := runConfigNullModifier(planmodifier.ObjectRequest{
		PlanValue:   planValue,
		ConfigValue: planValue,
		StateValue:  types.ObjectNull(hmAttrTypes),
	})
	if !resp.PlanValue.Equal(planValue) {
		t.Fatalf("expected plan value unchanged, got %s", resp.PlanValue)
	}
}

// Null plan value -> untouched.
func TestObjectUseStateForUnknownWhenConfigNull_NullPlanUntouched(t *testing.T) {
	t.Parallel()
	resp := runConfigNullModifier(planmodifier.ObjectRequest{
		PlanValue:   types.ObjectNull(hmAttrTypes),
		ConfigValue: types.ObjectNull(hmAttrTypes),
		StateValue:  types.ObjectNull(hmAttrTypes),
	})
	if !resp.PlanValue.IsNull() {
		t.Fatalf("expected plan value to stay null, got %s", resp.PlanValue)
	}
}
