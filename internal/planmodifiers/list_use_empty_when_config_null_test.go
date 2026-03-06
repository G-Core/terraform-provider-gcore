package planmodifiers_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/planmodifiers"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestUseEmptyListWhenConfigNull_NullConfig_WithState(t *testing.T) {
	t.Parallel()

	stateValue := types.ListValueMust(types.StringType, []attr.Value{
		types.StringValue("a"),
	})

	resp := &planmodifier.ListResponse{
		PlanValue: stateValue,
	}

	planmodifiers.UseEmptyListWhenConfigNull().PlanModifyList(context.Background(), planmodifier.ListRequest{
		ConfigValue: types.ListNull(types.StringType),
		StateValue:  stateValue,
		PlanValue:   stateValue,
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
	}
	if len(resp.PlanValue.Elements()) != 0 {
		t.Fatalf("expected empty list, got %d elements", len(resp.PlanValue.Elements()))
	}
	if resp.PlanValue.IsNull() {
		t.Fatal("expected empty list, not null")
	}
}

func TestUseEmptyListWhenConfigNull_NullConfig_NullState(t *testing.T) {
	t.Parallel()

	resp := &planmodifier.ListResponse{
		PlanValue: types.ListNull(types.StringType),
	}

	planmodifiers.UseEmptyListWhenConfigNull().PlanModifyList(context.Background(), planmodifier.ListRequest{
		ConfigValue: types.ListNull(types.StringType),
		StateValue:  types.ListNull(types.StringType),
		PlanValue:   types.ListNull(types.StringType),
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
	}
	// During create with null state, plan should remain as-is (null)
	if !resp.PlanValue.IsNull() {
		t.Fatalf("expected null plan for create, got %s", resp.PlanValue)
	}
}

func TestUseEmptyListWhenConfigNull_ExplicitValue(t *testing.T) {
	t.Parallel()

	configValue := types.ListValueMust(types.StringType, []attr.Value{
		types.StringValue("b"),
	})

	resp := &planmodifier.ListResponse{
		PlanValue: configValue,
	}

	planmodifiers.UseEmptyListWhenConfigNull().PlanModifyList(context.Background(), planmodifier.ListRequest{
		ConfigValue: configValue,
		StateValue:  types.ListNull(types.StringType),
		PlanValue:   configValue,
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
	}
	if len(resp.PlanValue.Elements()) != 1 {
		t.Fatalf("expected 1 element, got %d", len(resp.PlanValue.Elements()))
	}
}

func TestUseEmptyListWhenConfigNull_UnknownPlan_WithState(t *testing.T) {
	t.Parallel()

	stateValue := types.ListValueMust(types.StringType, []attr.Value{
		types.StringValue("a"),
	})

	resp := &planmodifier.ListResponse{
		PlanValue: types.ListUnknown(types.StringType),
	}

	planmodifiers.UseEmptyListWhenConfigNull().PlanModifyList(context.Background(), planmodifier.ListRequest{
		ConfigValue: types.ListValueMust(types.StringType, []attr.Value{types.StringValue("a")}),
		StateValue:  stateValue,
		PlanValue:   types.ListUnknown(types.StringType),
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
	}
	if !resp.PlanValue.Equal(stateValue) {
		t.Fatalf("expected state value preserved, got %s", resp.PlanValue)
	}
}
