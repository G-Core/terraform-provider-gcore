package planmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// StringRequiresReplaceUnlessAdopting is for create-only attributes the API
// does not return on GET (`no_refresh`).
//
// Such an attribute is null in state after `terraform import`, so a configured
// value looks like a change and an unconditional RequiresReplace() would plan a
// destroy+recreate of live infrastructure. Skipping replacement whenever prior
// state is null is NOT a safe generalisation, though: prior state is equally
// null for a resource created with the attribute omitted, and silently adopting
// there would record a value the infrastructure does not have.
//
// The modifier therefore skips replacement only when the resource carries a
// private-state marker written at import time and cleared once the value has
// been adopted. Every other transition keeps the built-in semantics, including
// removal (config set -> null) and an unknown planned value.
//
// The resource owning the attribute defines the private-state key and is
// responsible for setting it in ImportState and clearing it in Update.
//
// requiresReplaceUnlessAdopting below is the type-agnostic half of the
// decision, so a typed counterpart for another attribute type is a thin
// wrapper around it.

// privateStateReader is the subset of the framework's private state data these
// modifiers need. The concrete type lives in an internal framework package, so
// it is accepted through an interface rather than named directly. Its GetKey is
// safe to call on a nil receiver.
type privateStateReader interface {
	GetKey(ctx context.Context, key string) ([]byte, diag.Diagnostics)
}

// adoptionPending reports whether the resource is in the window between an
// import and the apply that adopts its create-only attributes into state.
func adoptionPending(ctx context.Context, private privateStateReader, key string) (bool, diag.Diagnostics) {
	if private == nil {
		return false, nil
	}

	value, diags := private.GetKey(ctx, key)

	return len(value) > 0, diags
}

// requiresReplaceUnlessAdopting is the decision shared by every typed modifier
// in this family, kept as a pure function so the whole matrix can be unit
// tested: the framework request types carry private state as an internal type
// that cannot be constructed from outside the framework module.
//
// The only case that skips replacement on a null prior value is an adoption
// following an import. In particular a resource created with the attribute
// omitted also has a null prior value, and setting it later is a real change
// the update API cannot honour, so it must still force replacement.
func requiresReplaceUnlessAdopting(stateRawNull, planRawNull, valueUnchanged, priorValueNull, adopting bool) bool {
	switch {
	case stateRawNull: // resource creation
		return false
	case planRawNull: // resource destruction
		return false
	case valueUnchanged:
		return false
	case priorValueNull && adopting: // one-time adoption after import
		return false
	default:
		return true
	}
}

// requiresReplaceUnlessAdoptingDescription is shared by every typed modifier in
// this family so their documented behaviour cannot drift apart.
const requiresReplaceUnlessAdoptingDescription = "Requires replacement when the value changes, except while adopting a value into " +
	"state immediately after terraform import."

// StringRequiresReplaceUnlessAdopting returns a plan modifier that requires
// replacement when a create-only string attribute changes, except during the
// one-time adoption that follows an import. See the package-level notes above.
func StringRequiresReplaceUnlessAdopting(privateKey string) planmodifier.String {
	return stringRequiresReplaceUnlessAdoptingModifier{privateKey: privateKey}
}

type stringRequiresReplaceUnlessAdoptingModifier struct {
	privateKey string
}

func (m stringRequiresReplaceUnlessAdoptingModifier) Description(_ context.Context) string {
	return requiresReplaceUnlessAdoptingDescription
}

func (m stringRequiresReplaceUnlessAdoptingModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m stringRequiresReplaceUnlessAdoptingModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	pending, diags := adoptionPending(ctx, req.Private, m.privateKey)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.RequiresReplace = requiresReplaceUnlessAdopting(
		req.State.Raw.IsNull(),
		req.Plan.Raw.IsNull(),
		req.PlanValue.Equal(req.StateValue),
		req.StateValue.IsNull(),
		pending,
	)
}

// BoolRequiresReplaceUnlessAdopting returns the Bool counterpart of
// StringRequiresReplaceUnlessAdopting. See the package-level notes above.
func BoolRequiresReplaceUnlessAdopting(privateKey string) planmodifier.Bool {
	return boolRequiresReplaceUnlessAdoptingModifier{privateKey: privateKey}
}

type boolRequiresReplaceUnlessAdoptingModifier struct {
	privateKey string
}

func (m boolRequiresReplaceUnlessAdoptingModifier) Description(_ context.Context) string {
	return requiresReplaceUnlessAdoptingDescription
}

func (m boolRequiresReplaceUnlessAdoptingModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m boolRequiresReplaceUnlessAdoptingModifier) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	pending, diags := adoptionPending(ctx, req.Private, m.privateKey)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.RequiresReplace = requiresReplaceUnlessAdopting(
		req.State.Raw.IsNull(),
		req.Plan.Raw.IsNull(),
		req.PlanValue.Equal(req.StateValue),
		req.StateValue.IsNull(),
		pending,
	)
}

// Int64RequiresReplaceUnlessAdopting returns the Int64 counterpart of
// StringRequiresReplaceUnlessAdopting. See the package-level notes above.
func Int64RequiresReplaceUnlessAdopting(privateKey string) planmodifier.Int64 {
	return int64RequiresReplaceUnlessAdoptingModifier{privateKey: privateKey}
}

type int64RequiresReplaceUnlessAdoptingModifier struct {
	privateKey string
}

func (m int64RequiresReplaceUnlessAdoptingModifier) Description(_ context.Context) string {
	return requiresReplaceUnlessAdoptingDescription
}

func (m int64RequiresReplaceUnlessAdoptingModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m int64RequiresReplaceUnlessAdoptingModifier) PlanModifyInt64(ctx context.Context, req planmodifier.Int64Request, resp *planmodifier.Int64Response) {
	pending, diags := adoptionPending(ctx, req.Private, m.privateKey)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.RequiresReplace = requiresReplaceUnlessAdopting(
		req.State.Raw.IsNull(),
		req.Plan.Raw.IsNull(),
		req.PlanValue.Equal(req.StateValue),
		req.StateValue.IsNull(),
		pending,
	)
}

// MapRequiresReplaceUnlessAdopting returns the Map counterpart of
// StringRequiresReplaceUnlessAdopting. See the package-level notes above.
func MapRequiresReplaceUnlessAdopting(privateKey string) planmodifier.Map {
	return mapRequiresReplaceUnlessAdoptingModifier{privateKey: privateKey}
}

type mapRequiresReplaceUnlessAdoptingModifier struct {
	privateKey string
}

func (m mapRequiresReplaceUnlessAdoptingModifier) Description(_ context.Context) string {
	return requiresReplaceUnlessAdoptingDescription
}

func (m mapRequiresReplaceUnlessAdoptingModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m mapRequiresReplaceUnlessAdoptingModifier) PlanModifyMap(ctx context.Context, req planmodifier.MapRequest, resp *planmodifier.MapResponse) {
	pending, diags := adoptionPending(ctx, req.Private, m.privateKey)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.RequiresReplace = requiresReplaceUnlessAdopting(
		req.State.Raw.IsNull(),
		req.Plan.Raw.IsNull(),
		req.PlanValue.Equal(req.StateValue),
		req.StateValue.IsNull(),
		pending,
	)
}

// ObjectRequiresReplaceUnlessAdopting returns the Object counterpart of
// StringRequiresReplaceUnlessAdopting. See the package-level notes above.
func ObjectRequiresReplaceUnlessAdopting(privateKey string) planmodifier.Object {
	return objectRequiresReplaceUnlessAdoptingModifier{privateKey: privateKey}
}

type objectRequiresReplaceUnlessAdoptingModifier struct {
	privateKey string
}

func (m objectRequiresReplaceUnlessAdoptingModifier) Description(_ context.Context) string {
	return requiresReplaceUnlessAdoptingDescription
}

func (m objectRequiresReplaceUnlessAdoptingModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m objectRequiresReplaceUnlessAdoptingModifier) PlanModifyObject(ctx context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
	pending, diags := adoptionPending(ctx, req.Private, m.privateKey)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.RequiresReplace = requiresReplaceUnlessAdopting(
		req.State.Raw.IsNull(),
		req.Plan.Raw.IsNull(),
		req.PlanValue.Equal(req.StateValue),
		req.StateValue.IsNull(),
		pending,
	)
}
