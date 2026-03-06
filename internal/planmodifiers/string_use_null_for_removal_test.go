package planmodifiers_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/planmodifiers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestUseNullForRemoval_UnknownConfig(t *testing.T) {
	t.Parallel()

	resp := &planmodifier.StringResponse{
		PlanValue: types.StringUnknown(),
	}

	planmodifiers.UseNullForRemoval().PlanModifyString(context.Background(), planmodifier.StringRequest{
		ConfigValue: types.StringUnknown(),
		StateValue:  types.StringValue("existing"),
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
	}
	if !resp.PlanValue.Equal(types.StringValue("existing")) {
		t.Fatalf("expected state value 'existing', got %s", resp.PlanValue)
	}
}

func TestUseNullForRemoval_NullConfig(t *testing.T) {
	t.Parallel()

	resp := &planmodifier.StringResponse{
		PlanValue: types.StringNull(),
	}

	planmodifiers.UseNullForRemoval().PlanModifyString(context.Background(), planmodifier.StringRequest{
		ConfigValue: types.StringNull(),
		StateValue:  types.StringValue("existing"),
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
	}
	if !resp.PlanValue.IsNull() {
		t.Fatalf("expected null, got %s", resp.PlanValue)
	}
}

func TestUseNullForRemoval_ExplicitValue(t *testing.T) {
	t.Parallel()

	resp := &planmodifier.StringResponse{
		PlanValue: types.StringValue("new-value"),
	}

	planmodifiers.UseNullForRemoval().PlanModifyString(context.Background(), planmodifier.StringRequest{
		ConfigValue: types.StringValue("new-value"),
		StateValue:  types.StringValue("old-value"),
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
	}
	if !resp.PlanValue.Equal(types.StringValue("new-value")) {
		t.Fatalf("expected 'new-value', got %s", resp.PlanValue)
	}
}

func TestUseNullForRemoval_NullState(t *testing.T) {
	t.Parallel()

	resp := &planmodifier.StringResponse{
		PlanValue: types.StringUnknown(),
	}

	planmodifiers.UseNullForRemoval().PlanModifyString(context.Background(), planmodifier.StringRequest{
		ConfigValue: types.StringUnknown(),
		StateValue:  types.StringNull(),
		State:       tfsdk.State{Raw: tftypes.NewValue(tftypes.Object{}, nil)},
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
	}
	if !resp.PlanValue.IsNull() {
		t.Fatalf("expected null (state value), got %s", resp.PlanValue)
	}
}
