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
