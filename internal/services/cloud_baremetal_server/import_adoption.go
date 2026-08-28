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
var replacementForcingAttributes = map[string]func(plan, state, config attr.Value) bool{
	"app_config":          replacesOverKnownPrior,
	"apptemplate_id":      replacesOverKnownPrior,
	"flavor":              replacesOverKnownPrior,
	"interfaces":          replacesOverKnownPrior,
	"name_template":       replacesOverKnownPrior,
	"password_wo_version": replacesOnAnyChange,
	"project_id":          replacesWhenConfigured,
	"region_id":           replacesWhenConfigured,
	"ssh_key_name":        replacesOverKnownPrior,
	"username":            replacesOverKnownPrior,
}

// replacesOverKnownPrior is the rule of ListRequiresReplaceIfNotNull and, while
// the marker is set, of the *RequiresReplaceUnlessAdopting family: a null prior
// value is either a creation or the adoption this warning is about.
func replacesOverKnownPrior(plan, state, _ attr.Value) bool {
	return !state.IsNull() && !plan.Equal(state)
}

// replacesWhenConfigured is the rule of RequiresReplaceIfConfigured.
func replacesWhenConfigured(plan, state, config attr.Value) bool {
	return !config.IsNull() && !plan.Equal(state)
}

// replacesOnAnyChange is the rule of RequiresReplace.
func replacesOnAnyChange(plan, state, _ attr.Value) bool {
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
func planForcesReplacement(ctx context.Context, plan tfsdk.Plan, state tfsdk.State, config tfsdk.Config) bool {
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

		if forcesReplacement(planValue, stateValue, configValue) {
			return true
		}
	}

	return false
}

// importAdoptionWarning names the attributes being adopted and the one way this
// can go wrong, which is a marker that outlived its import.
func importAdoptionWarning(adopted []string) diag.Diagnostic {
	return diag.NewWarningDiagnostic(
		"Adopting create-only attributes into state after import",
		fmt.Sprintf(
			"Terraform is recording %s from configuration instead of replacing the server, "+
				"because this server was imported and the API does not return these attributes.\n\n"+
				"Nothing is sent to the API for them, so the configured values are trusted as-is. "+
				"In particular no rebuild is performed for image_id or user_data during this "+
				"adoption, which is deliberate: rebuilding would reinstall the operating system "+
				"of a server that was only just imported.\n\n"+
				"If the configured values do not match the running server, state will be wrong and "+
				"no future refresh can correct it. Check them before applying. To force the values "+
				"to actually take effect, run `terraform apply -replace=` for this resource.",
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

// changeSurvivesAdoption reports whether a detected change is a real one the
// user asked for, rather than an artefact of adopting an imported server.
//
// It exists for image_id and user_data, the two attributes whose updates go
// through RebuildAndPoll. Both are no_refresh, so after an import they are null
// in prior state while set in config, which the plain comparison in Update reads
// as a change. Acting on that would reinstall the operating system of a server
// the user only just imported - and the plan would show a benign in-place
// update, so nothing would warn them first.
//
// Only the combination of an open adoption window AND a null prior value is
// suppressed. A change away from a known prior value is always real, and once
// the marker is cleared every change is real again.
func changeSurvivesAdoption(changed, adopting, priorValueNull bool) bool {
	return changed && !(adopting && priorValueNull)
}
