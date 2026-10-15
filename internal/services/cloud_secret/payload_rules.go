package cloud_secret

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.ResourceWithValidateConfig = (*CloudSecretResource)(nil)

const (
	missingPayloadSummary = "Missing secret payload"
	missingPayloadDetail  = "\"payload\" and \"payload_wo_version\" are required to create a secret, including " +
		"when a change (for example to \"name\" or \"expiration\") forces an existing secret to be replaced: the " +
		"API cannot update a secret or read its payload back. They may only be omitted for an existing secret " +
		"that is not being replaced, for example one imported without its original material."

	unpairedPayloadSummary = "\"payload\" and \"payload_wo_version\" must be set together"
	unpairedPayloadDetail  = "Configure both to create the secret, or to replace it with new material. Omit both " +
		"to manage an existing secret without its material; removing both from an existing secret keeps it " +
		"unchanged."
)

// payloadUnpaired reports whether exactly one of payload/payload_wo_version
// is configured. Either being unknown returns false: core validates the
// config again with known values in the apply-time plan, so skipping unknown
// values is safe here. This mirrors the storage_sftp config validator.
func payloadUnpaired(payload types.Object, version types.Int64) bool {
	if payload.IsUnknown() || version.IsUnknown() {
		return false
	}
	return payload.IsNull() != version.IsNull()
}

// payloadMissing reports whether either payload or payload_wo_version is
// null. Unknown is not null, so an unknown value never errors here either;
// the apply-time re-plan catches it once both are known.
func payloadMissing(payload types.Object, version types.Int64) bool {
	return payload.IsNull() || version.IsNull()
}

func (r *CloudSecretResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var payload types.Object
	var version types.Int64

	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("payload"), &payload)...)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("payload_wo_version"), &version)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if payloadUnpaired(payload, version) {
		resp.Diagnostics.AddError(unpairedPayloadSummary, unpairedPayloadDetail)
	}
}

// requirePayloadOnCreate is also the replacement guard. For every
// replacement (a RequiresReplace attribute, -replace, taint,
// create_before_destroy), Terraform core plans the create half again with a
// NULL prior state (terraform internal/terraform/node_resource_abstract_instance.go,
// action.IsReplace()) and fails the plan on any error from that call. The
// resource-level ModifyPlan cannot see attribute RequiresReplace decisions
// (resp.RequiresReplace starts empty). Do NOT gate this on
// resp.RequiresReplace, on the first plan, or on anything else besides
// State.Raw being null.
//
// No framework validator may target payload or its children: a
// resourcevalidator.RequiredTogether, AlsoRequires, ConflictsWith, or an
// objectvalidator attaches its diagnostic to configuredPaths[0], and
// Terraform then prints the offending config source line - which can hold
// the private key. Every diagnostic here is path-less.
func requirePayloadOnCreate(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if !(req.State.Raw.IsNull() && !req.Plan.Raw.IsNull()) {
		return
	}

	var payload types.Object
	var version types.Int64
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("payload"), &payload)...)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("payload_wo_version"), &version)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if payloadMissing(payload, version) {
		resp.Diagnostics.AddError(missingPayloadSummary, missingPayloadDetail)
	}
}
