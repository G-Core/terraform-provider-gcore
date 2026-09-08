package cloud_baremetal_server

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// importAdoptionPrivateKey marks a bare metal server whose create-only
// attributes still have to be adopted from config into state after a
// `terraform import`.
//
// It is written by ImportState, read by the plan modifiers on those attributes
// and by ModifyPlan, and cleared by Update once the adoption has happened.
//
// The marker exists because "prior state is null" does NOT by itself mean "just
// imported". It is equally true of a server created with the attribute omitted,
// and a user may set such an attribute on day two. Treating that as an adoption
// would plan a clean in-place update, send nothing the API acts on, and write
// the configured value to state - leaving state describing a server that does
// not exist. Keying on an explicit import marker keeps that case on the
// replacement path, which is the only path that applies the value.
//
// Unlike cloud_load_balancer, an adoption here cannot be corroborated against
// the API: the adopted attributes are no_refresh precisely because the server
// GET never echoes them. The residual risk is a marker that outlives its import
// - a config setting none of these attributes produces no-op plans, so no
// Update ever runs to retire it - and the mitigation is therefore a plan-time
// warning naming every attribute being adopted, rather than a verification.
// See ModifyPlan.
const importAdoptionPrivateKey = "baremetal_import_adoption_pending"

var importAdoptionPrivateValue = []byte(`{"pending":true}`)

// privateStateReader and privateStateWriter are the subsets of the framework's
// private state data used here. The concrete type lives in an internal framework
// package, so it is reached through interfaces rather than named directly.
type privateStateReader interface {
	GetKey(ctx context.Context, key string) ([]byte, diag.Diagnostics)
}

type privateStateWriter interface {
	SetKey(ctx context.Context, key string, value []byte) diag.Diagnostics
}

// markImportAdoptionPending records that the create-only attributes of a freshly
// imported server still have to be adopted into state.
func markImportAdoptionPending(ctx context.Context, private privateStateWriter) diag.Diagnostics {
	if private == nil {
		return nil
	}

	return private.SetKey(ctx, importAdoptionPrivateKey, importAdoptionPrivateValue)
}

// importAdoptionPending reports whether the server is in the window between an
// import and the apply that adopts its create-only attributes.
func importAdoptionPending(ctx context.Context, private privateStateReader) (bool, diag.Diagnostics) {
	if private == nil {
		return false, nil
	}

	value, diags := private.GetKey(ctx, importAdoptionPrivateKey)

	return len(value) > 0, diags
}

// clearImportAdoptionPending closes the adoption window. Setting an empty value
// deletes the key.
func clearImportAdoptionPending(ctx context.Context, private privateStateWriter) diag.Diagnostics {
	if private == nil {
		return nil
	}

	return private.SetKey(ctx, importAdoptionPrivateKey, nil)
}

// createOnlyAttributes maps each attribute that an adopting apply records into
// state without sending it to the API, to a predicate reporting whether the
// model has no value for it.
//
// image_id and user_data are in this list even though they are not on adopting
// plan modifiers: they are normally applied by a rebuild, and Update suppresses
// that rebuild while adopting, so during the adoption window they are written
// to state exactly like the others. ImportState back-fills flavor, image_id and
// ssh_key_name from the API response, so in practice those are rarely adopted -
// the predicate simply skips them when they are already populated.
//
// The names must match the `tfsdk:` tags in model.go, since they are what the
// warning shows the user.
var createOnlyAttributes = map[string]func(CloudBaremetalServerModel) bool{
	"apptemplate_id": func(m CloudBaremetalServerModel) bool { return m.ApptemplateID.IsNull() },
	"flavor":         func(m CloudBaremetalServerModel) bool { return m.Flavor.IsNull() },
	"image_id":       func(m CloudBaremetalServerModel) bool { return m.ImageID.IsNull() },
	"name_template":  func(m CloudBaremetalServerModel) bool { return m.NameTemplate.IsNull() },
	"ssh_key_name":   func(m CloudBaremetalServerModel) bool { return m.SSHKeyName.IsNull() },
	"user_data":      func(m CloudBaremetalServerModel) bool { return m.UserData.IsNull() },
	"username":       func(m CloudBaremetalServerModel) bool { return m.Username.IsNull() },
	"app_config":     func(m CloudBaremetalServerModel) bool { return m.AppConfig == nil || *m.AppConfig == nil },
}

// adoptedCreateOnlyAttributes returns the attributes this plan would record into
// state without sending them to the API: configured now, absent from prior
// state. The result is sorted so the warning text is stable.
func adoptedCreateOnlyAttributes(plan, state CloudBaremetalServerModel) []string {
	var adopted []string

	for name, isUnset := range createOnlyAttributes {
		if !isUnset(state) || isUnset(plan) {
			continue
		}
		adopted = append(adopted, name)
	}

	slices.Sort(adopted)

	return adopted
}

// replacementForcingAttributes maps every attribute whose plan modifier in
// schema.go can force a destroy+recreate to the rule that modifier applies.
// `password_wo` is absent on purpose: a write-only attribute is null in both
// plan and prior state, so its RequiresReplace can never fire.
// TestReplacementForcingAttributesMatchSchema fails if codegen adds or drops one.
var replacementForcingAttributes = map[string]replacementRule{
	"app_config":          replacesUnlessAdopting,
	"apptemplate_id":      replacesUnlessAdopting,
	"flavor":              replacesUnlessAdopting,
	"interfaces":          replacesOverKnownPrior,
	"name_template":       replacesUnlessAdopting,
	"password_wo_version": replacesOnAnyChange,
	"project_id":          replacesWhenConfigured,
	"region_id":           replacesWhenConfigured,
	"ssh_key_name":        replacesUnlessAdopting,
	"username":            replacesUnlessAdopting,
}

// replacementRule is the shape of every rule above. adopting matters to exactly
// one of them, but the table is uniform so a rule cannot be bound to a
// signature that silently drops it.
type replacementRule func(plan, state, config attr.Value, adopting bool) bool

// replacesOverKnownPrior is the rule of ListRequiresReplaceIfNotNull: a null
// prior value never replaces, whether or not an import is being adopted.
func replacesOverKnownPrior(plan, state, _ attr.Value, _ bool) bool {
	return !state.IsNull() && !plan.Equal(state)
}

// replacesUnlessAdopting is the rule of the *RequiresReplaceUnlessAdopting
// family. It replaces on any change, and carves out only the adoption that
// follows an import - which is not the same as carving out every null prior
// value. Setting one of these attributes for the first time on a resource that
// was created without it does replace, because replacement is the only path
// that applies the value.
func replacesUnlessAdopting(plan, state, _ attr.Value, adopting bool) bool {
	return !plan.Equal(state) && !(state.IsNull() && adopting)
}

// replacesWhenConfigured is the rule of RequiresReplaceIfConfigured.
func replacesWhenConfigured(plan, state, config attr.Value, _ bool) bool {
	return !config.IsNull() && !plan.Equal(state)
}

// replacesOnAnyChange is the rule of RequiresReplace.
func replacesOnAnyChange(plan, state, _ attr.Value, _ bool) bool {
	return !plan.Equal(state)
}

// planForcesReplacement reports whether some attribute already forces a
// destroy+recreate, in which case nothing is being adopted: Terraform runs
// Create and sends every configured value to the API.
//
// ModifyPlan cannot read the framework's replacement decision - the attribute
// plan modifiers that make it run earlier, against a response of their own, and
// ModifyPlanResponse.RequiresReplace arrives empty - so the same question is put
// to the values those modifiers compared, under the same rules.
//
// It covers replacement forced by an attribute change, and nothing else. A
// `terraform plan -replace=ADDR` and the `replace_triggered_by` lifecycle
// argument both land in core's forceReplace list - the second of them before
// planning even begins - but neither is ever put in the plan request, and core
// picks the final action only after the provider has answered. So no rule added
// here could see them. A tainted resource needs no rule because it arrives with
// a null prior state, which ModifyPlan returns on already.
//
// Moving the warnings to the second, null-prior-state PlanResourceChange call
// that core makes on a replacement is not an option either. Core takes that
// call's diagnostics only when at least one of them is an error, so as not to
// repeat warnings the first call already produced - which means a warning
// raised there on its own is dropped.
func planForcesReplacement(ctx context.Context, plan tfsdk.Plan, state tfsdk.State, config tfsdk.Config, adopting bool) bool {
	for name, forcesReplacement := range replacementForcingAttributes {
		attrPath := path.Root(name)

		var planValue, stateValue, configValue attr.Value

		// Advisory like the warning itself, so a decode problem must neither
		// fail the plan nor suppress the warning.
		if plan.GetAttribute(ctx, attrPath, &planValue).HasError() ||
			state.GetAttribute(ctx, attrPath, &stateValue).HasError() ||
			config.GetAttribute(ctx, attrPath, &configValue).HasError() {
			continue
		}

		if forcesReplacement(planValue, stateValue, configValue, adopting) {
			return true
		}
	}

	return false
}

// importAdoptionWarning names the attributes being adopted and the one way this
// can go wrong, which is a marker that outlived its import.
//
// Like the rebuild warning, this cannot assert which action Terraform planned -
// see the rebuildWarningSummary const block for why. Adoption is what an
// in-place update does; a replacement runs Create and sends every configured
// value to the API, so on that plan nothing below applies. The text says so
// rather than describing the update case as though it were the only one.
func importAdoptionWarning(adopted []string) diag.Diagnostic {
	return diag.NewWarningDiagnostic(
		"Adopting create-only attributes into state after import",
		fmt.Sprintf(
			"If the plan above shows this resource being replaced, none of the following "+
				"applies: replacing it creates the server again and every configured value "+
				"reaches the API.\n\n"+
				"Otherwise Terraform is recording %s from configuration without sending them "+
				"anywhere, because this server was imported and the API does not return these "+
				"attributes. The configured values are trusted as-is. Adopting them does not "+
				"by itself rebuild the server, which is deliberate: rebuilding would reinstall "+
				"the operating system of a server that was only just imported. If a disk erasure "+
				"warning appears alongside this one, some other change is driving a rebuild, and "+
				"that rebuild sends the adopted values too.\n\n"+
				"If the configured values do not match the running server, state will be wrong "+
				"and no future refresh can correct it. Check them before applying. To force the "+
				"values to actually take effect the server has to be replaced: run `terraform "+
				"apply -replace=` for this resource.",
			strings.Join(quoteAll(adopted), ", "),
		),
	)
}

// quoteAll wraps every name in double quotes for the warning text.
func quoteAll(names []string) []string {
	quoted := make([]string, 0, len(names))
	for _, name := range names {
		quoted = append(quoted, fmt.Sprintf("%q", name))
	}

	return quoted
}

// rebuildWarningSummary and rebuildWarningDetail are shown at plan time when a
// change to image_id or user_data will reinstall the server.
//
// Neither text may state which action Terraform planned, because ModifyPlan
// cannot know it. A `terraform plan -replace=ADDR` reaches the provider as an
// ordinary update request: the plan protocol carries no forced-replacement
// field, and core resolves -replace only after the provider has answered and
// its diagnostics have been collected. So this warning is emitted unchanged on
// a plan the CLI then prints as `-/+`, and a text asserting an in-place update
// would contradict the line above it.
//
// What is true under either action is the part the user actually needs, so the
// detail leads with it and qualifies only the mechanism that follows.
const (
	rebuildWarningSummary = "Applying this change erases the server's disks"

	rebuildWarningDetail = "Changing this attribute is applied by rebuilding the server: the " +
		"operating system is reinstalled and everything on the disks is lost. Terraform " +
		"carries that out either as an in-place update or as a replacement, and the disks " +
		"go either way.\n\nIf the plan above shows an ordinary in-place update, that is the " +
		"whole of it: the server keeps its identity, addresses and hardware, and nothing in " +
		"the plan output says the disks are about to go. If the plan above shows the resource " +
		"being replaced, more than the disks is at stake - the server is destroyed and built " +
		"again, and an address it was given rather than asked for does not come back.\n\n" +
		"If that is what you want, apply as normal. If it is not, revert this attribute to " +
		"its previous value - and if the plan shows a replacement, reverting is not enough " +
		"on its own, because the replacement has a separate cause: a -replace argument, a " +
		"replace_triggered_by, or another attribute that forces it. Re-run the plan and read " +
		"it again before applying."
)

// rebuildDecision is the answer rebuildsOnApply gives about a plan. It has
// three values rather than two because the two callers need different things
// from the same question: ModifyPlan can honestly say "it might", while Update
// has to either rebuild or not.
type rebuildDecision int

const (
	// rebuildNo means nothing that is applied by a rebuild is changing.
	rebuildNo rebuildDecision = iota

	// rebuildYes means image_id or user_data is changing to a known value, so
	// the apply reaches RebuildAndPoll.
	rebuildYes

	// rebuildUnknown means image_id or user_data is planned as unknown - it
	// depends on something Terraform has not resolved yet - so whether the apply
	// reaches RebuildAndPoll is not decided until the value is.
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

// rebuildsOnApply answers whether this plan reaches RebuildAndPoll, which
// reinstalls the operating system and erases the disks. The answer has three
// values, not two, so it is a rebuildDecision rather than a bool.
//
// Update calls it to decide whether to rebuild, and ModifyPlan calls it to
// decide whether to warn. That is the point of it being one function: the two
// answered the same question separately once, and a plan that warns about a
// rebuild it does not perform - or performs one it did not warn about - is
// worse than either behaviour on its own.
//
// An unknown planned value is reported as rebuildUnknown rather than folded
// into either of the other two answers. It is not known yet whether it will
// differ once resolved, so at plan time the honest answer to "will this erase
// the disks" is "it might", and ModifyPlan warns on it exactly as it warns on
// rebuildYes - silence would be the wrong way to express a maybe. Update cannot
// act on a maybe, though: Terraform resolves every planned value before it calls
// Update, so an unknown arriving there is a contradiction rather than a change,
// and rebuilding on it would reinstall the operating system while sending the
// API an empty value. Keeping the two apart is what lets one function serve both
// without either caller having to lie.
func rebuildsOnApply(plan, state CloudBaremetalServerModel, adopting bool) rebuildDecision {
	return leastResolved(
		rebuildDecisionFor(plan.ImageID, state.ImageID, adopting),
		rebuildDecisionFor(plan.UserData, state.UserData, adopting),
	)
}

// rebuildDecisionFor answers the same question for a single attribute.
//
// A null planned value is never a rebuild: removing image_id or user_data from
// the configuration is not something the API can act on, so nothing is sent. An
// unknown one is a change whose size is not known yet, which is why it cannot be
// compared: types.String.ValueString reports "" for an unknown value, so a naive
// comparison would read an unknown planned value against a null or empty prior
// as "no change" and neither warn nor rebuild.
func rebuildDecisionFor(plan, state types.String, adopting bool) rebuildDecision {
	if plan.IsNull() {
		return rebuildNo
	}

	changed := plan.IsUnknown() || plan.ValueString() != state.ValueString()

	// An unknown value gets no special treatment here on purpose: adopting an
	// imported server suppresses a change against a null prior value whether or
	// not the new value has been resolved yet, because in neither case is the
	// difference something the user asked to have applied.
	if !changeSurvivesAdoption(changed, adopting, state.IsNull()) {
		return rebuildNo
	}

	if plan.IsUnknown() {
		return rebuildUnknown
	}

	return rebuildYes
}

// leastResolved combines the per-attribute answers into one, and an unresolved
// attribute outranks a definite one.
//
// That ordering looks backwards until you ask what each caller does with the
// answer. ModifyPlan warns on rebuildYes and rebuildUnknown alike, so it cannot
// tell the two apart and the ordering does not reach the user. Update can tell
// them apart, and there the question is not "is a rebuild coming" but "may I act
// on this plan at all" - and it may not, if any single attribute is still
// unresolved. Reporting rebuildYes because image_id is definite would let an
// unknown user_data through the guard beside it and be sent to the API as an
// empty string, which is the exact harm the guard exists to prevent.
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

// rebuildAtUnknownValueError is raised when an unknown image_id or user_data
// reaches Update, which should not be reachable: Terraform resolves planned
// values before it applies them. Rebuilding anyway would reinstall the operating
// system and send the API an empty value, and skipping the rebuild silently
// would break the promise that the plan-time warning and the apply agree, so the
// apply stops instead.
func rebuildAtUnknownValueError() diag.Diagnostic {
	return diag.NewErrorDiagnostic(
		"Unknown image_id or user_data at apply time",
		"Terraform resolves every planned value before applying it, so by the time the "+
			"server is updated image_id and user_data are either known or absent. One of them "+
			"arrived unknown instead.\n\nRebuilding on an unresolved value would reinstall the "+
			"operating system and send the API an empty value, so nothing was done. This is a "+
			"bug in the provider; please report it.",
	)
}

// rebuildWarning tells the user that applying this plan erases the server's
// disks, without claiming which action Terraform planned. See the const block.
func rebuildWarning() diag.Diagnostic {
	return diag.NewWarningDiagnostic(rebuildWarningSummary, rebuildWarningDetail)
}

// changeSurvivesAdoption reports whether a detected change is a real one the
// user asked for, rather than an artefact of adopting an imported server.
//
// It exists for image_id and user_data, the two attributes whose updates go
// through RebuildAndPoll. Both are no_refresh, so after an import they are null
// in prior state while set in config, which the comparison in
// rebuildDecisionFor reads as a change. Acting on that would reinstall the
// operating system of a server the user only just imported - and the plan would
// show a benign in-place update, so nothing would warn them first.
//
// Only the combination of an open adoption window AND a null prior value is
// suppressed. A change away from a known prior value is always real, and once
// the marker is cleared every change is real again.
func changeSurvivesAdoption(changed, adopting, priorValueNull bool) bool {
	return changed && !(adopting && priorValueNull)
}
