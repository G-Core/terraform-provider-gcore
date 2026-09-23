package customvalidator

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var _ validator.Object = CredentialsValidator{}

// CredentialsValidator validates that either ssh_key_name is provided,
// or both username and password_wo (with password_wo_version) are provided,
// but not both authentication methods simultaneously.
//
// Used by both cloud_gpu_baremetal_cluster and cloud_gpu_virtual_cluster.
type CredentialsValidator struct{}

func (v CredentialsValidator) Description(_ context.Context) string {
	return "validates that either ssh_key_name is provided, or both username and password_wo (with password_wo_version) are provided"
}

func (v CredentialsValidator) MarkdownDescription(_ context.Context) string {
	return "validates that either ssh_key_name is provided, or both username and password_wo (with password_wo_version) are provided"
}

func (v CredentialsValidator) ValidateObject(_ context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	resp.Diagnostics.Append(checkCredentials(req.Path, req.ConfigValue)...)
}

// ValidateResolvedCredentials runs the same checks against the configured
// credentials at apply time, when every value is known. Terraform validates
// configuration, not resolved values, so a check the validator deferred on an
// unknown would otherwise never run: a key next to a password_wo_version that
// resolves only at apply would reach the API, which accepts both methods
// together. Resources call this in Create before sending anything.
func ValidateResolvedCredentials(p path.Path, credentials basetypes.ObjectValue) diag.Diagnostics {
	return checkCredentials(p, credentials)
}

// checkCredentials decides each rule on the values it can see and defers only
// when an unknown value could change the outcome. At validate and plan time
// that leaves the unknown cases to ValidateResolvedCredentials; at apply time
// nothing is unknown and every rule decides.
func checkCredentials(p path.Path, value basetypes.ObjectValue) diag.Diagnostics {
	var diags diag.Diagnostics

	if value.IsNull() || value.IsUnknown() {
		return diags
	}

	attrs := value.Attributes()

	sshKeyName := attrs["ssh_key_name"]
	username := attrs["username"]
	passwordWo := attrs["password_wo"]
	passwordWoVersion := attrs["password_wo_version"]

	known := func(v attr.Value) bool { return !v.IsNull() && !v.IsUnknown() }
	unknown := func(v attr.Value) bool { return v.IsUnknown() }

	hasSSHKey := known(sshKeyName)
	hasUsername := known(username)
	hasPasswordWo := known(passwordWo)
	hasPasswordWoVersion := known(passwordWoVersion)

	// The API requires a non-empty key name (min_length=1). Catch an empty
	// string here with a specific message rather than letting the API answer
	// with a generic validation error after the plan has already warned about
	// a re-image.
	if hasSSHKey {
		if key, ok := sshKeyName.(basetypes.StringValue); ok && key.ValueString() == "" {
			diags.AddAttributeError(
				p.AtName("ssh_key_name"),
				"Empty ssh_key_name",
				"'ssh_key_name' must not be empty. Set it to the name of an existing SSH key, or remove the 'credentials' block to use 'username' and 'password_wo' instead.",
			)

			return diags
		}
	}

	// Conflict: a known key next to any known password-method field. Two
	// known values conflict whatever the others resolve to.
	if hasSSHKey && (hasUsername || hasPasswordWo || hasPasswordWoVersion) {
		diags.AddAttributeError(
			p,
			"Conflicting credentials configuration",
			"Cannot specify 'ssh_key_name' together with 'username', 'password_wo', or 'password_wo_version'. Only one authentication method can be used.",
		)

		return diags
	}

	// Neither method: a null key cannot become one, and a pair with a null half
	// cannot complete, whatever the unknown values resolve to.
	if sshKeyName.IsNull() && (username.IsNull() || passwordWo.IsNull()) {
		diags.AddAttributeError(
			p,
			"Invalid credentials configuration",
			"Either 'ssh_key_name' must be provided, or both 'username' and 'password_wo' (with 'password_wo_version') must be provided together.",
		)
	}

	// Password without its version: decidable once the version is resolved.
	if hasPasswordWo && !hasPasswordWoVersion && !unknown(passwordWoVersion) {
		diags.AddAttributeError(
			p,
			"Missing password_wo_version",
			"When using 'password_wo', you must also provide 'password_wo_version'. This field is used to track password changes since write-only fields are not stored in state.",
		)
	}

	return diags
}
