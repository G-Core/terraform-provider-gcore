package planmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// ListRequiresReplaceIfNotNull returns a plan modifier that requires resource
// replacement when the list value changes, but only if the prior state value
// is not null.
//
// This is import-safe: after terraform import, if the state value is null
// (e.g., because the field uses no_refresh and the API doesn't return it),
// replacement is NOT triggered even though the config has a value. Instead,
// Terraform plans an update-in-place that populates the field in state from
// the user's config. After that one-time apply, future changes trigger
// replacement normally.
//
// For a Required field, null state can only occur after import (since creation
// always requires the field to be set).
func ListRequiresReplaceIfNotNull() planmodifier.List {
	return listRequiresReplaceIfNotNullModifier{}
}

type listRequiresReplaceIfNotNullModifier struct{}

func (m listRequiresReplaceIfNotNullModifier) Description(_ context.Context) string {
	return "Requires replacement when the value changes, but skips replacement " +
		"when prior state is null (e.g., after terraform import with a no_refresh field)."
}

func (m listRequiresReplaceIfNotNullModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m listRequiresReplaceIfNotNullModifier) PlanModifyList(_ context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
	// During resource creation there is no prior state — skip.
	if req.State.Raw.IsNull() {
		return
	}

	// If prior state value is null (post-import for no_refresh fields), skip
	// replacement. The plan will show an update-in-place that populates the
	// field in state from the user's config. The Update function handles this
	// as a no-op (no API call) but still persists the value to state.
	if req.StateValue.IsNull() {
		return
	}

	// If the plan value equals the state value, no replacement needed.
	if req.PlanValue.Equal(req.StateValue) {
		return
	}

	resp.RequiresReplace = true
}
