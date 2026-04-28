package planmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// RequiresReplaceIfPriorValueKnown returns a plan modifier that requires replacement
// only when the prior state value is known (non-null). This handles the import-then-plan
// scenario where the API doesn't return a field (e.g., ssl_certificate PEM) on GET:
//
//   - Create: state is null → no replacement (correct, resource is being created)
//   - Import then plan: state is null (API didn't return it) → no replacement
//     (certificate already exists in cloud, no need to recreate)
//   - Real config change: state has old value → replacement triggered (correct)
func RequiresReplaceIfPriorValueKnown() planmodifier.String {
	return requiresReplaceIfPriorValueKnownModifier{}
}

type requiresReplaceIfPriorValueKnownModifier struct{}

func (m requiresReplaceIfPriorValueKnownModifier) Description(_ context.Context) string {
	return "Requires replacement only when the prior state value is known (non-null), " +
		"avoiding spurious replacement after import when the API does not return the field."
}

func (m requiresReplaceIfPriorValueKnownModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m requiresReplaceIfPriorValueKnownModifier) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// If there's no prior state (create or import with null field), skip replacement.
	if req.StateValue.IsNull() || req.StateValue.IsUnknown() {
		return
	}

	// Prior state exists and differs from planned value → require replacement.
	if !req.PlanValue.Equal(req.StateValue) {
		resp.RequiresReplace = true
	}
}
