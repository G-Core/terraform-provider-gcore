package planmodifiers_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stainless-sdks/gcore-terraform/internal/planmodifiers"
)

func TestRequiresReplaceIfConfiguredPreservingState_NullConfig_WithState(t *testing.T) {
	t.Parallel()

	resp := &planmodifier.Int64Response{
		PlanValue: types.Int64Null(),
	}

	planmodifiers.RequiresReplaceIfConfiguredPreservingState().PlanModifyInt64(context.Background(), planmodifier.Int64Request{
		ConfigValue: types.Int64Null(),
		StateValue:  types.Int64Value(42),
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
	}
	if resp.RequiresReplace {
		t.Fatal("should not require replace when config is null")
	}
	if !resp.PlanValue.Equal(types.Int64Value(42)) {
		t.Fatalf("expected state value preserved, got %s", resp.PlanValue)
	}
}

func TestRequiresReplaceIfConfiguredPreservingState_NullConfig_NullState(t *testing.T) {
	t.Parallel()

	resp := &planmodifier.Int64Response{
		PlanValue: types.Int64Null(),
	}

	planmodifiers.RequiresReplaceIfConfiguredPreservingState().PlanModifyInt64(context.Background(), planmodifier.Int64Request{
		ConfigValue: types.Int64Null(),
		StateValue:  types.Int64Null(),
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
	}
	if resp.RequiresReplace {
		t.Fatal("should not require replace when config is null")
	}
}

func TestRequiresReplaceIfConfiguredPreservingState_SameValue(t *testing.T) {
	t.Parallel()

	resp := &planmodifier.Int64Response{
		PlanValue: types.Int64Value(42),
	}

	planmodifiers.RequiresReplaceIfConfiguredPreservingState().PlanModifyInt64(context.Background(), planmodifier.Int64Request{
		ConfigValue: types.Int64Value(42),
		StateValue:  types.Int64Value(42),
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
	}
	if resp.RequiresReplace {
		t.Fatal("should not require replace when values are the same")
	}
}

func TestRequiresReplaceIfConfiguredPreservingState_DifferentValue(t *testing.T) {
	t.Parallel()

	resp := &planmodifier.Int64Response{
		PlanValue: types.Int64Value(99),
	}

	planmodifiers.RequiresReplaceIfConfiguredPreservingState().PlanModifyInt64(context.Background(), planmodifier.Int64Request{
		ConfigValue: types.Int64Value(99),
		StateValue:  types.Int64Value(42),
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
	}
	if !resp.RequiresReplace {
		t.Fatal("should require replace when values differ")
	}
}

func TestRequiresReplaceIfConfiguredPreservingState_NewResource(t *testing.T) {
	t.Parallel()

	resp := &planmodifier.Int64Response{
		PlanValue: types.Int64Value(42),
	}

	planmodifiers.RequiresReplaceIfConfiguredPreservingState().PlanModifyInt64(context.Background(), planmodifier.Int64Request{
		ConfigValue: types.Int64Value(42),
		StateValue:  types.Int64Null(),
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
	}
	if resp.RequiresReplace {
		t.Fatal("should not require replace for new resource")
	}
}
