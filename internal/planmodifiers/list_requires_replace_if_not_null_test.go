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

func TestListRequiresReplaceIfNotNull_NullState_NoReplace(t *testing.T) {
	t.Parallel()

	planValue := types.ListValueMust(types.StringType, []attr.Value{
		types.StringValue("external"),
	})

	resp := &planmodifier.ListResponse{
		PlanValue: planValue,
	}

	planmodifiers.ListRequiresReplaceIfNotNull().PlanModifyList(context.Background(), planmodifier.ListRequest{
		StateValue: types.ListNull(types.StringType),
		PlanValue:  planValue,
		State: tfsdk.State{
			Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
		},
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
	}
	if resp.RequiresReplace {
		t.Fatal("expected no replacement when state value is null (import scenario)")
	}
	// Plan value should be preserved (config value) — an update-in-place will populate state
	if !resp.PlanValue.Equal(planValue) {
		t.Fatal("expected plan value to be preserved when state is null (import scenario)")
	}
}

func TestListRequiresReplaceIfNotNull_NullEntireState_NoReplace(t *testing.T) {
	t.Parallel()

	planValue := types.ListValueMust(types.StringType, []attr.Value{
		types.StringValue("external"),
	})

	resp := &planmodifier.ListResponse{
		PlanValue: planValue,
	}

	planmodifiers.ListRequiresReplaceIfNotNull().PlanModifyList(context.Background(), planmodifier.ListRequest{
		StateValue: types.ListNull(types.StringType),
		PlanValue:  planValue,
		State: tfsdk.State{
			Raw: tftypes.NewValue(tftypes.Object{}, nil),
		},
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
	}
	if resp.RequiresReplace {
		t.Fatal("expected no replacement during resource creation")
	}
}

func TestListRequiresReplaceIfNotNull_ValuesMatch_NoReplace(t *testing.T) {
	t.Parallel()

	value := types.ListValueMust(types.StringType, []attr.Value{
		types.StringValue("external"),
	})

	resp := &planmodifier.ListResponse{
		PlanValue: value,
	}

	planmodifiers.ListRequiresReplaceIfNotNull().PlanModifyList(context.Background(), planmodifier.ListRequest{
		StateValue: value,
		PlanValue:  value,
		State: tfsdk.State{
			Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
		},
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
	}
	if resp.RequiresReplace {
		t.Fatal("expected no replacement when values match")
	}
}

func TestListRequiresReplaceIfNotNull_ValuesDiffer_RequiresReplace(t *testing.T) {
	t.Parallel()

	stateValue := types.ListValueMust(types.StringType, []attr.Value{
		types.StringValue("external"),
	})
	planValue := types.ListValueMust(types.StringType, []attr.Value{
		types.StringValue("subnet"),
	})

	resp := &planmodifier.ListResponse{
		PlanValue: planValue,
	}

	planmodifiers.ListRequiresReplaceIfNotNull().PlanModifyList(context.Background(), planmodifier.ListRequest{
		StateValue: stateValue,
		PlanValue:  planValue,
		State: tfsdk.State{
			Raw: tftypes.NewValue(tftypes.Object{}, map[string]tftypes.Value{}),
		},
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
	}
	if !resp.RequiresReplace {
		t.Fatal("expected replacement when values differ")
	}
}
