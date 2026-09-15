package storage_sftp

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// The provider derives password_mode. It never uses "auto" because that mode
// creates and returns a random password.
const (
	passwordModeSet  = "set"
	passwordModeNone = "none"
)

// the API accepts passwords of 8 to 63 bytes
const (
	minPasswordLength = 8
	maxPasswordLength = 63
)

const (
	invalidLengthSummary = "Invalid password_wo length"
	invalidLengthDetail  = `"password_wo" string length must be between 8 and 63, got: %d.`

	missingVersionSummary = "Missing password_wo_version"
	missingVersionDetail  = `"password_wo" and "password_wo_version" must be configured together.` +
		"\n\n" +
		`Write-only values are never stored, so "password_wo_version" is the only record that ` +
		`Terraform set this storage's password. Without it the password would be sent once ` +
		`and could then never be rotated or cleared. Add "password_wo_version = 1".`

	missingPasswordSummary = "Missing password_wo"
	missingPasswordDetail  = `"password_wo" and "password_wo_version" must be configured together.` +
		"\n\n" +
		`Add "password_wo" back, or remove "password_wo_version" as well to clear password ` +
		`authentication from the storage in place.`
)

var _ resource.ConfigValidator = passwordPairValidator{}

// Terraform may print an attribute's source line without redacting a literal
// password. This validator avoids the password path and leaves unknown values
// for the plan check.
type passwordPairValidator struct{}

func (passwordPairValidator) Description(ctx context.Context) string {
	return "password_wo must be 8-63 characters and configured together with password_wo_version"
}

func (passwordPairValidator) MarkdownDescription(ctx context.Context) string {
	return "`password_wo` must be 8-63 characters and configured together with `password_wo_version`"
}

func (passwordPairValidator) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var configPassword types.String
	var configVersion types.Int64
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("password_wo"), &configPassword)...)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("password_wo_version"), &configVersion)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// a schema validator on password_wo would attribute its diagnostic to that
	// attribute, and Terraform prints the offending source line — the literal
	// password — so the length check must live here, without a path. It must
	// run before the unknown guard so it still catches a known password paired
	// with an unknown version
	if !configPassword.IsNull() && !configPassword.IsUnknown() {
		if l := len(configPassword.ValueString()); l < minPasswordLength || l > maxPasswordLength {
			resp.Diagnostics.AddError(invalidLengthSummary, fmt.Sprintf(invalidLengthDetail, l))
		}
	}

	if configPassword.IsUnknown() || configVersion.IsUnknown() {
		return
	}

	switch {
	case !configPassword.IsNull() && configVersion.IsNull():
		resp.Diagnostics.AddError(missingVersionSummary, missingVersionDetail)
	case configPassword.IsNull() && !configVersion.IsNull():
		resp.Diagnostics.AddError(missingPasswordSummary, missingPasswordDetail)
	}
}
