package planmodifiers_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/planmodifiers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestRequiresReplaceIfPriorValueKnown_NullState_Create(t *testing.T) {
	t.Parallel()

	resp := &planmodifier.StringResponse{
		PlanValue: types.StringValue("cert-pem"),
	}

	planmodifiers.RequiresReplaceIfPriorValueKnown().PlanModifyString(context.Background(), planmodifier.StringRequest{
		ConfigValue: types.StringValue("cert-pem"),
		PlanValue:   types.StringValue("cert-pem"),
		StateValue:  types.StringNull(),
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
	}
	if resp.RequiresReplace {
		t.Fatal("expected no replacement on create (null state)")
	}
}

func TestRequiresReplaceIfPriorValueKnown_NullState_ImportScenario(t *testing.T) {
	t.Parallel()

	// After import, state is null because API doesn't return the PEM.
	// Config has the cert PEM. Should NOT trigger replacement.
	resp := &planmodifier.StringResponse{
		PlanValue: types.StringValue("cert-pem"),
	}

	planmodifiers.RequiresReplaceIfPriorValueKnown().PlanModifyString(context.Background(), planmodifier.StringRequest{
		ConfigValue: types.StringValue("cert-pem"),
		PlanValue:   types.StringValue("cert-pem"),
		StateValue:  types.StringNull(),
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
	}
	if resp.RequiresReplace {
		t.Fatal("expected no replacement after import (null state, config set)")
	}
}

func TestRequiresReplaceIfPriorValueKnown_SameValue_NoReplacement(t *testing.T) {
	t.Parallel()

	resp := &planmodifier.StringResponse{
		PlanValue: types.StringValue("cert-pem"),
	}

	planmodifiers.RequiresReplaceIfPriorValueKnown().PlanModifyString(context.Background(), planmodifier.StringRequest{
		ConfigValue: types.StringValue("cert-pem"),
		PlanValue:   types.StringValue("cert-pem"),
		StateValue:  types.StringValue("cert-pem"),
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
	}
	if resp.RequiresReplace {
		t.Fatal("expected no replacement when state and plan are the same")
	}
}

func TestRequiresReplaceIfPriorValueKnown_RemovedFromConfig_RequiresReplace(t *testing.T) {
	t.Parallel()

	// The attribute is create-only and the user deleted it from the config.
	// The Update API cannot honour the removal, so the resource must be
	// replaced rather than silently nulled out in state.
	resp := &planmodifier.StringResponse{
		PlanValue: types.StringNull(),
	}

	planmodifiers.RequiresReplaceIfPriorValueKnown().PlanModifyString(context.Background(), planmodifier.StringRequest{
		ConfigValue: types.StringNull(),
		PlanValue:   types.StringNull(),
		StateValue:  types.StringValue("net-a"),
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
	}
	if !resp.RequiresReplace {
		t.Fatal("expected replacement when a configured value is removed from the config")
	}
}

func TestRequiresReplaceIfPriorValueKnown_UnknownPlanValue_RequiresReplace(t *testing.T) {
	t.Parallel()

	// The config references an attribute of a resource that is created or
	// replaced in the same apply, so the value is unknown at plan time. It is
	// still a change away from the known prior value.
	resp := &planmodifier.StringResponse{
		PlanValue: types.StringUnknown(),
	}

	planmodifiers.RequiresReplaceIfPriorValueKnown().PlanModifyString(context.Background(), planmodifier.StringRequest{
		ConfigValue: types.StringUnknown(),
		PlanValue:   types.StringUnknown(),
		StateValue:  types.StringValue("net-a"),
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
	}
	if !resp.RequiresReplace {
		t.Fatal("expected replacement when the planned value is unknown and prior state is known")
	}
}

func TestRequiresReplaceIfPriorValueKnown_DifferentValue_RequiresReplace(t *testing.T) {
	t.Parallel()

	resp := &planmodifier.StringResponse{
		PlanValue: types.StringValue("new-cert-pem"),
	}

	planmodifiers.RequiresReplaceIfPriorValueKnown().PlanModifyString(context.Background(), planmodifier.StringRequest{
		ConfigValue: types.StringValue("new-cert-pem"),
		PlanValue:   types.StringValue("new-cert-pem"),
		StateValue:  types.StringValue("old-cert-pem"),
	}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error: %s", resp.Diagnostics.Errors())
	}
	if !resp.RequiresReplace {
		t.Fatal("expected replacement when state and plan differ")
	}
}
