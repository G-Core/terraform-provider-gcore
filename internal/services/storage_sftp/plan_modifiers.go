package storage_sftp

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// PasswordArgumentsRequiredTogether checks unknown values that the config
// validator must skip. It avoids password_wo diagnostics because Terraform may
// print the source line with a literal password.
func PasswordArgumentsRequiredTogether(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() {
		return
	}

	var configPassword types.String
	var configVersion types.Int64
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("password_wo"), &configPassword)...)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("password_wo_version"), &configVersion)...)
	if resp.Diagnostics.HasError() {
		return
	}

	switch {
	case !configPassword.IsNull() && configVersion.IsNull():
		resp.Diagnostics.AddError(missingVersionSummary, missingVersionDetail)
	case configPassword.IsNull() && !configVersion.IsNull():
		resp.Diagnostics.AddAttributeError(path.Root("password_wo_version"), missingPasswordSummary, missingPasswordDetail)
	}
}

// Setting has_password after refresh turns detected password removal into an
// update.
func PlanHasPassword(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() {
		return
	}

	var configPassword types.String
	var planVersion, stateVersion types.Int64
	var stateHasPassword types.Bool
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("password_wo"), &configPassword)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("password_wo_version"), &planVersion)...)
	creating := req.State.Raw.IsNull()
	if !creating {
		resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("password_wo_version"), &stateVersion)...)
		resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("has_password"), &stateHasPassword)...)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	action, decided := passwordActionFor(creating, configPassword, planVersion, stateVersion, stateHasPassword)
	if !decided {
		return
	}

	switch action {
	case passwordSet:
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("has_password"), types.BoolValue(true))...)

		if creating || replacementPlanned(ctx, req, resp) {
			return
		}

		knownTrue := !stateHasPassword.IsNull() && !stateHasPassword.IsUnknown() && stateHasPassword.ValueBool()
		if stateVersion.IsNull() && knownTrue {
			resp.Diagnostics.AddAttributeWarning(
				path.Root("password_wo_version"),
				"Replacing a password Terraform did not set",
				`This storage already has a password (has_password = true) that was set outside Terraform — `+
					`by an import, an older provider version, or the portal. Applying replaces it in place with `+
					`"password_wo"; the old password stops working. SSH keys and every other setting are untouched.`+
					"\n\n"+
					`Remove "password_wo" and "password_wo_version" to leave the existing password alone.`,
			)
		}
	case passwordClear:
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("has_password"), types.BoolValue(false))...)
	case passwordUntouched:
	}
}

// Keep this list in sync with RequiresReplace attributes in the schema.
var replacementAttributes = []string{"name", "location_name"}

// replacementPlanned checks replacement attributes because their plan modifier
// results are not available in this response.
func replacementPlanned(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) bool {
	for _, attribute := range replacementAttributes {
		var planned, prior types.String
		resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root(attribute), &planned)...)
		resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root(attribute), &prior)...)
		if resp.Diagnostics.HasError() {
			return false
		}
		if !planned.Equal(prior) {
			return true
		}
	}
	return false
}
