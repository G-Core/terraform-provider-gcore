package planmodifiers

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// StringUseStateForCaseInsensitiveEqual returns a plan modifier for a
// Computed+Optional attribute validated case-insensitively but reported by the
// API in canonical case (cloud_volume type_name / volume_type): when the
// configured value equals prior state under strings.EqualFold, the plan keeps
// the prior (canonical) value, so a following RequiresReplace* modifier sees no
// change. Terraform accepts a planned value equal to the prior value for a
// Computed attribute even when the configuration differs.
//
// Place it before the replace modifier. It does nothing on create, on a null
// or unknown configured value, or on a null prior value.
//
// "No replacement" is not "no plan": a configured value that differs from
// prior state only in case still makes the proposed new state differ from
// prior state, and the framework marks every Computed attribute without a
// UseStateForUnknown modifier unknown BEFORE attribute plan modifiers run
// (internal/fwserver MarkComputedNilsAsUnknown). The plan therefore shows an
// in-place update of known-after-apply computed attributes on every run, the
// apply sends nothing and state keeps the canonical value; only the canonical
// spelling gives a clean plan. Observed on real infra (QA S7).
func StringUseStateForCaseInsensitiveEqual() planmodifier.String {
	return stringUseStateForCaseInsensitiveEqual{}
}

type stringUseStateForCaseInsensitiveEqual struct{}

func (stringUseStateForCaseInsensitiveEqual) Description(_ context.Context) string {
	return "Keeps the prior state value when the configured value equals it ignoring case."
}

func (m stringUseStateForCaseInsensitiveEqual) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (stringUseStateForCaseInsensitiveEqual) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.State.Raw.IsNull() || req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	if req.StateValue.IsNull() || req.StateValue.IsUnknown() {
		return
	}
	if strings.EqualFold(req.ConfigValue.ValueString(), req.StateValue.ValueString()) {
		resp.PlanValue = req.StateValue
	}
}
