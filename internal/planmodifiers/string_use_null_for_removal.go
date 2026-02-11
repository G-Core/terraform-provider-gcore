// Package planmodifiers provides custom Terraform plan modifiers.
package planmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// UseNullForRemoval returns a plan modifier that:
// - Uses the state value when the config value is unknown (computed behavior)
// - Uses null when the config value is explicitly null (removal behavior)
// - Uses the config value when explicitly set
//
// This solves the "computed_optional removal" problem where Terraform doesn't
// distinguish between "not set" and "explicitly removed".
func UseNullForRemoval() planmodifier.String {
	return useNullForRemovalModifier{}
}

type useNullForRemovalModifier struct{}

func (m useNullForRemovalModifier) Description(_ context.Context) string {
	return "Uses null when the config value is explicitly removed, preserving state only when unknown."
}

func (m useNullForRemovalModifier) MarkdownDescription(_ context.Context) string {
	return "Uses null when the config value is explicitly removed, preserving state only when unknown."
}

func (m useNullForRemovalModifier) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// If the config value is unknown, use the state value (computed behavior)
	if req.ConfigValue.IsUnknown() {
		resp.PlanValue = req.StateValue
		return
	}

	// If the config value is null (explicitly removed or never set),
	// explicitly set the plan value to null - this will trigger an update to remove the value
	if req.ConfigValue.IsNull() {
		resp.PlanValue = types.StringNull()
		return
	}

	// Otherwise, use the config value
	resp.PlanValue = req.ConfigValue
}
