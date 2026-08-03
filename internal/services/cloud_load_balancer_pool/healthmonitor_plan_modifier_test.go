package cloud_load_balancer_pool

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// healthmonitorAttr returns the healthmonitor schema attribute and its object
// type, failing the test if the attribute stopped being a Computed+Optional
// single nested attribute — the shape the plan modifier below compensates for.
func healthmonitorAttr(t *testing.T, ctx context.Context) (schema.SingleNestedAttribute, types.ObjectType) {
	t.Helper()

	attrRaw, ok := ResourceSchema(ctx).Attributes["healthmonitor"]
	if !ok {
		t.Fatal("healthmonitor attribute missing from resource schema")
	}
	hm, ok := attrRaw.(schema.SingleNestedAttribute)
	if !ok {
		t.Fatalf("healthmonitor is %T, expected schema.SingleNestedAttribute", attrRaw)
	}
	if !hm.Computed || !hm.Optional {
		// If the attribute is no longer Computed+Optional the framework can no
		// longer inject the unknown, so the modifier may be obsolete —
		// re-evaluate rather than blindly re-adding it.
		t.Fatalf("healthmonitor configurability changed (Computed=%v Optional=%v); re-evaluate the unknown-guard plan modifier", hm.Computed, hm.Optional)
	}
	objType, ok := hm.GetType().(types.ObjectType)
	if !ok {
		t.Fatalf("healthmonitor type is %T, expected types.ObjectType", hm.GetType())
	}
	return hm, objType
}

// planHealthmonitor runs the attribute's configured plan modifiers in order and
// returns the resulting planned value.
func planHealthmonitor(ctx context.Context, hm schema.SingleNestedAttribute, plan, config, state types.Object) basetypes.ObjectValue {
	current := plan
	for _, pm := range hm.PlanModifiers {
		resp := &planmodifier.ObjectResponse{PlanValue: current}
		pm.PlanModifyObject(ctx, planmodifier.ObjectRequest{
			PlanValue:   current,
			ConfigValue: config,
			StateValue:  state,
		}, resp)
		current = resp.PlanValue
	}
	return current
}

// TestHealthmonitorOmittedBlockPlansNull guards the custom-code fix for the
// "Value Conversion Error ... Path: healthmonitor" apply crash. The model
// stores healthmonitor as a plain Go pointer (model.go), which cannot hold the
// unknown value the framework injects for a Computed+Optional attribute that is
// absent from config, so Plan.Get fails during Create. The plan modifier on the
// schema attribute must resolve that unknown to null at plan time.
//
// This asserts behaviour rather than the modifier's identity, so swapping in a
// different modifier that does not null the plan still fails. If a regeneration
// or custom-code reseal drops the wiring, this test turns that silent revert
// into a failure.
func TestHealthmonitorOmittedBlockPlansNull(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	hm, objType := healthmonitorAttr(t, ctx)

	got := planHealthmonitor(ctx, hm,
		types.ObjectUnknown(objType.AttrTypes), // framework marked the computed attribute unknown
		types.ObjectNull(objType.AttrTypes),    // block omitted from config
		types.ObjectNull(objType.AttrTypes),    // no prior state: create
	)

	if got.IsUnknown() {
		t.Fatal("healthmonitor still plans as unknown when the block is omitted; creating a pool without a healthmonitor block will fail with 'Received unknown value, however the target type cannot handle unknown values'")
	}
	if !got.IsNull() {
		t.Fatalf("expected a null planned healthmonitor on create, got %v", got)
	}
}

// TestHealthmonitorRemovalPlansNull covers removal-by-omission. Update dispatches
// to the dedicated health monitor DELETE endpoint only when the prior state has a
// monitor and the planned value is nil (resource.go). A modifier that copied the
// prior state here instead would silently make removal a no-op and leave that
// DELETE branch unreachable.
func TestHealthmonitorRemovalPlansNull(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	hm, objType := healthmonitorAttr(t, ctx)

	priorState, diags := types.ObjectValue(objType.AttrTypes, priorHealthmonitorAttrValues(objType))
	if diags.HasError() {
		t.Fatalf("failed to build prior state object: %v", diags)
	}

	got := planHealthmonitor(ctx, hm,
		types.ObjectUnknown(objType.AttrTypes), // computed attribute marked unknown
		types.ObjectNull(objType.AttrTypes),    // block removed from config
		priorState,                             // prior state has a monitor
	)

	if !got.IsNull() {
		t.Fatalf("removing the healthmonitor block must plan null so Update reaches the health monitor DELETE endpoint, got %v", got)
	}
}

// priorHealthmonitorAttrValues builds a null-valued attribute map for the
// healthmonitor object type. The values themselves are irrelevant to the
// modifier under test; only "the object is non-null" matters.
func priorHealthmonitorAttrValues(objType types.ObjectType) map[string]attr.Value {
	values := make(map[string]attr.Value, len(objType.AttrTypes))
	for name, attrType := range objType.AttrTypes {
		values[name] = attrType.ValueType(context.Background())
	}
	return values
}
