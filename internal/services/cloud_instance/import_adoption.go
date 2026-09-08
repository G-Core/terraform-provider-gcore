package cloud_instance

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/G-Core/gcore-go"
	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/packages/param"
	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/G-Core/gcore-go/option"
	"github.com/G-Core/terraform-provider-gcore/internal/logging"
)

// importAdoptionPrivateKey marks an instance whose create-only attributes still
// have to be adopted from config into state after a `terraform import`.
//
// It is written by ImportState, read by the plan modifiers on those attributes
// and by ModifyPlan, advanced through its lifecycle by Read, and cleared by
// Update once the adoption has happened.
//
// The marker exists because "prior state is null" does NOT by itself mean "just
// imported". It is equally true of an instance created with the attribute
// omitted, and a user may set such an attribute on day two. Treating that as an
// adoption would plan a clean in-place update, send nothing the API acts on, and
// write the configured value to state - leaving state describing an instance
// that does not exist. Keying on an explicit import marker keeps that case on
// the replacement path, which is the only path that actually applies the value.
//
// Unlike cloud_load_balancer, most adoptions here cannot be corroborated
// against the API: the instance GET never echoes user_data, username or a
// servergroup placement, and nothing else in the response stands in for them.
// ssh_key_name is the one exception - the response does carry it, so ImportState
// resolves it and adoption only ever covers the residual case where that lookup
// failed.
//
// A marker used to be able to outlive its import indefinitely: a config setting
// none of these attributes produces no-op plans, so no Update ever ran to retire
// it, and a day-two attribute set weeks later was silently adopted instead of
// planning the replacement that would actually apply it. That hole is closed by
// the phased lifecycle below, which Read advances until the marker retires
// itself. What remains is bounded and mostly loud:
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
//     consumes the window, and the adopting apply then plans a replacement -
//     loud, and recoverable by running that apply with `-refresh=false` or by
//     re-importing.
const importAdoptionPrivateKey = "instance_import_adoption_pending"

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

// markImportAdoptionPending records that the create-only attributes of a freshly
// imported instance still have to be adopted into state.
func markImportAdoptionPending(ctx context.Context, private privateStateWriter) diag.Diagnostics {
	if private == nil {
		return nil
	}

	return private.SetKey(ctx, importAdoptionPrivateKey, importAdoptionValuePending)
}

// importAdoptionPhaseOf parses a marker value. Absent -> "". A non-empty value
// that does not parse (including the pre-phase {"pending":true} written by
// earlier releases of this pattern) is treated as armed: one apply of window,
// then retirement.
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

// importAdoptionPending reports whether the instance is in the window between an
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

// instanceInPlacementGroup reports whether the instance is already a member of
// the placement group, which is what tells a bookkeeping adoption apart from a
// real placement change. servergroup_id is create-only and never echoed by the
// instance GET, so this is the only way to establish it.
func instanceInPlacementGroup(ctx context.Context, client *gcore.Client, servergroupID, instanceID string, data *CloudInstanceModel) (bool, error) {
	params := cloud.PlacementGroupGetParams{}

	if !data.ProjectID.IsNull() {
		params.ProjectID = param.NewOpt(data.ProjectID.ValueInt64())
	}

	if !data.RegionID.IsNull() {
		params.RegionID = param.NewOpt(data.RegionID.ValueInt64())
	}

	group, err := client.Cloud.PlacementGroups.Get(
		ctx,
		servergroupID,
		params,
		option.WithMiddleware(logging.Middleware(ctx)),
	)
	if err != nil {
		return false, err
	}

	for _, instance := range group.Instances {
		if instance.InstanceID == instanceID {
			return true, nil
		}
	}

	return false, nil
}

// createOnlyAttributes maps each create-only attribute's Terraform name to a
// predicate reporting whether the model has no value for it. These are the
// attributes carrying an adopting plan modifier in schema.go: the six tagged
// no_refresh in model.go, plus ssh_key_name, which ImportState back-fills from
// the API response and which is therefore rarely adopted in practice - the
// predicate simply skips it once it is populated. password_wo is write-only and
// password_wo_version keeps unconditional replacement, so neither appears here.
//
// The names must match the `tfsdk:` tags in model.go, since they are what the
// warning shows the user.
var createOnlyAttributes = map[string]func(CloudInstanceModel) bool{
	"allow_app_ports": func(m CloudInstanceModel) bool { return m.AllowAppPorts.IsNull() },
	"name_template":   func(m CloudInstanceModel) bool { return m.NameTemplate.IsNull() },
	"servergroup_id":  func(m CloudInstanceModel) bool { return m.ServergroupID.IsNull() },
	"ssh_key_name":    func(m CloudInstanceModel) bool { return m.SSHKeyName.IsNull() },
	"user_data":       func(m CloudInstanceModel) bool { return m.UserData.IsNull() },
	"username":        func(m CloudInstanceModel) bool { return m.Username.IsNull() },
	"configuration":   func(m CloudInstanceModel) bool { return m.Configuration == nil || *m.Configuration == nil },
}

// adoptedCreateOnlyAttributes returns the create-only attributes this plan would
// record into state without sending them to the API: configured now, absent from
// prior state. The result is sorted so the warning text is stable.
func adoptedCreateOnlyAttributes(plan, state CloudInstanceModel) []string {
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

// importAdoptionWarning names the attributes being adopted and the one way this
// can go wrong, which is a configured value that does not match the running
// instance. It also tells the user the window is closing, so a later attribute
// set planning a replacement does not come as a surprise.
//
// Most adopted attributes are only written to state. "servergroup_id" is the
// exception - Update really does place the instance in the configured group -
// so the warning must not promise that nothing reaches the API when it is in
// the set.
func importAdoptionWarning(adopted []string) diag.Diagnostic {
	detail := fmt.Sprintf(
		"Terraform is recording %s from configuration instead of replacing the instance, "+
			"because this instance was imported and the API does not return these attributes.\n\n",
		strings.Join(quoteAll(adopted), ", "),
	)

	if slices.Contains(adopted, "servergroup_id") {
		detail += "These values are trusted as-is rather than sent to the API, with one exception: " +
			"for \"servergroup_id\" Terraform checks the placement group's members and adds the " +
			"instance to that group if it is not already in it, so applying can change the " +
			"instance's actual placement.\n\n"
	} else {
		detail += "Nothing is sent to the API for them, so the configured values are trusted as-is.\n\n"
	}

	detail += "If a recorded value does not match the running instance, state will be wrong and no " +
		"future refresh can correct it - these attributes are never returned by the API.\n\n" +
		"Check them against the instance before applying. To force the values to actually " +
		"take effect instead, run `terraform apply -replace=` for this resource. " +
		"This adoption window closes after this apply; afterwards, setting these attributes plans a replacement."

	return diag.NewWarningDiagnostic(
		"Adopting create-only attributes into state after import",
		detail,
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
