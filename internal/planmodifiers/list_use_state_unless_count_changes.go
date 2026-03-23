package planmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// UseStateUnlessCountChanges returns a plan modifier that preserves the state value
// for a list attribute unless:
//   - The resource is being created (no prior state)
//   - The resource is being replaced (id becomes unknown)
//   - The specified count attribute is changing (which affects the list)
//
// This is useful for computed list attributes like servers_ids that change when
// the cluster is resized (servers_count changes) but should be preserved otherwise.
func UseStateUnlessCountChanges(countAttr string) planmodifier.List {
	return useStateUnlessCountChangesModifier{countAttr: countAttr}
}

type useStateUnlessCountChangesModifier struct {
	countAttr string
}

func (m useStateUnlessCountChangesModifier) Description(_ context.Context) string {
	return "Preserves state value unless resource is replaced or " + m.countAttr + " changes"
}

func (m useStateUnlessCountChangesModifier) MarkdownDescription(_ context.Context) string {
	return "Preserves state value unless resource is replaced or " + m.countAttr + " changes"
}

func (m useStateUnlessCountChangesModifier) PlanModifyList(ctx context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
	// If there's no state (new resource), nothing to preserve
	if req.State.Raw.IsNull() {
		return
	}

	// If the planned value is already known, don't override it
	if !resp.PlanValue.IsUnknown() {
		return
	}

	// Check if the resource is being replaced by checking if id is becoming unknown
	var planID types.String
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("id"), &planID)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if planID.IsUnknown() {
		// Resource is being replaced, don't preserve state
		return
	}

	// Check if the count attribute is changing
	var stateCount, planCount types.Int64
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root(m.countAttr), &stateCount)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root(m.countAttr), &planCount)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !stateCount.Equal(planCount) {
		// Count is changing, list will change too
		return
	}

	// Safe to preserve state value
	resp.PlanValue = req.StateValue
}
