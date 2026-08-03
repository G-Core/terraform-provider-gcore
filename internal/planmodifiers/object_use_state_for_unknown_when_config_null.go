package planmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// ObjectUseStateForUnknownWhenConfigNull replaces an unknown planned object
// with the prior state value — including a null prior state, and including at
// create time when no prior state exists — but only when the attribute is
// null in configuration.
//
// Computed+Optional object attributes with no Default are marked "(known
// after apply)" whenever the config omits them. When the resource model
// stores the attribute as a plain Go pointer (*Model), that unknown cannot be
// decoded at apply time and req.Plan.Get fails with "Received unknown value,
// however the target type cannot handle unknown values". This modifier
// resolves the unknown at plan time instead.
//
// Only use this on attributes the API never populates implicitly when they
// are omitted from the request: the planned value becomes known (null on
// create), so a server-supplied value in the create/update response would
// make Terraform fail with "Provider produced inconsistent result after
// apply". For attributes the API does populate from the response, use
// UseStateForUnknownIncludingNullObject instead.
func ObjectUseStateForUnknownWhenConfigNull() planmodifier.Object {
	return objectUseStateForUnknownWhenConfigNullModifier{}
}

type objectUseStateForUnknownWhenConfigNullModifier struct{}

func (m objectUseStateForUnknownWhenConfigNullModifier) Description(_ context.Context) string {
	return "Uses the prior state value (null on create) for an unknown planned value when the attribute is omitted from configuration"
}

func (m objectUseStateForUnknownWhenConfigNullModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m objectUseStateForUnknownWhenConfigNullModifier) PlanModifyObject(_ context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
	if !req.PlanValue.IsUnknown() {
		return
	}
	// A non-null (possibly unknown) config value is not "omitted": an unknown
	// config expression resolves by the apply-time re-plan, and a concrete
	// block never yields a whole-object unknown. Leave the plan alone.
	if !req.ConfigValue.IsNull() {
		return
	}
	// req.StateValue is a typed null object during create, so this yields null
	// when there is no prior value and preserves the prior value otherwise.
	resp.PlanValue = req.StateValue
}
