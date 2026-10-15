package cloud_secret

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// TestPayloadRemovedInPlace covers the full 4-bool table. Only
// (false,false,false,true) - a real resource whose prior value was set and
// whose planned value is null - counts as an in-place removal.
func TestPayloadRemovedInPlace(t *testing.T) {
	t.Parallel()

	bools := []bool{false, true}
	for _, stateRawNull := range bools {
		for _, planRawNull := range bools {
			for _, priorValueNull := range bools {
				for _, plannedValueNull := range bools {
					stateRawNull, planRawNull, priorValueNull, plannedValueNull := stateRawNull, planRawNull, priorValueNull, plannedValueNull
					name := fmt.Sprintf("stateRawNull=%v/planRawNull=%v/priorValueNull=%v/plannedValueNull=%v",
						stateRawNull, planRawNull, priorValueNull, plannedValueNull)
					t.Run(name, func(t *testing.T) {
						t.Parallel()

						want := !stateRawNull && !planRawNull && !priorValueNull && plannedValueNull
						if got := payloadRemovedInPlace(stateRawNull, planRawNull, priorValueNull, plannedValueNull); got != want {
							t.Errorf("payloadRemovedInPlace(%v,%v,%v,%v) = %v, want %v",
								stateRawNull, planRawNull, priorValueNull, plannedValueNull, got, want)
						}
					})
				}
			}
		}
	}
}

// TestWrappersDelegateToSharedModifier pins that each local wrapper's inner
// modifier is exactly the shared planmodifiers.*RequiresReplaceUnlessAdopting
// type, by %T - a hand-resolved seal that swaps in a different constructor
// (or the framework's unconditional RequiresReplace()) is invisible to a
// behavioural test alone.
func TestWrappersDelegateToSharedModifier(t *testing.T) {
	t.Parallel()

	objectInner := objectReplaceUnlessAdoptingOrRemoved("k").(objectReplaceUnlessAdoptingOrRemovedModifier).inner
	if got := fmt.Sprintf("%T", objectInner); got != "planmodifiers.objectRequiresReplaceUnlessAdoptingModifier" {
		t.Errorf("object wrapper inner = %s, want planmodifiers.objectRequiresReplaceUnlessAdoptingModifier", got)
	}

	int64Inner := int64ReplaceUnlessAdoptingOrRemoved("k").(int64ReplaceUnlessAdoptingOrRemovedModifier).inner
	if got := fmt.Sprintf("%T", int64Inner); got != "planmodifiers.int64RequiresReplaceUnlessAdoptingModifier" {
		t.Errorf("int64 wrapper inner = %s, want planmodifiers.int64RequiresReplaceUnlessAdoptingModifier", got)
	}
}

var payloadAttrTypes = map[string]attr.Type{
	"certificate_wo":       types.StringType,
	"certificate_chain_wo": types.StringType,
	"private_key_wo":       types.StringType,
}

// TestPayloadModifiersEndToEnd drives both wrappers through
// planmodifier.Object/Int64Request end to end, with Private nil (not
// adopting) - the state these modifiers are in for every plan except the one
// right after an import.
func TestPayloadModifiersEndToEnd(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	resourceSchema := ResourceSchema(ctx)
	objType := resourceSchema.Type().TerraformType(ctx).(tftypes.Object)

	allNull := map[string]tftypes.Value{}
	for name, attrType := range objType.AttributeTypes {
		allNull[name] = tftypes.NewValue(attrType, nil)
	}
	liveRaw := tftypes.NewValue(objType, allNull)
	nullRaw := tftypes.NewValue(objType, nil)

	liveState := tfsdk.State{Schema: resourceSchema, Raw: liveRaw}
	livePlan := tfsdk.Plan{Schema: resourceSchema, Raw: liveRaw}
	creatingState := tfsdk.State{Schema: resourceSchema, Raw: nullRaw}
	destroyingPlan := tfsdk.Plan{Schema: resourceSchema, Raw: nullRaw}

	setObject := types.ObjectValueMust(payloadAttrTypes, map[string]attr.Value{
		"certificate_wo":       types.StringNull(),
		"certificate_chain_wo": types.StringNull(),
		"private_key_wo":       types.StringNull(),
	})
	equalObject := types.ObjectValueMust(payloadAttrTypes, map[string]attr.Value{
		"certificate_wo":       types.StringNull(),
		"certificate_chain_wo": types.StringNull(),
		"private_key_wo":       types.StringNull(),
	})
	nullObject := types.ObjectNull(payloadAttrTypes)
	unknownObject := types.ObjectUnknown(payloadAttrTypes)

	t.Run("Object", func(t *testing.T) {
		t.Parallel()

		cases := []struct {
			name               string
			state              tfsdk.State
			plan               tfsdk.Plan
			stateValue         types.Object
			planValue          types.Object
			wantRequireReplace bool
		}{
			{"set->null (removed in place)", liveState, livePlan, setObject, nullObject, false},
			{"null->set (re-added, no marker)", liveState, livePlan, nullObject, setObject, true},
			{"set->equal (unchanged)", liveState, livePlan, setObject, equalObject, false},
			{"set->unknown (changed)", liveState, livePlan, setObject, unknownObject, true},
			{"create (State.Raw null)", creatingState, livePlan, nullObject, setObject, false},
			{"destroy (Plan.Raw null)", liveState, destroyingPlan, setObject, nullObject, false},
		}

		for _, tc := range cases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				resp := &planmodifier.ObjectResponse{PlanValue: tc.planValue}
				objectReplaceUnlessAdoptingOrRemoved("adoption_key").PlanModifyObject(ctx, planmodifier.ObjectRequest{
					State:      tc.state,
					Plan:       tc.plan,
					StateValue: tc.stateValue,
					PlanValue:  tc.planValue,
				}, resp)

				if resp.Diagnostics.HasError() {
					t.Fatalf("unexpected error diagnostics: %v", resp.Diagnostics)
				}
				if resp.RequiresReplace != tc.wantRequireReplace {
					t.Errorf("RequiresReplace = %v, want %v", resp.RequiresReplace, tc.wantRequireReplace)
				}
			})
		}
	})

	t.Run("Int64", func(t *testing.T) {
		t.Parallel()

		one := types.Int64Value(1)
		two := types.Int64Value(2)
		null := types.Int64Null()

		cases := []struct {
			name               string
			state              tfsdk.State
			plan               tfsdk.Plan
			stateValue         types.Int64
			planValue          types.Int64
			wantRequireReplace bool
		}{
			{"set->null (removed in place)", liveState, livePlan, one, null, false},
			{"null->set (re-added, no marker)", liveState, livePlan, null, one, true},
			{"set->equal (unchanged)", liveState, livePlan, one, one, false},
			{"1->2 (changed)", liveState, livePlan, one, two, true},
			{"create (State.Raw null)", creatingState, livePlan, null, one, false},
			{"destroy (Plan.Raw null)", liveState, destroyingPlan, one, null, false},
		}

		for _, tc := range cases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				resp := &planmodifier.Int64Response{PlanValue: tc.planValue}
				int64ReplaceUnlessAdoptingOrRemoved("adoption_key").PlanModifyInt64(ctx, planmodifier.Int64Request{
					State:      tc.state,
					Plan:       tc.plan,
					StateValue: tc.stateValue,
					PlanValue:  tc.planValue,
				}, resp)

				if resp.Diagnostics.HasError() {
					t.Fatalf("unexpected error diagnostics: %v", resp.Diagnostics)
				}
				if resp.RequiresReplace != tc.wantRequireReplace {
					t.Errorf("RequiresReplace = %v, want %v", resp.RequiresReplace, tc.wantRequireReplace)
				}
			})
		}
	})
}

// TestCarryComputedAttributesForward pins the fix for the review finding
// that called the seven read-only attributes going unknown on an
// Update-shaped plan "cosmetic": nothing restores them until Update runs at
// apply time, so a plan-time reader (a `count`, another resource's
// RequiresReplace attribute) sees the unknown, not the value Update ends up
// writing. This must copy every one of them back from state whenever the
// plan shows it unknown, on any plan where both State.Raw and Plan.Raw are
// non-null - not only the payload-removal one.
func TestCarryComputedAttributesForward(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	resourceSchema := ResourceSchema(ctx)
	objType := resourceSchema.Type().TerraformType(ctx).(tftypes.Object)

	known := map[string]tftypes.Value{}
	for name, attrType := range objType.AttributeTypes {
		known[name] = tftypes.NewValue(attrType, nil)
	}
	known["status"] = tftypes.NewValue(tftypes.String, "ACTIVE")
	known["algorithm"] = tftypes.NewValue(tftypes.String, "rsa")
	known["bit_length"] = tftypes.NewValue(tftypes.Number, int64(2048))
	known["mode"] = tftypes.NewValue(tftypes.String, "CBC")
	known["secret_type"] = tftypes.NewValue(tftypes.String, "certificate")
	known["created"] = tftypes.NewValue(tftypes.String, "2026-01-01T00:00:00Z")
	known["content_types"] = tftypes.NewValue(objType.AttributeTypes["content_types"], map[string]tftypes.Value{
		"default": tftypes.NewValue(tftypes.String, "application/octet-stream"),
	})
	stateRaw := tftypes.NewValue(objType, known)

	unknown := map[string]tftypes.Value{}
	for name, v := range known {
		unknown[name] = v
	}
	for _, name := range computedAttributesWithoutModifiers {
		unknown[name] = tftypes.NewValue(objType.AttributeTypes[name], tftypes.UnknownValue)
	}
	planRaw := tftypes.NewValue(objType, unknown)

	state := tfsdk.State{Schema: resourceSchema, Raw: stateRaw}
	plan := tfsdk.Plan{Schema: resourceSchema, Raw: planRaw}

	resp := &resource.ModifyPlanResponse{Plan: plan}
	carryComputedAttributesForward(ctx, resource.ModifyPlanRequest{State: state, Plan: plan}, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error diagnostics: %v", resp.Diagnostics)
	}

	for _, name := range computedAttributesWithoutModifiers {
		var got, want attr.Value
		if diags := resp.Plan.GetAttribute(ctx, path.Root(name), &got); diags.HasError() {
			t.Fatalf("GetAttribute(%q) on result: %v", name, diags)
		}
		if diags := state.GetAttribute(ctx, path.Root(name), &want); diags.HasError() {
			t.Fatalf("GetAttribute(%q) on state: %v", name, diags)
		}
		if got.IsUnknown() {
			t.Errorf("%s is still unknown after carryComputedAttributesForward", name)
			continue
		}
		if !got.Equal(want) {
			t.Errorf("%s = %v, want %v (copied from state)", name, got, want)
		}
	}
}

// TestCarryComputedAttributesForwardSkipsCreateAndDestroy pins that the
// helper never touches the plan when State.Raw or Plan.Raw is null: a
// create/destroy plan, and the null-prior-state re-plan of a replacement's
// create half (see requirePayloadOnCreate's comment), must be left alone so a
// genuine create still sees these attributes as unknown.
func TestCarryComputedAttributesForwardSkipsCreateAndDestroy(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	resourceSchema := ResourceSchema(ctx)
	objType := resourceSchema.Type().TerraformType(ctx).(tftypes.Object)
	nullRaw := tftypes.NewValue(objType, nil)

	withUnknownStatus := map[string]tftypes.Value{}
	for name, attrType := range objType.AttributeTypes {
		withUnknownStatus[name] = tftypes.NewValue(attrType, nil)
	}
	withUnknownStatus["status"] = tftypes.NewValue(tftypes.String, tftypes.UnknownValue)
	liveRaw := tftypes.NewValue(objType, withUnknownStatus)

	cases := []struct {
		name  string
		state tfsdk.State
		plan  tfsdk.Plan
	}{
		{"create (State.Raw null)", tfsdk.State{Schema: resourceSchema, Raw: nullRaw}, tfsdk.Plan{Schema: resourceSchema, Raw: liveRaw}},
		{"destroy (Plan.Raw null)", tfsdk.State{Schema: resourceSchema, Raw: liveRaw}, tfsdk.Plan{Schema: resourceSchema, Raw: nullRaw}},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			resp := &resource.ModifyPlanResponse{Plan: tc.plan}
			carryComputedAttributesForward(ctx, resource.ModifyPlanRequest{State: tc.state, Plan: tc.plan}, resp)

			if resp.Diagnostics.HasError() {
				t.Fatalf("unexpected error diagnostics: %v", resp.Diagnostics)
			}
			if !resp.Plan.Raw.Equal(tc.plan.Raw) {
				t.Errorf("plan was modified when it should have been left untouched")
			}
		})
	}
}

// TestExpirationRequiresReplace pins that an imported "...Z" expiration does
// not replace the secret when the config names the same instant with an
// offset, while a different instant, a set, a removal or an unknown still
// does.
func TestExpirationRequiresReplace(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	resourceSchema := ResourceSchema(ctx)
	objType := resourceSchema.Type().TerraformType(ctx).(tftypes.Object)
	allNull := map[string]tftypes.Value{}
	for name, attrType := range objType.AttributeTypes {
		allNull[name] = tftypes.NewValue(attrType, nil)
	}
	liveRaw := tftypes.NewValue(objType, allNull)

	str := types.StringValue
	cases := []struct {
		name        string
		state, plan types.String
		wantReplace bool
	}{
		{"same Z", str("2030-12-31T23:59:59Z"), str("2030-12-31T23:59:59Z"), false},
		{"utc offset", str("2030-12-31T23:59:59Z"), str("2030-12-31T23:59:59+00:00"), false},
		{"other offset same instant", str("2030-12-31T23:59:59Z"), str("2031-01-01T02:59:59+03:00"), false},
		{"different instant", str("2030-12-31T23:59:59Z"), str("2030-12-31T23:59:59+03:00"), true},
		{"set", types.StringNull(), str("2030-12-31T23:59:59Z"), true},
		{"removed", str("2030-12-31T23:59:59Z"), types.StringNull(), true},
		{"unknown", str("2030-12-31T23:59:59Z"), types.StringUnknown(), true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req := planmodifier.StringRequest{
				Path:        path.Root("expiration"),
				StateValue:  tc.state,
				PlanValue:   tc.plan,
				ConfigValue: tc.plan,
				State:       tfsdk.State{Schema: resourceSchema, Raw: liveRaw},
				Plan:        tfsdk.Plan{Schema: resourceSchema, Raw: liveRaw},
			}
			resp := &planmodifier.StringResponse{PlanValue: tc.plan}
			expirationRequiresReplace().PlanModifyString(ctx, req, resp)
			if resp.Diagnostics.HasError() {
				t.Fatalf("diagnostics: %v", resp.Diagnostics)
			}
			if resp.RequiresReplace != tc.wantReplace {
				t.Errorf("RequiresReplace = %v, want %v", resp.RequiresReplace, tc.wantReplace)
			}
		})
	}
}

func TestSameInstant(t *testing.T) {
	t.Parallel()

	z := timetypes.NewRFC3339ValueMust("2030-12-31T23:59:59Z")
	cases := []struct {
		name string
		b    timetypes.RFC3339
		want bool
	}{
		{"same", timetypes.NewRFC3339ValueMust("2030-12-31T23:59:59Z"), true},
		{"other offset", timetypes.NewRFC3339ValueMust("2031-01-01T02:59:59+03:00"), true},
		{"other instant", timetypes.NewRFC3339ValueMust("2030-12-31T23:59:58Z"), false},
		{"fraction dropped on the wire", timetypes.NewRFC3339ValueMust("2031-01-01T02:59:59.5+03:00"), true},
		{"null", timetypes.NewRFC3339Null(), false},
		{"unknown", timetypes.NewRFC3339Unknown(), false},
	}
	for _, tc := range cases {
		if got := sameInstant(z.StringValue, tc.b.StringValue); got != tc.want {
			t.Errorf("%s: sameInstant = %v, want %v", tc.name, got, tc.want)
		}
	}
}
