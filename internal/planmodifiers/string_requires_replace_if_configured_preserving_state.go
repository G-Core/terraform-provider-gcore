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
// whose API response doesn't include this write-only field), replacement is NOT triggered
// even if the config has a value. This prevents unnecessary destroy+recreate cycles
// after terraform import.
func StringRequiresReplaceIfConfiguredPreservingState() planmodifier.String {
	return stringRequiresReplaceIfConfiguredPreservingStateModifier{}
}

type stringRequiresReplaceIfConfiguredPreservingStateModifier struct{}

func (m stringRequiresReplaceIfConfiguredPreservingStateModifier) Description(_ context.Context) string {
	return "Requires replacement only when the user has explicitly configured a different value " +
		"and the prior state has a known value. Import-safe: does not force replacement when " +
		"state is null (e.g., after importing a resource with a write-only field)."
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
	// don't require replacement. We can't compare against a value we don't have.
	if req.StateValue.IsNull() || req.StateValue.IsUnknown() {
		return
	}

	// Both config and state have known values — require replacement if they differ.
	if !req.ConfigValue.Equal(req.StateValue) {
		resp.RequiresReplace = true
	}
}
