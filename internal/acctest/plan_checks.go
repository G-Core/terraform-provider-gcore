package acctest

import (
	"context"
	"fmt"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

// ExpectNotUnknownValue returns a plan check asserting that the attribute at the
// given path is NOT planned as "(known after apply)".
//
// plancheck.ExpectKnownValue cannot express this: it reads the plan's "after"
// object, where an unknown value and a known null are both JSON null. This check
// reads "after_unknown" instead, so it distinguishes a computed attribute whose
// prior state was preserved (including a preserved null) from one that the
// framework re-marked as unknown because it lacks a plan modifier.
//
// "after_unknown" is a SPARSE structure. Terraform only emits entries for the
// parts of the planned object that are, or contain, unknown values; anything
// fully known is omitted entirely. So:
//
//   - path absent from "after_unknown"             -> fully known -> pass
//   - path present, false or all-false when nested -> fully known -> pass
//   - path present with any true inside            -> unknown     -> fail
//
// Because absence from "after_unknown" is what makes this check pass, the same
// path is also traversed in "after" so that a bad path cannot pass silently.
// Note that "after" is sparse too, but in the opposite direction: Terraform runs
// the planned object through an omit-unknowns filter, so a fully unknown
// attribute is dropped from "after" entirely, while a KNOWN null is kept as JSON
// null. A path missing from "after" therefore means one of two things -- the
// attribute is planned as unknown, or the name is wrong -- and both are
// failures, so the check reports them together.
//
// Nested unknowns count: if the path resolves to an object or a list, the check
// fails when any value inside it is unknown.
func ExpectNotUnknownValue(resourceAddress string, attributePath tfjsonpath.Path) plancheck.PlanCheck {
	return expectNotUnknownValue{
		resourceAddress: resourceAddress,
		attributePath:   attributePath,
	}
}

type expectNotUnknownValue struct {
	resourceAddress string
	attributePath   tfjsonpath.Path
}

func (e expectNotUnknownValue) CheckPlan(_ context.Context, req plancheck.CheckPlanRequest, resp *plancheck.CheckPlanResponse) {
	if req.Plan == nil {
		resp.Error = fmt.Errorf("plan is nil")
		return
	}

	for _, resourceChange := range req.Plan.ResourceChanges {
		if resourceChange.Address != e.resourceAddress {
			continue
		}

		// "after" keeps known nulls but drops fully unknown attributes, so a path
		// that does not resolve there is either planned as unknown or misspelled.
		if _, err := tfjsonpath.Traverse(resourceChange.Change.After, e.attributePath); err != nil {
			resp.Error = fmt.Errorf("%s.%s: absent from the planned value, meaning it is either planned as "+
				"unknown (\"known after apply\") because it lacks a stabilising plan modifier, or the "+
				"attribute name is wrong: %w",
				e.resourceAddress, e.attributePath.String(), err)
			return
		}

		// "after_unknown" is sparse: a path that does not resolve there is fully
		// known, which is exactly what this check wants. The path itself has
		// already been validated against "after" above.
		afterUnknown, err := tfjsonpath.Traverse(resourceChange.Change.AfterUnknown, e.attributePath)
		if err != nil {
			return
		}

		if containsUnknown(afterUnknown) {
			resp.Error = fmt.Errorf("%s.%s is planned as unknown (\"known after apply\"); "+
				"expected the prior state value to be preserved",
				e.resourceAddress, e.attributePath.String())
		}
		return
	}

	resp.Error = fmt.Errorf("%s - resource not found in plan", e.resourceAddress)
}

// containsUnknown reports whether an "after_unknown" fragment marks anything
// unknown. Terraform encodes it as true, false, or a nested object/array of the
// same.
func containsUnknown(afterUnknown any) bool {
	switch value := afterUnknown.(type) {
	case bool:
		return value
	case map[string]any:
		for _, nested := range value {
			if containsUnknown(nested) {
				return true
			}
		}
	case []any:
		for _, nested := range value {
			if containsUnknown(nested) {
				return true
			}
		}
	}
	return false
}

// ExpectResourceNotReplaced returns a plan check asserting that the resource is
// not being destroyed, recreated or replaced. Creating and updating both pass,
// as does a no-op.
//
// plancheck.ExpectResourceAction cannot express this: it pins the action to one
// exact value, so a test guarding "import must not replace this" has to name the
// action it does expect and then fails whenever an unrelated attribute changes
// whether the plan is an update or a no-op. This check states the invariant
// directly, which keeps import-no-replace regression tests from being coupled to
// whichever incidental diffs an import currently leaves behind.
func ExpectResourceNotReplaced(resourceAddress string) plancheck.PlanCheck {
	return expectResourceNotReplaced{resourceAddress: resourceAddress}
}

type expectResourceNotReplaced struct {
	resourceAddress string
}

func (e expectResourceNotReplaced) CheckPlan(_ context.Context, req plancheck.CheckPlanRequest, resp *plancheck.CheckPlanResponse) {
	if req.Plan == nil {
		resp.Error = fmt.Errorf("plan is nil")
		return
	}

	for _, resourceChange := range req.Plan.ResourceChanges {
		if resourceChange.Address != e.resourceAddress {
			continue
		}

		for _, action := range resourceChange.Change.Actions {
			if action == tfjson.ActionDelete {
				resp.Error = fmt.Errorf("%s is planned for %s (actions: %v); expected it to be kept in place",
					e.resourceAddress, action, resourceChange.Change.Actions)
				return
			}
		}

		return
	}

	resp.Error = fmt.Errorf("%s - resource not found in plan", e.resourceAddress)
}
