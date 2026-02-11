package planmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// UseStateForUnknownIncludingNullString returns a plan modifier that preserves the
// state value when the plan value is unknown, INCLUDING when the state value is null.
//
// This differs from the built-in UseStateForUnknown() which leaves unknown plans as
// unknown when state is null. Use this when a computed_optional attribute should preserve
// its null state across updates rather than being re-computed.
func UseStateForUnknownIncludingNullString() planmodifier.String {
	return useStateForUnknownIncludingNullStringModifier{}
}

type useStateForUnknownIncludingNullStringModifier struct{}

func (m useStateForUnknownIncludingNullStringModifier) Description(_ context.Context) string {
	return "Preserves the state value (including null) when the plan value is unknown"
}

func (m useStateForUnknownIncludingNullStringModifier) MarkdownDescription(_ context.Context) string {
	return "Preserves the state value (including null) when the plan value is unknown"
}

func (m useStateForUnknownIncludingNullStringModifier) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// Do nothing if there is no state (resource is being created)
	if req.StateValue.IsNull() && req.PlanValue.IsUnknown() && req.State.Raw.IsNull() {
		return
	}

	// If the plan value is not unknown, do nothing
	if !req.PlanValue.IsUnknown() {
		return
	}

	// Preserve the state value (including null) in the plan
	resp.PlanValue = req.StateValue
}
