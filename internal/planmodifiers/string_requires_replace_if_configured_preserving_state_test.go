package planmodifiers_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/planmodifiers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestStringRequiresReplaceIfConfiguredPreservingState_NullConfig_WithState(t *testing.T) {
	t.Parallel()

	resp := &planmodifier.StringResponse{
		PlanValue: types.StringNull(),
	}

	planmodifiers.StringRequiresReplaceIfConfiguredPreservingState().PlanModifyString(context.Background(), planmodifier.StringRequest{
		ConfigValue: types.StringNull(),
		StateValue:  types.StringValue("subnet"),
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
	}
	if resp.RequiresReplace {
		t.Fatal("should not require replace when config is null")
	}
}

func TestStringRequiresReplaceIfConfiguredPreservingState_NullConfig_NullState(t *testing.T) {
	t.Parallel()

	resp := &planmodifier.StringResponse{
		PlanValue: types.StringNull(),
	}

	planmodifiers.StringRequiresReplaceIfConfiguredPreservingState().PlanModifyString(context.Background(), planmodifier.StringRequest{
		ConfigValue: types.StringNull(),
		StateValue:  types.StringNull(),
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
	}
	if resp.RequiresReplace {
		t.Fatal("should not require replace when config is null")
	}
}

func TestStringRequiresReplaceIfConfiguredPreservingState_SameValue(t *testing.T) {
	t.Parallel()

	resp := &planmodifier.StringResponse{
		PlanValue: types.StringValue("subnet"),
	}

	planmodifiers.StringRequiresReplaceIfConfiguredPreservingState().PlanModifyString(context.Background(), planmodifier.StringRequest{
		ConfigValue: types.StringValue("subnet"),
		StateValue:  types.StringValue("subnet"),
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
	}
	if resp.RequiresReplace {
		t.Fatal("should not require replace when values are the same")
	}
}

func TestStringRequiresReplaceIfConfiguredPreservingState_DifferentValue(t *testing.T) {
	t.Parallel()

	resp := &planmodifier.StringResponse{
		PlanValue: types.StringValue("external"),
	}

	planmodifiers.StringRequiresReplaceIfConfiguredPreservingState().PlanModifyString(context.Background(), planmodifier.StringRequest{
		ConfigValue: types.StringValue("external"),
		StateValue:  types.StringValue("subnet"),
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
	}
	if !resp.RequiresReplace {
		t.Fatal("should require replace when values differ")
	}
}

func TestStringRequiresReplaceIfConfiguredPreservingState_NullState_ConfigSet(t *testing.T) {
	t.Parallel()

	// This is the import scenario: state is null (API didn't return the field),
	// but the user's config specifies a value. The modifier should adopt the
	// config value into the plan without requiring replacement.
	resp := &planmodifier.StringResponse{
		PlanValue: types.StringValue("subnet"),
	}

	planmodifiers.StringRequiresReplaceIfConfiguredPreservingState().PlanModifyString(context.Background(), planmodifier.StringRequest{
		ConfigValue: types.StringValue("subnet"),
		StateValue:  types.StringNull(),
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
	}
	if resp.RequiresReplace {
		t.Fatal("should not require replace when state is null (import case)")
	}
	if !resp.PlanValue.Equal(types.StringValue("subnet")) {
		t.Fatalf("expected plan value to be config value 'subnet', got %s", resp.PlanValue)
	}
}

func TestStringRequiresReplaceIfConfiguredPreservingState_UnknownState_ConfigSet(t *testing.T) {
	t.Parallel()

	resp := &planmodifier.StringResponse{
		PlanValue: types.StringValue("any_subnet"),
	}

	planmodifiers.StringRequiresReplaceIfConfiguredPreservingState().PlanModifyString(context.Background(), planmodifier.StringRequest{
		ConfigValue: types.StringValue("any_subnet"),
		StateValue:  types.StringUnknown(),
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
	}
	if resp.RequiresReplace {
		t.Fatal("should not require replace when state is unknown")
	}
	if !resp.PlanValue.Equal(types.StringValue("any_subnet")) {
		t.Fatalf("expected plan value to be config value 'any_subnet', got %s", resp.PlanValue)
	}
}

func TestStringRequiresReplaceIfConfiguredPreservingState_UnknownConfig(t *testing.T) {
	t.Parallel()

	resp := &planmodifier.StringResponse{
		PlanValue: types.StringUnknown(),
	}

	planmodifiers.StringRequiresReplaceIfConfiguredPreservingState().PlanModifyString(context.Background(), planmodifier.StringRequest{
		ConfigValue: types.StringUnknown(),
		StateValue:  types.StringValue("subnet"),
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
	}
	if resp.RequiresReplace {
		t.Fatal("should not require replace when config is unknown")
	}
}
