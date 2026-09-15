package storage_sftp

import (
	"github.com/G-Core/gcore-go/option"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type passwordAction int

const (
	passwordUntouched passwordAction = iota
	passwordSet
	passwordClear
)

func (a passwordAction) String() string {
	switch a {
	case passwordSet:
		return "set"
	case passwordClear:
		return "none"
	default:
		return "untouched"
	}
}

// A sibling-only PATCH must not resend an ephemeral password because its value
// may change on every run. decided is false when planning needs a known version.
func passwordActionFor(creating bool, configPassword types.String, planVersion, stateVersion types.Int64, stateHasPassword types.Bool) (action passwordAction, decided bool) {
	pairConfigured := !configPassword.IsNull()

	if creating {
		// POST always requires password_mode.
		if pairConfigured {
			return passwordSet, true
		}
		return passwordClear, true
	}

	if !pairConfigured {
		if !stateVersion.IsNull() {
			// Terraform managed this password, so removing the pair clears it.
			return passwordClear, true
		}
		// Keep passwords that Terraform did not manage.
		return passwordUntouched, true
	}

	if planVersion.IsUnknown() {
		// Adoption and drift repair always send the password. Other cases need
		// the resolved version.
		if stateVersion.IsNull() || isKnownFalse(stateHasPassword) {
			return passwordSet, true
		}
		return passwordUntouched, false
	}

	if planVersion.IsNull() {
		return passwordUntouched, true
	}

	if !planVersion.Equal(stateVersion) {
		return passwordSet, true
	}

	if isKnownFalse(stateHasPassword) {
		// Repair a password removed outside Terraform.
		return passwordSet, true
	}

	return passwordUntouched, true
}

func isKnownFalse(b types.Bool) bool {
	return !b.IsNull() && !b.IsUnknown() && !b.ValueBool()
}

// passwordRequestOptions must run after option.WithRequestBody. That option
// replaces the full body and would discard fields added before it. Clear mode
// must not include a password field.
func passwordRequestOptions(action passwordAction, configPassword types.String) []option.RequestOption {
	switch action {
	case passwordSet:
		return []option.RequestOption{
			option.WithJSONSet("password_mode", passwordModeSet),
			option.WithJSONSet("password", configPassword.ValueString()),
		}
	case passwordClear:
		return []option.RequestOption{
			option.WithJSONSet("password_mode", passwordModeNone),
		}
	default:
		return nil
	}
}
