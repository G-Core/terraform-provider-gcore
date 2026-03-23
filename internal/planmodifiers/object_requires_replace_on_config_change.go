package planmodifiers

import (
	"context"

	t "github.com/G-Core/terraform-provider-gcore/internal/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// RequiresReplaceOnConfigChange returns a plan modifier that triggers replacement
// only when user-specified fields in the config actually change from state,
// ignoring computed fields that may show as unknown during planning.
//
// The optional ignoreFields parameter specifies top-level attribute names within
// this object that should be excluded from the comparison (e.g. fields that are
// updatable in-place and should not force replacement).
//
// This is useful for nested objects that contain both user-specified (Required/Optional)
// and computed fields. Using the built-in RequiresReplace() would trigger replacement
// whenever any computed field becomes unknown, even if the user didn't change anything.
func RequiresReplaceOnConfigChange(ignoreFields ...string) planmodifier.Object {
	ignore := make(map[string]struct{}, len(ignoreFields))
	for _, f := range ignoreFields {
		ignore[f] = struct{}{}
	}
	return requiresReplaceOnConfigChangeModifier{ignoreFields: ignore}
}

type requiresReplaceOnConfigChangeModifier struct {
	ignoreFields map[string]struct{}
}

func (m requiresReplaceOnConfigChangeModifier) Description(_ context.Context) string {
	return "Requires replacement only when user-specified server settings change"
}

func (m requiresReplaceOnConfigChangeModifier) MarkdownDescription(_ context.Context) string {
	return "Requires replacement only when user-specified server settings change"
}

func (m requiresReplaceOnConfigChangeModifier) PlanModifyObject(_ context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
	// If there's no state (new resource), no replacement needed
	if req.StateValue.IsNull() || req.StateValue.IsUnknown() {
		return
	}

	// If config is null but state exists, resource is being removed
	if req.ConfigValue.IsNull() {
		return
	}

	// Compare config value with state value to determine if user-specified fields changed.
	// We compare config (what user wrote) vs state (what was applied), not plan vs state,
	// because plan may have computed fields marked as unknown.
	if configValuesChangedFromState(req.ConfigValue, req.StateValue, m.ignoreFields) {
		resp.RequiresReplace = true
	}
}

// configValuesChangedFromState performs a deep comparison between config and state values,
// only checking fields that are explicitly set in config (not null/unknown).
// This allows computed fields in state to differ without triggering replacement.
// The ignoreFields set contains top-level attribute names to skip during comparison;
// it is only applied at the first level of recursion and passed as nil for deeper levels.
func configValuesChangedFromState(configVal, stateVal attr.Value, ignoreFields map[string]struct{}) bool {
	// If config value is null or unknown, consider it unchanged (computed field)
	if configVal.IsNull() || configVal.IsUnknown() {
		return false
	}

	// If state is null but config is not, something changed
	if stateVal.IsNull() {
		return true
	}

	// Handle objects/maps - compare only attributes present in config
	if ok, configAttrs := t.ChildAttributes(configVal); ok {
		if ok, stateAttrs := t.ChildAttributes(stateVal); ok {
			for field, configFieldVal := range configAttrs {
				// Skip explicitly ignored fields
				if _, ignored := ignoreFields[field]; ignored {
					continue
				}

				// Skip fields not specified in config
				if configFieldVal.IsNull() || configFieldVal.IsUnknown() {
					continue
				}

				stateFieldVal, exists := stateAttrs[field]
				if !exists {
					return true
				}

				// Don't propagate ignoreFields to deeper levels
				if configValuesChangedFromState(configFieldVal, stateFieldVal, nil) {
					return true
				}
			}
			return false
		}
		return true
	}

	// Handle lists/tuples/sets - compare elements by index
	if ok, configElems := t.ChildItems(configVal); ok {
		if ok, stateElems := t.ChildItems(stateVal); ok {
			if len(configElems) != len(stateElems) {
				return true
			}
			for i, configElem := range configElems {
				if configValuesChangedFromState(configElem, stateElems[i], nil) {
					return true
				}
			}
			return false
		}
		return true
	}

	// For primitive types, use direct comparison
	return !configVal.Equal(stateVal)
}
