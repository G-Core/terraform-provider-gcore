// This file wraps the shared planmodifiers.*RequiresReplaceUnlessAdopting
// modifiers so that removing "payload" or "payload_wo_version" from
// configuration records null in state in place, instead of forcing
// replacement.
//
// The secret API has no update endpoint and never returns the payload once
// set, so removing these two arguments cannot be sent anywhere: the only way
// to honour it is to drop the values from state without touching the live
// secret. Forcing replacement here would destroy a secret whose material is
// no longer configured, for no operational reason - the secret still exists
// and still works.
//
// Adding or changing payload/payload_wo_version still replaces, exactly as
// the wrapped modifier says, unless the import-adoption marker is pending.
// After a removal no marker exists, so re-adding them later goes through the
// normal replace-unless-adopting path (a real change, not an adoption).
package cloud_secret

import (
	"context"
	"time"

	"github.com/G-Core/terraform-provider-gcore/internal/planmodifiers"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// computedAttributesWithoutModifiers lists this resource's Computed
// attributes that have no plan modifier of their own. Keep this list in sync
// with schema.go: adding UseStateForUnknown to one of them there makes its
// entry here redundant, not wrong, but removing it from schema.go without
// adding it here would let it go unknown again.
var computedAttributesWithoutModifiers = []string{
	"algorithm", "bit_length", "created", "mode", "secret_type", "status", "content_types",
}

// carryComputedAttributesForward preserves computedAttributesWithoutModifiers
// across any plan that reaches Update, so a downstream reference to one of
// them (for example `count = gcore_cloud_secret.s.status == "ACTIVE" ? 1 :
// 0`, or passing `.content_types`/`.created` into another resource's
// RequiresReplace attribute) does not see a plan-time unknown just because
// something else in this resource changed.
//
// The framework marks every Computed attribute whose config value is null as
// unknown whenever the planned state differs at all from the prior state
// (see MarkComputedNilsAsUnknown in the framework's PlanResourceChange),
// regardless of which attribute actually changed. None of
// computedAttributesWithoutModifiers has a plan modifier of its own -
// UseStateForUnknown on them in schema.go would mean more generated-file
// hunks - so without this helper they show "(known after apply)" on every
// Update-shaped plan, not only the payload-removal one. Update (resource.go)
// already copies them back from state before writing the final state, but
// that runs at apply time, too late for anything reading them at plan time.
// This calls that a documented tradeoff worth closing at plan time instead.
//
// This never reaches a replacement plan. For every replacement, Terraform
// core re-plans the create half again with a NULL prior state and uses that
// call's PlannedState instead (see payload_rules.go's requirePayloadOnCreate
// comment for the citation): req.State.Raw is null on that call, so this
// returns immediately and never overwrites the unknown values a real create
// needs.
func carryComputedAttributesForward(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}

	for _, name := range computedAttributesWithoutModifiers {
		attrPath := path.Root(name)

		var planned attr.Value
		resp.Diagnostics.Append(resp.Plan.GetAttribute(ctx, attrPath, &planned)...)
		if resp.Diagnostics.HasError() {
			return
		}
		if planned == nil || !planned.IsUnknown() {
			continue
		}

		var prior attr.Value
		resp.Diagnostics.Append(req.State.GetAttribute(ctx, attrPath, &prior)...)
		if resp.Diagnostics.HasError() {
			return
		}

		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, attrPath, prior)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
}

// payloadRemovedInPlace reports whether a plan drops a previously-set value
// to null: both state and plan are non-null objects (a real resource, not
// create/destroy), the prior value was set, and the planned value is null.
//
// An unknown planned value is not null, so it never counts as removal here:
// PlanModifyObject/PlanModifyInt64 fall through to the wrapped modifier for
// it, same as for any other change.
func payloadRemovedInPlace(stateRawNull, planRawNull, priorValueNull, plannedValueNull bool) bool {
	return !stateRawNull && !planRawNull && !priorValueNull && plannedValueNull
}

// requiresReplaceUnlessAdoptingOrRemovedDescription is shared by both typed
// wrappers below so their documented behaviour cannot drift apart.
const requiresReplaceUnlessAdoptingOrRemovedDescription = "Requires replacement when the value is set or changed, " +
	"except while adopting a value right after terraform import; removing it from configuration keeps the resource."

// objectReplaceUnlessAdoptingOrRemoved returns a plan modifier that behaves
// like planmodifiers.ObjectRequiresReplaceUnlessAdopting, except that
// removing the value from configuration never forces replacement.
func objectReplaceUnlessAdoptingOrRemoved(key string) planmodifier.Object {
	return objectReplaceUnlessAdoptingOrRemovedModifier{inner: planmodifiers.ObjectRequiresReplaceUnlessAdopting(key)}
}

type objectReplaceUnlessAdoptingOrRemovedModifier struct {
	inner planmodifier.Object
}

func (m objectReplaceUnlessAdoptingOrRemovedModifier) Description(_ context.Context) string {
	return requiresReplaceUnlessAdoptingOrRemovedDescription
}

func (m objectReplaceUnlessAdoptingOrRemovedModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m objectReplaceUnlessAdoptingOrRemovedModifier) PlanModifyObject(ctx context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
	if payloadRemovedInPlace(req.State.Raw.IsNull(), req.Plan.Raw.IsNull(), req.StateValue.IsNull(), req.PlanValue.IsNull()) {
		return
	}
	m.inner.PlanModifyObject(ctx, req, resp)
}

// int64ReplaceUnlessAdoptingOrRemoved is the Int64 counterpart of
// objectReplaceUnlessAdoptingOrRemoved, for payload_wo_version.
func int64ReplaceUnlessAdoptingOrRemoved(key string) planmodifier.Int64 {
	return int64ReplaceUnlessAdoptingOrRemovedModifier{inner: planmodifiers.Int64RequiresReplaceUnlessAdopting(key)}
}

type int64ReplaceUnlessAdoptingOrRemovedModifier struct {
	inner planmodifier.Int64
}

func (m int64ReplaceUnlessAdoptingOrRemovedModifier) Description(_ context.Context) string {
	return requiresReplaceUnlessAdoptingOrRemovedDescription
}

func (m int64ReplaceUnlessAdoptingOrRemovedModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m int64ReplaceUnlessAdoptingOrRemovedModifier) PlanModifyInt64(ctx context.Context, req planmodifier.Int64Request, resp *planmodifier.Int64Response) {
	if payloadRemovedInPlace(req.State.Raw.IsNull(), req.Plan.Raw.IsNull(), req.StateValue.IsNull(), req.PlanValue.IsNull()) {
		return
	}
	m.inner.PlanModifyInt64(ctx, req, resp)
}

// expirationRequiresReplace replaces the secret only when the configured
// expiration is a different instant from prior state. The API echoes
// expiration in UTC with a Z suffix, so after `terraform import` state holds
// "...Z" while a config written with an offset ("+00:00", "+03:00") names the
// same instant; plain RequiresReplace would replace an unchanged secret. A
// spelling-only change plans in place instead, Update (which sends nothing)
// records the configured spelling, and Read keeps it (see sameInstant).
func expirationRequiresReplace() planmodifier.String {
	return stringplanmodifier.RequiresReplaceIf(
		func(_ context.Context, req planmodifier.StringRequest, resp *stringplanmodifier.RequiresReplaceIfFuncResponse) {
			resp.RequiresReplace = !sameInstant(req.StateValue, req.PlanValue)
		},
		"Requires replacement when the configured expiration is a different instant.",
		"Requires replacement when the configured expiration is a different instant.",
	)
}

// sameInstant reports whether two known, non-null RFC 3339 values name the same
// instant to the second, whatever their offsets. Seconds, because apijson
// sends and reads expiration with time.RFC3339, which drops any fraction: the
// API only ever stores whole seconds. timetypes.RFC3339 semantic equality is
// not enough here: it compares the offset too, so "+03:00" never equals the
// API's "Z".
func sameInstant(a, b basetypes.StringValue) bool {
	if a.IsNull() || a.IsUnknown() || b.IsNull() || b.IsUnknown() {
		return false
	}
	at, aErr := time.Parse(time.RFC3339, a.ValueString())
	bt, bErr := time.Parse(time.RFC3339, b.ValueString())
	return aErr == nil && bErr == nil && at.Truncate(time.Second).Equal(bt.Truncate(time.Second))
}
