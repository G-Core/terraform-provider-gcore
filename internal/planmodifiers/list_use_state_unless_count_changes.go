package planmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// UseStateUnlessCountChanges returns a plan modifier that preserves the state value
// for a list attribute unless:
//   - The resource is being created or replaced (no prior state; Terraform core
//     re-plans a replacement create with null prior state)
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
	// If there's no state (new resource, or the create half of a replace —
	// Terraform core re-plans replacement creates with null prior state),
	// there is nothing to preserve.
	if req.State.Raw.IsNull() {
		return
	}

	// If the planned value is already known, don't override it
	if !resp.PlanValue.IsUnknown() {
		return
	}

	// Note: replacement is intentionally NOT detected via plan.id being
	// unknown. Attribute plan modifiers run in nondeterministic order, so
	// this modifier may observe id before its own UseStateForUnknown-style
	// modifier restores it, wrongly treating a no-op plan as a replace.

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
