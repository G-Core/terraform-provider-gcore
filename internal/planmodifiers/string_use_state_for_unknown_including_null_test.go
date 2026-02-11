package planmodifiers_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stainless-sdks/gcore-terraform/internal/planmodifiers"
)

func TestUseStateForUnknownIncludingNullString_UnknownPlan_NullState(t *testing.T) {
	t.Parallel()

	resp := &planmodifier.StringResponse{
		PlanValue: types.StringUnknown(),
	}

	// Existing resource with null state value
	planmodifiers.UseStateForUnknownIncludingNullString().PlanModifyString(context.Background(), planmodifier.StringRequest{
		PlanValue:  types.StringUnknown(),
		StateValue: types.StringNull(),
		State:      tfsdk.State{Raw: tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"x": tftypes.String}}, map[string]tftypes.Value{"x": tftypes.NewValue(tftypes.String, nil)})},
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
	}
	if !resp.PlanValue.IsNull() {
		t.Fatalf("expected null (state value), got %s", resp.PlanValue)
	}
}

func TestUseStateForUnknownIncludingNullString_UnknownPlan_WithState(t *testing.T) {
	t.Parallel()

	resp := &planmodifier.StringResponse{
		PlanValue: types.StringUnknown(),
	}

	planmodifiers.UseStateForUnknownIncludingNullString().PlanModifyString(context.Background(), planmodifier.StringRequest{
		PlanValue:  types.StringUnknown(),
		StateValue: types.StringValue("existing"),
		State:      tfsdk.State{Raw: tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"x": tftypes.String}}, map[string]tftypes.Value{"x": tftypes.NewValue(tftypes.String, "existing")})},
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
	}
	if !resp.PlanValue.Equal(types.StringValue("existing")) {
		t.Fatalf("expected 'existing', got %s", resp.PlanValue)
	}
}

func TestUseStateForUnknownIncludingNullString_KnownPlan(t *testing.T) {
	t.Parallel()

	resp := &planmodifier.StringResponse{
		PlanValue: types.StringValue("configured"),
	}

	planmodifiers.UseStateForUnknownIncludingNullString().PlanModifyString(context.Background(), planmodifier.StringRequest{
		PlanValue:  types.StringValue("configured"),
		StateValue: types.StringValue("old"),
		State:      tfsdk.State{Raw: tftypes.NewValue(tftypes.Object{AttributeTypes: map[string]tftypes.Type{"x": tftypes.String}}, map[string]tftypes.Value{"x": tftypes.NewValue(tftypes.String, "old")})},
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
	}
	if !resp.PlanValue.Equal(types.StringValue("configured")) {
		t.Fatalf("expected 'configured', got %s", resp.PlanValue)
	}
}

func TestUseStateForUnknownIncludingNullString_NewResource(t *testing.T) {
	t.Parallel()

	resp := &planmodifier.StringResponse{
		PlanValue: types.StringUnknown(),
	}

	// New resource: state.Raw is null
	planmodifiers.UseStateForUnknownIncludingNullString().PlanModifyString(context.Background(), planmodifier.StringRequest{
		PlanValue:  types.StringUnknown(),
		StateValue: types.StringNull(),
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
