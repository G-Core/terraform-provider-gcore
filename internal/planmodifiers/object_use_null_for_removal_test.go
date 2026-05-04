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

var objectAttrTypes = map[string]attr.Type{
	"name": types.StringType,
}

func newObjectValue(t *testing.T, name string) types.Object {
	t.Helper()
	v, diags := types.ObjectValue(objectAttrTypes, map[string]attr.Value{
		"name": types.StringValue(name),
	})
	if diags.HasError() {
		t.Fatalf("failed to build object: %s", diags.Errors())
	}
	return v
}

func TestObjectUseNullForRemoval_UnknownConfig_StateHasValue(t *testing.T) {
	t.Parallel()

	state := newObjectValue(t, "existing")
	resp := &planmodifier.ObjectResponse{PlanValue: types.ObjectUnknown(objectAttrTypes)}

	planmodifiers.ObjectUseNullForRemoval().PlanModifyObject(context.Background(), planmodifier.ObjectRequest{
		ConfigValue: types.ObjectUnknown(objectAttrTypes),
		StateValue:  state,
		PlanValue:   types.ObjectUnknown(objectAttrTypes),
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
	}
	if !resp.PlanValue.Equal(state) {
		t.Fatalf("expected state value, got %s", resp.PlanValue)
	}
}

func TestObjectUseNullForRemoval_NullConfig(t *testing.T) {
	t.Parallel()

	state := newObjectValue(t, "existing")
	resp := &planmodifier.ObjectResponse{PlanValue: types.ObjectUnknown(objectAttrTypes)}

	planmodifiers.ObjectUseNullForRemoval().PlanModifyObject(context.Background(), planmodifier.ObjectRequest{
		ConfigValue: types.ObjectNull(objectAttrTypes),
		StateValue:  state,
		PlanValue:   types.ObjectUnknown(objectAttrTypes),
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
	}
	if !resp.PlanValue.IsNull() {
		t.Fatalf("expected null, got %s", resp.PlanValue)
	}
}

func TestObjectUseNullForRemoval_ExplicitValue(t *testing.T) {
	t.Parallel()

	cfg := newObjectValue(t, "new-value")
	state := newObjectValue(t, "old-value")
	resp := &planmodifier.ObjectResponse{PlanValue: cfg}

	planmodifiers.ObjectUseNullForRemoval().PlanModifyObject(context.Background(), planmodifier.ObjectRequest{
		ConfigValue: cfg,
		StateValue:  state,
		PlanValue:   cfg,
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
	}
	if !resp.PlanValue.Equal(cfg) {
		t.Fatalf("expected config value, got %s", resp.PlanValue)
	}
}

func TestObjectUseNullForRemoval_UnknownConfig_NullState(t *testing.T) {
	t.Parallel()

	resp := &planmodifier.ObjectResponse{PlanValue: types.ObjectUnknown(objectAttrTypes)}

	planmodifiers.ObjectUseNullForRemoval().PlanModifyObject(context.Background(), planmodifier.ObjectRequest{
		ConfigValue: types.ObjectUnknown(objectAttrTypes),
		StateValue:  types.ObjectNull(objectAttrTypes),
		PlanValue:   types.ObjectUnknown(objectAttrTypes),
		State:       tfsdk.State{Raw: tftypes.NewValue(tftypes.Object{}, nil)},
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
	}
	if !resp.PlanValue.IsUnknown() {
		t.Fatalf("expected unknown, got %s", resp.PlanValue)
	}
}
