package cloud_gpu_baremetal_cluster

import (
	"context"
	"fmt"

	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/G-Core/terraform-provider-gcore/internal/planmodifiers"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// This file holds the custom logic behind the re-image warning: which
// attribute changes re-image a cluster's servers, how that decision is read out
// of a plan, and the diagnostics built on it. resource.go only calls into it.
//
// Invariant: ModifyPlan warns and Update acts on ONE decision, rebuildsOnApply,
// read from the same planned values. A warning about a re-image that does not
// happen, or a re-image nobody was warned about, would each be worse than the
// silence this replaces.

// attributeGetter is the part of tfsdk.Plan and tfsdk.State this file needs.
type attributeGetter interface {
	GetAttribute(ctx context.Context, p path.Path, target interface{}) diag.Diagnostics
}

// rebuildTriggers holds the attributes the API hashes into a cluster's spec. A
// change to any of them makes the running servers out of date, and applying it
// re-images them. The set mirrors the API's spec hash for the fields this
// resource exposes; username and password_wo_version are not hashed, cannot be
// updated, and are rejected by the credential guard before this is consulted.
type rebuildTriggers struct {
	imageID    types.String
	sshKeyName types.String
	userData   types.String
}

// rebuildTriggersOf reads the three triggers out of a plan or a prior state.
//
// It asks only for the two root attributes and descends by hand. A nested
// GetAttribute path through an unknown parent does not fail: the framework
// returns a NULL leaf (fwschemadata.Data.ValueAtPath swallows
// tftypes.ErrInvalidStep), and null would read as "no change".
func rebuildTriggersOf(ctx context.Context, src attributeGetter) rebuildTriggers {
	var imageID types.String
	if src.GetAttribute(ctx, path.Root("image_id"), &imageID).HasError() {
		imageID = types.StringUnknown()
	}

	var settings types.Object
	if src.GetAttribute(ctx, path.Root("servers_settings"), &settings).HasError() {
		return rebuildTriggers{
			imageID:    imageID,
			sshKeyName: types.StringUnknown(),
			userData:   types.StringUnknown(),
		}
	}

	return rebuildTriggers{
		imageID:    imageID,
		sshKeyName: stringAt(objectAt(settings, "credentials"), "ssh_key_name"),
		userData:   stringAt(settings, "user_data"),
	}
}

// objectAt returns the named object attribute of obj. It fails closed: a null
// parent gives a null child, but an unknown parent, a missing attribute, or a
// value of the wrong type all give UNKNOWN. credentialsObject in resource.go
// answers null in those cases, which is right for the credential guard and
// wrong here, where an unreadable value might be a re-image nobody was warned
// about.
func objectAt(obj types.Object, name string) types.Object {
	if obj.IsNull() {
		return types.ObjectNull(map[string]attr.Type{})
	}
	if obj.IsUnknown() {
		return types.ObjectUnknown(map[string]attr.Type{})
	}
	val, ok := obj.Attributes()[name]
	if !ok {
		return types.ObjectUnknown(map[string]attr.Type{})
	}
	child, ok := val.(types.Object)
	if !ok {
		return types.ObjectUnknown(map[string]attr.Type{})
	}
	return child
}

// stringAt returns the named string attribute of obj, under the same
// fail-closed rule as objectAt.
func stringAt(obj types.Object, name string) types.String {
	if obj.IsNull() {
		return types.StringNull()
	}
	if obj.IsUnknown() {
		return types.StringUnknown()
	}
	val, ok := obj.Attributes()[name]
	if !ok {
		return types.StringUnknown()
	}
	str, ok := val.(types.String)
	if !ok {
		return types.StringUnknown()
	}
	return str
}

// rebuildDecision is the answer rebuildsOnApply gives about a plan. Three
// values, not two: ModifyPlan can honestly say "it might", while Update has to
// either re-image or not.
type rebuildDecision int

const (
	// rebuildNo means nothing that is applied by a re-image is changing.
	rebuildNo rebuildDecision = iota

	// rebuildYes means a hashed attribute is changing to a known value, so the
	// apply reaches ApplySettingsAndPoll.
	rebuildYes

	// rebuildUnknown means a hashed attribute is planned as unknown, or could
	// not be read, so the answer is not decided until the value is.
	rebuildUnknown
)

func (d rebuildDecision) String() string {
	switch d {
	case rebuildNo:
		return "rebuildNo"
	case rebuildYes:
		return "rebuildYes"
	case rebuildUnknown:
		return "rebuildUnknown"
	default:
		return fmt.Sprintf("rebuildDecision(%d)", int(d))
	}
}

// rebuildsOnApply answers whether this plan re-images the cluster's servers.
// It takes the trigger structs rather than the getters so that Update builds
// its request body from the same values it decided on.
func rebuildsOnApply(planned, prior rebuildTriggers) rebuildDecision {
	return leastResolved(
		rebuildDecisionFor(planned.imageID, prior.imageID),
		leastResolved(
			rebuildDecisionFor(planned.sshKeyName, prior.sshKeyName),
			rebuildDecisionFor(planned.userData, prior.userData),
		),
	)
}

// rebuildDecisionFor answers the same question for a single attribute.
//
// A null planned value is never a re-image: there is nothing to send. An
// unknown value on either side cannot be compared - types.String.ValueString
// reports "" for it, so a naive comparison would read an unknown planned value
// against a null prior as "no change".
//
// The comparison is plan against prior state, which is correct for a
// computed_optional attribute under the pinned framework: when a configuration
// omits user_data, stringplanmodifier.UseStateForUnknown copies the prior value
// into the plan, null included (v1.19.0 use_state_for_unknown.go documents null
// as a known value it copies). So a rename plans user_data equal to prior and
// answers rebuildNo. TestRebuildTriggerAttributeShapes pins that modifier.
func rebuildDecisionFor(plan, state types.String) rebuildDecision {
	if plan.IsNull() {
		return rebuildNo
	}

	if plan.IsUnknown() || state.IsUnknown() {
		return rebuildUnknown
	}

	if plan.Equal(state) {
		return rebuildNo
	}

	return rebuildYes
}

// leastResolved combines per-attribute answers, and an unresolved attribute
// outranks a definite one: Update may not act while any single trigger is
// unresolved, or it would re-image to apply the definite one while dropping
// the unresolved one from the payload.
func leastResolved(a, b rebuildDecision) rebuildDecision {
	switch {
	case a == rebuildUnknown || b == rebuildUnknown:
		return rebuildUnknown
	case a == rebuildYes || b == rebuildYes:
		return rebuildYes
	default:
		return rebuildNo
	}
}

// sshKeyRemoved reports whether the plan drops a key the cluster has. There is
// no API path for it: the settings PATCH cannot send a null ssh_key_name, and a
// running server keeps the key it was imaged with. Left alone, Update would
// send nothing and never write state, and Terraform would report an
// inconsistent result after apply.
//
// An empty string counts as removal too. CredentialsValidator rejects a
// literal "" at validate time; this covers a value that was unknown then and
// resolved to "" at plan, which the API would otherwise reject with a generic
// validation error after the re-image warning had already been shown.
func sshKeyRemoved(planned, prior rebuildTriggers) bool {
	blank := func(v types.String) bool { return v.IsNull() || (!v.IsUnknown() && v.ValueString() == "") }
	return blank(planned.sshKeyName) && !prior.sshKeyName.IsUnknown() && !blank(prior.sshKeyName)
}

// planReplacesCluster reports whether this plan replaces the cluster because
// an attribute with a RequiresReplace-family modifier changed. On such a plan
// Create runs instead of Update: nothing is re-imaged, nothing is patched, and
// the credential guards do not apply, since a new cluster can be built with
// any credentials or none.
//
// The attribute modifiers run first and record their answer in the
// framework's own response; the resource's ModifyPlan is handed a separate,
// empty RequiresReplace. So the rules are re-derived here from the values the
// modifiers compared. TestReplacementRulesMatchSchema fails if a modifier is
// added, dropped or changed by codegen.
//
// A `terraform plan -replace=ADDR` or a replace_triggered_by is resolved by
// core after this returns and never reaches the provider, so it is not
// detected: on such a plan the guards below still apply and the warning may
// still fire. The diagnostics and docs say so rather than assert an action.
func planReplacesCluster(ctx context.Context, plan tfsdk.Plan, state tfsdk.State, config tfsdk.Config) bool {
	// Without a config there is nothing to compare for the config-driven rules,
	// and the test harness builds requests without one.
	haveConfig := config.Schema != nil && !config.Raw.IsNull()

	// flavor: RequiresReplace - the modifier's rule is !plan.Equal(state), and
	// an unknown planned value counts as a change, the same way core treats an
	// unknown on a replace path.
	var planFlavor, stateFlavor types.String
	if !plan.GetAttribute(ctx, path.Root("flavor"), &planFlavor).HasError() &&
		!state.GetAttribute(ctx, path.Root("flavor"), &stateFlavor).HasError() &&
		!planFlavor.Equal(stateFlavor) {
		return true
	}

	// project_id, region_id: RequiresReplaceIfConfigured - a change while the
	// attribute is set in the configuration.
	if haveConfig {
		for _, name := range []string{"project_id", "region_id"} {
			var configVal, planVal, stateVal types.Int64
			if config.GetAttribute(ctx, path.Root(name), &configVal).HasError() ||
				plan.GetAttribute(ctx, path.Root(name), &planVal).HasError() ||
				state.GetAttribute(ctx, path.Root(name), &stateVal).HasError() {
				continue
			}
			if !configVal.IsNull() && !planVal.Equal(stateVal) {
				return true
			}
		}

		// servers_settings: RequiresReplaceOnConfigChange, ignoring the two
		// fields that are applied by a re-image instead.
		var configSettings, stateSettings types.Object
		if !config.GetAttribute(ctx, path.Root("servers_settings"), &configSettings).HasError() &&
			!state.GetAttribute(ctx, path.Root("servers_settings"), &stateSettings).HasError() &&
			planmodifiers.ConfigChangedFromState(configSettings, stateSettings, serversSettingsReplaceIgnore...) {
			return true
		}
	}

	return false
}

// serversSettingsReplaceIgnore is the ignore list handed to
// RequiresReplaceOnConfigChange on servers_settings in schema.go. It is
// duplicated here rather than read from the modifier because the modifier
// exposes nothing but its description.
var serversSettingsReplaceIgnore = []string{"credentials", "user_data"}

// applyTriggerPayload fills the trigger-derived fields of a cluster settings
// PATCH from the planned values, never from the model: by the time the
// settings PATCH is built, earlier branches have folded API responses into
// the model, and UnmarshalComputed does not honour user_data's no_refresh tag,
// so the model can carry the API's user_data where the plan had none.
func applyTriggerPayload(params *cloud.GPUBaremetalClusterUpdateParams, planned rebuildTriggers, imageChanged bool) {
	if imageChanged {
		params.ImageID = param.NewOpt(planned.imageID.ValueString())
	}
	if !planned.sshKeyName.IsNull() && !planned.sshKeyName.IsUnknown() {
		params.ServersSettings.Credentials.SSHKeyName = param.NewOpt(planned.sshKeyName.ValueString())
	}
	if !planned.userData.IsNull() && !planned.userData.IsUnknown() {
		params.ServersSettings.UserData = param.NewOpt(planned.userData.ValueString())
	}
}

// Diagnostics. None of them states which action Terraform planned: a
// `terraform plan -replace=ADDR` reaches the provider as an ordinary update, so
// a text asserting an in-place update could contradict the plan printed above
// it. The re-image warning also fires on an unresolved value, where Update may
// yet find nothing to do; over-warning is the intended side to err on.
const (
	rebuildWarningSummary = "Applying this change re-images servers in the cluster"

	rebuildWarningDetail = "image_id, servers_settings.credentials.ssh_key_name or " +
		"servers_settings.user_data is changing. The API applies that by reinstalling the " +
		"operating system on every server whose settings are out of date, and anything on " +
		"their local disks is lost.\n\n" +
		"To avoid it, revert the attribute to its previous value and plan again. If the plan " +
		"above replaces the cluster, reverting does not prevent that. See the resource " +
		"documentation for how re-imaging works."

	rebuildUnknownSummary = "Unknown image_id, ssh_key_name or user_data at apply time"

	rebuildUnknownDetail = "Terraform resolves every planned value before applying it, so " +
		"image_id, servers_settings.credentials.ssh_key_name and servers_settings.user_data " +
		"should be known or absent by now. One arrived unknown, or could not be read from " +
		"the plan.\n\n" +
		"Re-imaging on an unresolved value would reinstall the servers while sending the API " +
		"an empty value, so the cluster settings were left alone. Any other change in this " +
		"plan may already have reached the API; none of it was written to state, and the " +
		"refresh at the start of the next plan reconciles it. This is a bug in the provider; " +
		"please report it."

	sshKeyRemovalSummary = "Unsupported credential removal"

	sshKeyRemovalDetail = "Removing servers_settings.credentials.ssh_key_name is not supported: " +
		"the API has no operation that takes a key away from running servers. Keep the key in " +
		"the configuration. To build the cluster again without it, change an attribute that " +
		"forces replacement (such as flavor) or destroy and recreate it; a -replace on the " +
		"command line does not bypass this check."
)

func rebuildWarning() diag.Diagnostic {
	return diag.NewWarningDiagnostic(rebuildWarningSummary, rebuildWarningDetail)
}

// rebuildAtUnknownValueError stops the settings rollout in Update. The guard
// sits after the name/tags/resize branches so an unreadable trigger cannot
// block them, and it writes nothing to state: the framework returns the prior
// state on error, and the planned model would record a key the API never
// received.
func rebuildAtUnknownValueError() diag.Diagnostic {
	return diag.NewErrorDiagnostic(rebuildUnknownSummary, rebuildUnknownDetail)
}

func sshKeyRemovalError() diag.Diagnostic {
	return diag.NewErrorDiagnostic(sshKeyRemovalSummary, sshKeyRemovalDetail)
}
