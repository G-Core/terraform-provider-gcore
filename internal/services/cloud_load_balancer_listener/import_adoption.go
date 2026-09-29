package cloud_load_balancer_listener

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// importAdoptionPrivateKey marks a listener whose create-only attributes still
// have to be adopted from config into state after a `terraform import`.
//
// It is written by ImportState, read by the plan modifiers on those attributes
// and by ModifyPlan, advanced through its lifecycle by Read, and cleared by
// Update once the adoption has happened.
//
// The marker exists because "prior state is null" does NOT by itself mean "just
// imported". It is equally true of a listener created with the attribute
// omitted, and a user may set such an attribute on day two. Treating that as an
// adoption would plan a clean in-place update, send nothing the API acts on,
// and write the configured value to state - leaving state describing a
// listener that does not exist. Keying on an explicit import marker keeps that
// case on the replacement path, which is the only path that actually applies
// the value.
//
// Neither adopted attribute can be corroborated against the API: the listener
// GET never echoes default_pool_id or insert_x_forwarded, and nothing else in
// the response stands in for them.
//
// A marker that only cleared in Update could outlive its import indefinitely: a
// config setting neither attribute produces a no-op plan, so no Update ever
// runs to retire it, and a day-two attribute set weeks later would be silently
// adopted instead of planning the replacement that would actually apply it.
// That hole is closed by the phased lifecycle below, which Read advances until
// the marker retires itself. What remains is bounded and mostly loud:
//
//   - the first apply after import is still an adoption window - at that point
//     "finishing the import" and "day-two set" are indistinguishable in
//     everything a provider can observe, so the window is the contract, and the
//     plan-time warning in ModifyPlan names every attribute being adopted;
//   - `import` blocks get one extra apply of window, because their plan runs in
//     the same walk as the import itself;
//   - `-refresh=false` skips Read and therefore does not advance the marker,
//     leaving the window open until the first refreshed operation;
//   - a refresh-only operation between the import and the adopting apply
//     consumes the window, and the adopting apply then plans a replacement;
//   - a failed dependency consumes it too, and this is an accepted, documented
//     limitation rather than a solved one. If the adopting apply also changes
//     a resource the listener depends on (typically its load balancer, through
//     load_balancer_id) and that change fails, Terraform skips this listener's
//     Update but has already refreshed it, so the phase advanced with nothing
//     adopted. A retry spends the rest of the window and plans a replacement.
//     No plan modifier can see another resource's apply outcome, so this cannot
//     be prevented here; it is described in the adoption warning and in the
//     resource's import docs.
//
// Recovery from a consumed window is `terraform state rm` followed by
// `terraform import` - import refuses an address that is still in state, so
// the entry has to go first. `-refresh=false` only helps while the saved state
// still holds the marker; once a refreshed apply has deleted it, nothing
// brings it back. The durable fix is the API returning both attributes on
// GET, which would make this whole file unnecessary.
//
// Mirrors cloud_instance and cloud_baremetal_server.
const importAdoptionPrivateKey = "load_balancer_listener_import_adoption_pending"

// The marker is a three-phase lifecycle advanced by Read, the only provider
// method whose private-state writes Terraform persists in every `terraform
// apply` - including a no-op one, which keeps the refreshed prior state.
// ImportState writes "pending"; the Read Terraform runs immediately after
// import advances it to "armed"; the refresh at the start of the next apply
// advances it to "spent"; the Read after that deletes it. "spent" is
// deliberately NON-EMPTY: every reader tests only presence, so the plan of the
// invocation that spends the marker still sees "adopting" - that is what lets
// the first apply after import adopt. Do not collapse spent into deletion:
// deleting on armed breaks adoption on the `terraform import` path.
//
// Net effect: the marker can influence plans up to and including the first
// apply after import (one apply later for `import` blocks, whose plan runs in
// the same walk as the import), and is mechanically dead afterwards. Update
// retires it early on any successful apply, and re-arms it on a failed one so
// a transient API error does not turn the retry into a destroy+recreate.
var (
	importAdoptionValuePending = []byte(`{"phase":"pending"}`)
	importAdoptionValueArmed   = []byte(`{"phase":"armed"}`)
	importAdoptionValueSpent   = []byte(`{"phase":"spent"}`)
)

const (
	importAdoptionPhasePending = "pending"
	importAdoptionPhaseArmed   = "armed"
	importAdoptionPhaseSpent   = "spent"
)

// privateStateReader and privateStateWriter are the subsets of the framework's
// private state data used here. The concrete type lives in an internal framework
// package, so it is reached through interfaces rather than named directly.
type privateStateReader interface {
	GetKey(ctx context.Context, key string) ([]byte, diag.Diagnostics)
}

type privateStateWriter interface {
	SetKey(ctx context.Context, key string, value []byte) diag.Diagnostics
}

type privateStateReadWriter interface {
	privateStateReader
	privateStateWriter
}

// markImportAdoptionPending records that the create-only attributes of a
// freshly imported listener still have to be adopted into state.
func markImportAdoptionPending(ctx context.Context, private privateStateWriter) diag.Diagnostics {
	if private == nil {
		return nil
	}

	return private.SetKey(ctx, importAdoptionPrivateKey, importAdoptionValuePending)
}

// importAdoptionPhaseOf parses a marker value. Absent -> "". A non-empty value
// that does not parse is treated as armed: one apply of window, then
// retirement.
func importAdoptionPhaseOf(value []byte) string {
	if len(value) == 0 {
		return ""
	}
	var marker struct {
		Phase string `json:"phase"`
	}
	if err := json.Unmarshal(value, &marker); err != nil || marker.Phase == "" {
		return importAdoptionPhaseArmed
	}
	switch marker.Phase {
	case importAdoptionPhasePending, importAdoptionPhaseArmed, importAdoptionPhaseSpent:
		return marker.Phase
	default:
		return importAdoptionPhaseArmed
	}
}

// advanceImportAdoption moves the marker one step along its lifecycle:
// pending->armed, armed->spent, spent->deleted. Call exactly once per Read, on
// the success path, before the final state write. No-op when no marker exists.
func advanceImportAdoption(ctx context.Context, private privateStateReadWriter) diag.Diagnostics {
	if private == nil {
		return nil
	}
	value, diags := private.GetKey(ctx, importAdoptionPrivateKey)
	if diags.HasError() || len(value) == 0 {
		return diags
	}
	switch importAdoptionPhaseOf(value) {
	case importAdoptionPhasePending:
		diags.Append(private.SetKey(ctx, importAdoptionPrivateKey, importAdoptionValueArmed)...)
	case importAdoptionPhaseArmed:
		diags.Append(private.SetKey(ctx, importAdoptionPrivateKey, importAdoptionValueSpent)...)
	case importAdoptionPhaseSpent:
		diags.Append(private.SetKey(ctx, importAdoptionPrivateKey, nil)...)
	}
	return diags
}

// rearmImportAdoption restores a full apply of adoption window. Called from
// Update's error path while adopting, so a transient API failure does not let
// the next refresh retire the marker and turn the retry into a replacement.
func rearmImportAdoption(ctx context.Context, private privateStateWriter) diag.Diagnostics {
	if private == nil {
		return nil
	}
	return private.SetKey(ctx, importAdoptionPrivateKey, importAdoptionValueArmed)
}

// importAdoptionPending reports whether the listener is in the window between
// an import and the apply that adopts its create-only attributes.
func importAdoptionPending(ctx context.Context, private privateStateReader) (bool, diag.Diagnostics) {
	if private == nil {
		return false, nil
	}

	value, diags := private.GetKey(ctx, importAdoptionPrivateKey)

	return len(value) > 0, diags
}

// clearImportAdoptionPending closes the adoption window. Setting an empty
// value deletes the key.
func clearImportAdoptionPending(ctx context.Context, private privateStateWriter) diag.Diagnostics {
	if private == nil {
		return nil
	}

	return private.SetKey(ctx, importAdoptionPrivateKey, nil)
}

// createOnlyAttributes maps each create-only attribute's Terraform name to a
// predicate reporting whether the model has no value for it. These are the
// attributes carrying an adopting plan modifier in schema.go: default_pool_id
// and insert_x_forwarded, both tagged no_refresh in model.go.
//
// The names must match the `tfsdk:` tags in model.go, since they are what the
// warning shows the user.
var createOnlyAttributes = map[string]func(CloudLoadBalancerListenerModel) bool{
	"default_pool_id":    func(m CloudLoadBalancerListenerModel) bool { return m.DefaultPoolID.IsNull() },
	"insert_x_forwarded": func(m CloudLoadBalancerListenerModel) bool { return m.InsertXForwarded.IsNull() },
}

// adoptedCreateOnlyAttributes returns the create-only attributes this plan
// would record into state without sending them to the API: configured now,
// absent from prior state. The result is sorted so the warning text is
// stable.
func adoptedCreateOnlyAttributes(plan, state CloudLoadBalancerListenerModel) []string {
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

// patchBaselineFor returns the prior state that Update diffs the PATCH body
// against. While adopting, it carries the planned create-only values, so the
// diff leaves them out of a request to an endpoint that never accepts them.
func patchBaselineFor(state, plan CloudLoadBalancerListenerModel, adopting bool) CloudLoadBalancerListenerModel {
	if adopting {
		state.DefaultPoolID = plan.DefaultPoolID
		state.InsertXForwarded = plan.InsertXForwarded
	}

	return state
}

// replacementForcingAttributes maps every attribute whose plan modifier in
// schema.go can force a destroy+recreate to the rule that modifier applies.
// TestReplacementForcingAttributesMatchSchema fails if codegen adds, drops or
// re-modifies one.
var replacementForcingAttributes = map[string]replacementRule{
	"default_pool_id":    replacesUnlessAdopting,
	"insert_x_forwarded": replacesUnlessAdopting,
	"load_balancer_id":   replacesOnAnyChange,
	"project_id":         replacesWhenConfigured,
	"protocol":           replacesOnAnyChange,
	"protocol_port":      replacesOnAnyChange,
	"region_id":          replacesWhenConfigured,
}

// replacementRule is the shape of every rule above. adopting matters to one of
// them, but the table is uniform so a rule cannot be bound to a signature that
// silently drops it.
type replacementRule func(plan, state, config attr.Value, adopting bool) bool

// replacesUnlessAdopting is the rule of the *RequiresReplaceUnlessAdopting
// family: any change replaces, except a null prior value during the adoption
// that follows an import.
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
// destroy+recreate. Then nothing is adopted: Terraform runs Create and sends
// every configured value to the API, so the adoption warning must stay quiet.
//
// ModifyPlan cannot read the framework's replacement decision - the plan
// modifiers that make it run earlier, against a response of their own - so
// the same question is put to the values they compared, under the same rules.
//
// It covers replacement forced by an attribute change only. `-replace=ADDR`
// and `replace_triggered_by` are decided by Terraform core after the provider
// answers and never reach the plan request, so the warning text still opens
// with a caveat for that case. A tainted resource needs no rule: it arrives
// with a null prior state, which ModifyPlan returns on already.
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

// importAdoptionWarning names the attributes being adopted and the one way
// this can go wrong, which is a configured value that does not match the live
// listener. It also tells the user the window is closing, so a later
// attribute set planning a replacement does not come as a surprise.
//
// ModifyPlan does not emit it when an attribute change forces a replacement
// (see planForcesReplacement). The opening caveat covers what it cannot see:
// a `-replace` on the command line or `replace_triggered_by`, where Terraform
// runs Create and every configured value reaches the API.
func importAdoptionWarning(adopted []string) diag.Diagnostic {
	return diag.NewWarningDiagnostic(
		"Adopting create-only attributes into state after import",
		fmt.Sprintf(
			"If the plan above shows this resource being replaced, none of the following applies: "+
				"replacing it creates the listener again and every configured value reaches the API.\n\n"+
				"Otherwise, Terraform is recording %s from configuration instead of replacing the "+
				"listener, because this listener was imported and the API does not return these "+
				"attributes. Nothing is sent to the API for them, so the configured values are trusted "+
				"as-is.\n\n"+
				"If a recorded value does not match the running listener, state will be wrong and no "+
				"future refresh can correct it - these attributes are never returned by the API.\n\n"+
				"Check them against the listener before applying. To force the values to actually "+
				"take effect instead, run `terraform apply -replace=` for this resource. "+
				"This adoption window closes after this apply; afterwards, setting these attributes "+
				"plans a replacement.\n\n"+
				"Apply this on its own if you can. If the same apply also changes a resource this "+
				"listener depends on, such as its load balancer, and that change fails, Terraform "+
				"skips the listener but the window still closes, and the next plan replaces the "+
				"listener. To recover from that, run `terraform state rm` for this resource and "+
				"import it again.",
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
