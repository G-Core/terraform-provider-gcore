package customvalidator

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	attrs := req.ConfigValue.Attributes()

	sshKeyName := attrs["ssh_key_name"]
	username := attrs["username"]
	passwordWo := attrs["password_wo"]
	passwordWoVersion := attrs["password_wo_version"]

	hasSSHKey := !sshKeyName.IsNull() && !sshKeyName.IsUnknown()
	hasUsername := !username.IsNull() && !username.IsUnknown()
	hasPasswordWo := !passwordWo.IsNull() && !passwordWo.IsUnknown()
	hasPasswordWoVersion := !passwordWoVersion.IsNull() && !passwordWoVersion.IsUnknown()

	if hasSSHKey && (hasUsername || hasPasswordWo || hasPasswordWoVersion) {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Conflicting credentials configuration",
			"Cannot specify 'ssh_key_name' together with 'username', 'password_wo', or 'password_wo_version'. Only one authentication method can be used.",
		)
	} else if !hasSSHKey && !(hasUsername && hasPasswordWo) {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid credentials configuration",
			"Either 'ssh_key_name' must be provided, or both 'username' and 'password_wo' (with 'password_wo_version') must be provided together.",
		)
	} else if hasPasswordWo && !hasPasswordWoVersion {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Missing password_wo_version",
			"When using 'password_wo', you must also provide 'password_wo_version'. This field is used to track password changes since write-only fields are not stored in state.",
		)
	}
}
