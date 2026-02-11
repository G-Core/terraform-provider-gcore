package planmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// RequiresReplaceIfConfiguredPreservingState returns a plan modifier that only requires
// replacement when the user has explicitly configured a different value. If the config
// value is null (user relies on environment variables or defaults), the state value is
// preserved without forcing replacement.
//
// This differs from int64planmodifier.RequiresReplaceIfConfigured() by also preserving
// the state value when config is null, which handles import scenarios where
// project_id/region_id are in state (from import ID) but not in the user's .tf config.
func RequiresReplaceIfConfiguredPreservingState() planmodifier.Int64 {
	return requiresReplaceIfConfiguredPreservingStateModifier{}
}

type requiresReplaceIfConfiguredPreservingStateModifier struct{}

func (m requiresReplaceIfConfiguredPreservingStateModifier) Description(_ context.Context) string {
	return "Requires replacement only when the user has explicitly configured a different value. " +
		"If config is null (e.g., using environment variables), preserves state value without replacement."
}

func (m requiresReplaceIfConfiguredPreservingStateModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m requiresReplaceIfConfiguredPreservingStateModifier) PlanModifyInt64(_ context.Context, req planmodifier.Int64Request, resp *planmodifier.Int64Response) {
	// If config is null (user didn't specify the value, likely using env vars),
	// preserve the state value. This handles import scenarios where project_id/region_id
	// are in state (from import ID) but not in the user's .tf config.
	if req.ConfigValue.IsNull() {
		if !req.StateValue.IsNull() && !req.StateValue.IsUnknown() {
			// Preserve state value - user wants to use the imported/existing value
			resp.PlanValue = req.StateValue
		}
		// Don't require replacement - user wants to use env vars / defaults
		return
	}

	// If creating a new resource (state is null), no replacement needed
	if req.StateValue.IsNull() {
		return
	}

	// If config has an explicit value different from state, require replacement
	if !req.ConfigValue.Equal(req.StateValue) {
		resp.RequiresReplace = true
	}
}
