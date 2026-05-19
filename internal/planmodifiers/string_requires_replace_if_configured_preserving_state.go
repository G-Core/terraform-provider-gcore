package planmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// StringRequiresReplaceIfConfiguredPreservingState returns a plan modifier that only
// requires replacement when the user has explicitly configured a different value AND
// the prior state has a known, non-null value to compare against.
//
// This is import-safe: if the state value is null (e.g., after importing a resource
// whose API response doesn't include this write-only field), the config value is
// adopted directly into the plan without triggering replacement or showing a diff.
// This prevents unnecessary destroy+recreate cycles and eliminates spurious
// in-place update diffs after terraform import.
func StringRequiresReplaceIfConfiguredPreservingState() planmodifier.String {
	return stringRequiresReplaceIfConfiguredPreservingStateModifier{}
}

type stringRequiresReplaceIfConfiguredPreservingStateModifier struct{}

func (m stringRequiresReplaceIfConfiguredPreservingStateModifier) Description(_ context.Context) string {
	return "Requires replacement only when the user has explicitly configured a different value " +
		"and the prior state has a known value. Import-safe: adopts the config value into the " +
		"plan when state is null (e.g., after importing a resource with a write-only field), " +
		"suppressing both replacement and spurious update diffs."
}

func (m stringRequiresReplaceIfConfiguredPreservingStateModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m stringRequiresReplaceIfConfiguredPreservingStateModifier) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// If config is null (user didn't specify the value), no replacement needed.
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	// If state is null or unknown (new resource or post-import with write-only field),
	// don't require replacement. Adopt the config value directly into the plan so
	// that importing doesn't produce a spurious diff for fields the API cannot return.
	if req.StateValue.IsNull() || req.StateValue.IsUnknown() {
		resp.PlanValue = req.ConfigValue
		return
	}

	// Both config and state have known values — require replacement if they differ.
	if !req.ConfigValue.Equal(req.StateValue) {
		resp.RequiresReplace = true
	}
}
