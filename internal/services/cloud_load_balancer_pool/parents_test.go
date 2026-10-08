package cloud_load_balancer_pool

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/G-Core/terraform-provider-gcore/internal/importadoption/adoptiontest"
)

func listeners(ctx context.Context, ids ...string) customfield.NestedObjectList[CloudLoadBalancerPoolListenersModel] {
	out := []CloudLoadBalancerPoolListenersModel{}
	for _, id := range ids {
		out = append(out, CloudLoadBalancerPoolListenersModel{ID: types.StringValue(id)})
	}
	return customfield.NewObjectListMust(ctx, out)
}

func loadbalancers(ctx context.Context, ids ...string) customfield.NestedObjectList[CloudLoadBalancerPoolLoadbalancersModel] {
	out := []CloudLoadBalancerPoolLoadbalancersModel{}
	for _, id := range ids {
		out = append(out, CloudLoadBalancerPoolLoadbalancersModel{ID: types.StringValue(id)})
	}
	return customfield.NewObjectListMust(ctx, out)
}

func TestBackfillParentIDs(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	tests := []struct {
		name         string
		in           CloudLoadBalancerPoolModel
		wantListener types.String
		wantLB       types.String
	}{
		{
			name:         "null + lists -> filled",
			in:           CloudLoadBalancerPoolModel{Listeners: listeners(ctx, "l1", "l2"), Loadbalancers: loadbalancers(ctx, "lb1")},
			wantListener: types.StringValue("l1"),
			wantLB:       types.StringValue("lb1"),
		},
		{
			name:         "unknown + lists -> filled",
			in:           CloudLoadBalancerPoolModel{ListenerID: types.StringUnknown(), LoadBalancerID: types.StringUnknown(), Listeners: listeners(ctx, "l1"), Loadbalancers: loadbalancers(ctx, "lb1")},
			wantListener: types.StringValue("l1"),
			wantLB:       types.StringValue("lb1"),
		},
		{
			name:         "unknown + empty lists -> null (state cannot hold unknown)",
			in:           CloudLoadBalancerPoolModel{ListenerID: types.StringUnknown(), LoadBalancerID: types.StringUnknown(), Listeners: listeners(ctx), Loadbalancers: loadbalancers(ctx)},
			wantListener: types.StringNull(),
			wantLB:       types.StringNull(),
		},
		{
			name:         "unknown + null lists -> null",
			in:           CloudLoadBalancerPoolModel{ListenerID: types.StringUnknown(), LoadBalancerID: types.StringUnknown(), Listeners: customfield.NullObjectList[CloudLoadBalancerPoolListenersModel](ctx), Loadbalancers: customfield.NullObjectList[CloudLoadBalancerPoolLoadbalancersModel](ctx)},
			wantListener: types.StringNull(),
			wantLB:       types.StringNull(),
		},
		{
			name:         "unknown + unknown lists -> null",
			in:           CloudLoadBalancerPoolModel{ListenerID: types.StringUnknown(), LoadBalancerID: types.StringUnknown(), Listeners: customfield.UnknownObjectList[CloudLoadBalancerPoolListenersModel](ctx), Loadbalancers: customfield.UnknownObjectList[CloudLoadBalancerPoolLoadbalancersModel](ctx)},
			wantListener: types.StringNull(),
			wantLB:       types.StringNull(),
		},
		{
			// Create and Update must store a known planned value unchanged.
			name:         "non-null keeps even when not in the list",
			in:           CloudLoadBalancerPoolModel{ListenerID: types.StringValue("stale"), LoadBalancerID: types.StringValue("stale-lb"), Listeners: listeners(ctx, "l1"), Loadbalancers: loadbalancers(ctx, "lb1")},
			wantListener: types.StringValue("stale"),
			wantLB:       types.StringValue("stale-lb"),
		},
		{
			name:         "one set, one null -> only the null is filled",
			in:           CloudLoadBalancerPoolModel{LoadBalancerID: types.StringValue("lb-cfg"), Listeners: listeners(ctx, "l1"), Loadbalancers: loadbalancers(ctx, "lb1")},
			wantListener: types.StringValue("l1"),
			wantLB:       types.StringValue("lb-cfg"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			data := tc.in
			if diags := backfillParentIDs(ctx, &data); diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}
			if !data.ListenerID.Equal(tc.wantListener) {
				t.Errorf("listener_id = %v, want %v", data.ListenerID, tc.wantListener)
			}
			if !data.LoadBalancerID.Equal(tc.wantLB) {
				t.Errorf("load_balancer_id = %v, want %v", data.LoadBalancerID, tc.wantLB)
			}
		})
	}
	if diags := backfillParentIDs(ctx, nil); diags.HasError() {
		t.Fatalf("nil model returned diagnostics: %v", diags)
	}
}

func TestBackfillParentIDsReturnsConversionDiagnostics(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	invalid := types.ListValueMust(types.StringType, []attr.Value{types.StringValue("invalid")})
	data := CloudLoadBalancerPoolModel{
		Listeners:     customfield.NestedObjectList[CloudLoadBalancerPoolListenersModel]{ListValue: invalid},
		Loadbalancers: customfield.NestedObjectList[CloudLoadBalancerPoolLoadbalancersModel]{ListValue: invalid},
	}

	diags := backfillParentIDs(ctx, &data)
	if got := len(diags.Errors()); got != 2 {
		t.Fatalf("conversion errors = %d, want 2: %v", got, diags)
	}
}

func TestRefreshParentIDs(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	tests := []struct {
		name         string
		in           CloudLoadBalancerPoolModel
		wantListener types.String
		wantLB       types.String
	}{
		{
			name:         "listed values are kept, including a non-first listener",
			in:           CloudLoadBalancerPoolModel{ListenerID: types.StringValue("l2"), LoadBalancerID: types.StringValue("lb1"), Listeners: listeners(ctx, "l1", "l2"), Loadbalancers: loadbalancers(ctx, "lb1")},
			wantListener: types.StringValue("l2"),
			wantLB:       types.StringValue("lb1"),
		},
		{
			name:         "null is filled from the lists",
			in:           CloudLoadBalancerPoolModel{Listeners: listeners(ctx, "l1"), Loadbalancers: loadbalancers(ctx, "lb1")},
			wantListener: types.StringValue("l1"),
			wantLB:       types.StringValue("lb1"),
		},
		{
			// A parent recorded without the pool being attached to it.
			name:         "unlisted values are replaced by the first listed",
			in:           CloudLoadBalancerPoolModel{ListenerID: types.StringValue("stale"), LoadBalancerID: types.StringValue("stale-lb"), Listeners: listeners(ctx, "l1"), Loadbalancers: loadbalancers(ctx, "lb1")},
			wantListener: types.StringValue("l1"),
			wantLB:       types.StringValue("lb1"),
		},
		{
			// The listener was deleted or its L7 policy removed.
			name:         "unlisted values with empty lists become null",
			in:           CloudLoadBalancerPoolModel{ListenerID: types.StringValue("gone"), LoadBalancerID: types.StringValue("lb1"), Listeners: listeners(ctx), Loadbalancers: loadbalancers(ctx, "lb1")},
			wantListener: types.StringNull(),
			wantLB:       types.StringValue("lb1"),
		},
		{
			name:         "null lists keep known values",
			in:           CloudLoadBalancerPoolModel{ListenerID: types.StringValue("l1"), LoadBalancerID: types.StringValue("lb1"), Listeners: customfield.NullObjectList[CloudLoadBalancerPoolListenersModel](ctx), Loadbalancers: customfield.NullObjectList[CloudLoadBalancerPoolLoadbalancersModel](ctx)},
			wantListener: types.StringValue("l1"),
			wantLB:       types.StringValue("lb1"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			data := tc.in
			if diags := refreshParentIDs(ctx, &data); diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}
			if !data.ListenerID.Equal(tc.wantListener) {
				t.Errorf("listener_id = %v, want %v", data.ListenerID, tc.wantListener)
			}
			if !data.LoadBalancerID.Equal(tc.wantLB) {
				t.Errorf("load_balancer_id = %v, want %v", data.LoadBalancerID, tc.wantLB)
			}
		})
	}
	if diags := refreshParentIDs(ctx, nil); diags.HasError() {
		t.Fatalf("nil model returned diagnostics: %v", diags)
	}
}

// TestParentAttributesPinned: both parent ids are Computed+Optional (removal
// from config is a no-op) with UseNonNullStateForUnknown FIRST and the
// membership-aware replace modifier - no null-adopting modifier.
func TestParentAttributesPinned(t *testing.T) {
	t.Parallel()

	attrs := ResourceSchema(context.Background()).Attributes
	adoptiontest.AssertComputedOptional(t, attrs, "listener_id", "load_balancer_id")
	for _, name := range []string{"listener_id", "load_balancer_id"} {
		adoptiontest.AssertModifierTypes(t, attrs, name,
			"stringplanmodifier.useNonNullStateForUnknown", "planmodifiers.stringRequiresReplaceUnlessMemberOfModifier")
	}
}

// TestParentChainDecisions runs the real listener_id / load_balancer_id chain:
// a configured parent that is in the relationship list is accepted in place,
// any other changed value replaces. %T cannot see
// which list the modifier reads, so this also pins the list/id names.
func TestParentChainDecisions(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	s := ResourceSchema(ctx)
	objType := s.Type().TerraformType(ctx).(tftypes.Object)
	idList := func(ids ...string) tftypes.Value {
		elemType := tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String}}
		var elems []tftypes.Value
		for _, id := range ids {
			elems = append(elems, tftypes.NewValue(elemType, map[string]tftypes.Value{"id": tftypes.NewValue(tftypes.String, id)}))
		}
		return tftypes.NewValue(tftypes.List{ElementType: elemType}, elems)
	}
	run := func(attr, list, recorded string, ids []string, config string) (types.String, bool) {
		vals := map[string]tftypes.Value{}
		for name, ty := range objType.AttributeTypes {
			vals[name] = tftypes.NewValue(ty, nil)
		}
		vals[attr] = tftypes.NewValue(tftypes.String, recorded)
		vals[list] = idList(ids...)
		raw := tftypes.NewValue(objType, vals)
		req := planmodifier.StringRequest{
			StateValue: types.StringValue(recorded), ConfigValue: types.StringValue(config), PlanValue: types.StringValue(config),
			State: tfsdk.State{Raw: raw, Schema: s}, Plan: tfsdk.Plan{Raw: raw, Schema: s},
		}
		resp := &planmodifier.StringResponse{PlanValue: req.PlanValue}
		for _, m := range s.Attributes[attr].(schema.StringAttribute).PlanModifiers {
			req.PlanValue = resp.PlanValue
			m.PlanModifyString(ctx, req, resp)
		}
		return resp.PlanValue, resp.RequiresReplace
	}

	for _, pair := range [][2]string{{"listener_id", "listeners"}, {"load_balancer_id", "loadbalancers"}} {
		attr, list := pair[0], pair[1]
		// Import back-filled the first element; the config names the second.
		if plan, replace := run(attr, list, "p1", []string{"p1", "p2"}, "p2"); replace || !plan.Equal(types.StringValue("p2")) {
			t.Errorf("%s: configured member p2 over back-filled p1: plan = %v, replace = %v; want p2 in place", attr, plan, replace)
		}
		if _, replace := run(attr, list, "p1", []string{"p1", "p2"}, "p3"); !replace {
			t.Errorf("%s: configured non-member did not plan a replacement", attr)
		}
		// State records a parent the pool is not attached to. Terraform ignores
		// a replacement for an unchanged value, so Read must not keep it; see
		// TestRefreshParentIDs.
		if plan, replace := run(attr, list, "stale", []string{"p1"}, "p1"); replace || !plan.Equal(types.StringValue("p1")) {
			t.Errorf("%s: configuring a listed parent over a wrong record: plan = %v, replace = %v; want p1 in place", attr, plan, replace)
		}
		if _, replace := run(attr, list, "stale", []string{"p1"}, "p2"); !replace {
			t.Errorf("%s: an unlisted value over a wrong record did not plan a replacement", attr)
		}
	}
}

// TestParentBackfillWiringPinned: Create and Update back-fill after every
// decode, twice in Update because it has two exits that write state; Read and
// ImportState refresh. The calls are custom hunks in generated resource.go.
func TestParentBackfillWiringPinned(t *testing.T) {
	t.Parallel()

	const needle = "resp.Diagnostics.Append(backfillParentIDs(ctx,data)...)"
	const refresh = "resp.Diagnostics.Append(refreshParentIDs(ctx,data)...)"
	adoptiontest.AssertWired(t, "resource.go", adoptiontest.Wiring{
		Receiver: "CloudLoadBalancerPoolResource",
		Needles: map[string][]string{
			"Create":      {needle},
			"Read":        {refresh},
			"ImportState": {refresh},
			"Update":      {needle},
		},
	})
	// Each state-writing exit of Update back-fills immediately before its
	// State.Set: a bare count of two is satisfied by a duplicated call on the
	// PATCH path while the {} short-circuit writes an unknown parent.
	const beforeSet = needle + "ifresp.Diagnostics.HasError(){return}resp.Diagnostics.Append(resp.State.Set("
	if n := adoptiontest.CountInMethod(t, "resource.go", "CloudLoadBalancerPoolResource", "Update", beforeSet); n != 2 {
		t.Errorf("Update back-fills immediately before State.Set %d times, want 2 (the {} short-circuit and the PATCH path)", n)
	}
}
