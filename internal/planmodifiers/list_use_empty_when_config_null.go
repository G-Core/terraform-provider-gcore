package planmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// UseEmptyListWhenConfigNull returns a plan modifier that sets the plan value to
// an empty list when the config value is null (attribute omitted from .tf file).
//
// Use this for computed_optional list attributes where omitting the attribute
// should mean "clear the list" rather than "keep state value".
func UseEmptyListWhenConfigNull() planmodifier.List {
	return useEmptyListWhenConfigNullModifier{}
}

type useEmptyListWhenConfigNullModifier struct{}

func (m useEmptyListWhenConfigNullModifier) Description(_ context.Context) string {
	return "Sets plan to empty list when config is null, preserving state when plan is unknown during create."
}

func (m useEmptyListWhenConfigNullModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m useEmptyListWhenConfigNullModifier) PlanModifyList(ctx context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
	// If config is null (user omitted the attribute), set plan to empty list
	// This ensures "no members in config" means "clear all members"
	if req.ConfigValue.IsNull() {
		// Only set to empty list if we have existing state (update scenario)
		// During create, let the plan remain unknown so API can provide defaults
		if !req.StateValue.IsNull() && !req.StateValue.IsUnknown() {
			resp.PlanValue = types.ListValueMust(req.StateValue.ElementType(ctx), []attr.Value{})
		}
		return
	}

	// If plan is unknown and state has a value, use state (like UseStateForUnknown)
	// This handles computed defaults during create
	if req.PlanValue.IsUnknown() && !req.StateValue.IsNull() && !req.StateValue.IsUnknown() {
		resp.PlanValue = req.StateValue
	}
}
