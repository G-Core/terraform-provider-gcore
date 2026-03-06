// Package planmodifiers provides custom Terraform plan modifiers.
package planmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// baseUseStateForUnknownInclNull provides shared description methods for all
// type-specific UseStateForUnknownInclNull plan modifiers.
type baseUseStateForUnknownInclNull struct{}

func (m baseUseStateForUnknownInclNull) Description(_ context.Context) string {
	return "Preserves the prior state value (including null) when the plan value is unknown."
}

func (m baseUseStateForUnknownInclNull) MarkdownDescription(_ context.Context) string {
	return "Preserves the prior state value (including null) when the plan value is unknown."
}

// shouldUseStateValue returns true if the state value should be used.
// This is the common logic shared by all type-specific modifiers.
func shouldUseStateValue(stateRawIsNull, planValueIsUnknown bool) bool {
	// Do nothing if there is no state value (new resource)
	if stateRawIsNull {
		return false
	}
	// Only use state value if plan value is unknown
	return planValueIsUnknown
}

// BoolUseStateForUnknownInclNull returns a Bool plan modifier that preserves
// the prior state value when the plan value is unknown, INCLUDING when
// the state value is null. This differs from the standard UseStateForUnknown()
// which only preserves non-null state values.
func BoolUseStateForUnknownInclNull() planmodifier.Bool {
	return boolUseStateForUnknownInclNull{}
}

type boolUseStateForUnknownInclNull struct {
	baseUseStateForUnknownInclNull
}

func (m boolUseStateForUnknownInclNull) PlanModifyBool(_ context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	if shouldUseStateValue(req.State.Raw.IsNull(), req.PlanValue.IsUnknown()) {
		resp.PlanValue = req.StateValue
	}
}

// Int64UseStateForUnknownInclNull returns an Int64 plan modifier that preserves
// the prior state value when the plan value is unknown, INCLUDING when
// the state value is null.
func Int64UseStateForUnknownInclNull() planmodifier.Int64 {
	return int64UseStateForUnknownInclNull{}
}

type int64UseStateForUnknownInclNull struct {
	baseUseStateForUnknownInclNull
}

func (m int64UseStateForUnknownInclNull) PlanModifyInt64(_ context.Context, req planmodifier.Int64Request, resp *planmodifier.Int64Response) {
	if shouldUseStateValue(req.State.Raw.IsNull(), req.PlanValue.IsUnknown()) {
		resp.PlanValue = req.StateValue
	}
}

// StringUseStateForUnknownInclNull returns a String plan modifier that preserves
// the prior state value when the plan value is unknown, INCLUDING when
// the state value is null.
func StringUseStateForUnknownInclNull() planmodifier.String {
	return stringUseStateForUnknownInclNull{}
}

type stringUseStateForUnknownInclNull struct {
	baseUseStateForUnknownInclNull
}

func (m stringUseStateForUnknownInclNull) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if shouldUseStateValue(req.State.Raw.IsNull(), req.PlanValue.IsUnknown()) {
		resp.PlanValue = req.StateValue
	}
}

// SetUseStateForUnknownInclNull returns a Set plan modifier that preserves
// the prior state value when the plan value is unknown, INCLUDING when
// the state value is null.
func SetUseStateForUnknownInclNull() planmodifier.Set {
	return setUseStateForUnknownInclNull{}
}

type setUseStateForUnknownInclNull struct {
	baseUseStateForUnknownInclNull
}

func (m setUseStateForUnknownInclNull) PlanModifySet(_ context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse) {
	if shouldUseStateValue(req.State.Raw.IsNull(), req.PlanValue.IsUnknown()) {
		resp.PlanValue = req.StateValue
	}
}

// ObjectUseStateForUnknownInclNull returns an Object plan modifier that preserves
// the prior state value when the plan value is unknown, INCLUDING when
// the state value is null.
func ObjectUseStateForUnknownInclNull() planmodifier.Object {
	return objectUseStateForUnknownInclNull{}
}

type objectUseStateForUnknownInclNull struct {
	baseUseStateForUnknownInclNull
}

func (m objectUseStateForUnknownInclNull) PlanModifyObject(_ context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
	if shouldUseStateValue(req.State.Raw.IsNull(), req.PlanValue.IsUnknown()) {
		resp.PlanValue = req.StateValue
	}
}

// ObjectPreserveNullState returns an Object plan modifier that:
// - If plan is unknown AND state is null → plan becomes null (don't compute)
// - If plan is unknown AND state is not null → plan becomes state value
// This is useful for computed+optional nested objects that should not be
// computed when neither config nor state specifies them.
func ObjectPreserveNullState() planmodifier.Object {
	return objectPreserveNullState{}
}

type objectPreserveNullState struct{}

func (m objectPreserveNullState) Description(_ context.Context) string {
	return "Preserves the state value (including null) when the plan value is unknown."
}

func (m objectPreserveNullState) MarkdownDescription(_ context.Context) string {
	return "Preserves the state value (including null) when the plan value is unknown."
}

func (m objectPreserveNullState) PlanModifyObject(_ context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
	// Do nothing if there is no state (new resource being created)
	if req.State.Raw.IsNull() {
		return
	}

	// If plan value is unknown, use state value (even if state value is null)
	if req.PlanValue.IsUnknown() {
		resp.PlanValue = req.StateValue
	}
}
