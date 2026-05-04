// Package planmodifiers provides custom Terraform plan modifiers.
package planmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ObjectUseNullForRemoval is the object-typed analogue of UseNullForRemoval. It:
//   - Uses the state value when the config value is unknown and state has a value
//     (computed behavior during updates)
//   - Keeps the plan value unknown when the config value is unknown and state is
//     null (resource creation referencing another resource's computed attribute)
//   - Uses null when the config value is explicitly null (removal behavior)
//   - Uses the config value when explicitly set
//
// This solves the "computed_optional removal" problem for nested object
// attributes where Terraform doesn't distinguish between "not set" and
// "explicitly removed".
func ObjectUseNullForRemoval() planmodifier.Object {
	return objectUseNullForRemovalModifier{}
}

type objectUseNullForRemovalModifier struct{}

func (m objectUseNullForRemovalModifier) Description(_ context.Context) string {
	return "Uses null when the config value is explicitly removed, preserving state only when unknown."
}

func (m objectUseNullForRemovalModifier) MarkdownDescription(_ context.Context) string {
	return "Uses null when the config value is explicitly removed, preserving state only when unknown."
}

func (m objectUseNullForRemovalModifier) PlanModifyObject(ctx context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
	// Config is unknown (e.g. references another resource's computed attribute):
	// fall back to state during updates; leave plan unknown during create.
	if req.ConfigValue.IsUnknown() {
		if req.StateValue.IsNull() {
			resp.PlanValue = types.ObjectUnknown(req.PlanValue.AttributeTypes(ctx))
			return
		}
		resp.PlanValue = req.StateValue
		return
	}

	// Config is explicitly null (omitted from .tf): mark plan null so the
	// resource is updated to remove the value rather than preserving state.
	if req.ConfigValue.IsNull() {
		resp.PlanValue = types.ObjectNull(req.PlanValue.AttributeTypes(ctx))
		return
	}

	// Config has a value: leave the framework's already-computed plan alone so
	// inner attribute defaults and plan modifiers continue to apply. Setting
	// resp.PlanValue = req.ConfigValue here would discard those defaults
	// (e.g. a Default: booldefault.StaticBool(true) on a child attribute).
}
