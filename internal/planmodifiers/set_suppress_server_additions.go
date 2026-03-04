package planmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// SetSuppressServerAdditions returns a plan modifier that suppresses drift when
// the API enriches a user-provided set with additional server-managed elements.
//
// Use this for computed_optional set attributes where the API automatically adds
// required or default elements that the user did not specify. For example, when
// an API always includes certain HTTP headers in a set regardless of user input.
//
// Behavior:
//   - If config is null or unknown: no action (other modifiers handle these).
//   - If state is null (new resource): no action.
//   - If every element in the config also exists in state AND state has more
//     elements: use the state value (suppress the server-added elements from
//     showing as drift).
//   - If the config contains elements not in state: normal diff (update needed).
func SetSuppressServerAdditions() planmodifier.Set {
	return setSuppressServerAdditionsModifier{}
}

type setSuppressServerAdditionsModifier struct{}

func (m setSuppressServerAdditionsModifier) Description(_ context.Context) string {
	return "Suppresses drift when the API adds server-managed elements to a user-provided set."
}

func (m setSuppressServerAdditionsModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m setSuppressServerAdditionsModifier) PlanModifySet(_ context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse) {
	// Nothing to do for null/unknown plan or null state (create)
	if req.PlanValue.IsNull() || req.PlanValue.IsUnknown() {
		return
	}
	if req.StateValue.IsNull() || req.StateValue.IsUnknown() {
		return
	}

	configElems := req.ConfigValue.Elements()
	stateElems := req.StateValue.Elements()

	// Build a set of state element string representations for lookup
	stateSet := make(map[string]struct{}, len(stateElems))
	for _, e := range stateElems {
		stateSet[e.String()] = struct{}{}
	}

	// Check if every config element exists in state
	for _, e := range configElems {
		if _, ok := stateSet[e.String()]; !ok {
			// Config has an element not in state — user is adding something new.
			// Let the normal diff proceed.
			return
		}
	}

	// All config elements are in state, and state may have extra server-added
	// elements. Use state to suppress the drift from those additions.
	if len(stateElems) >= len(configElems) {
		resp.PlanValue = req.StateValue
	}
}
