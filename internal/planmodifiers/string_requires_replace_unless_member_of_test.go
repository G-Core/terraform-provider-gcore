package planmodifiers_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/planmodifiers"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// memberOfSchema mirrors the pool's shape: a parent id next to a nested list
// of {id} objects the API reports the parents through.
var memberOfSchema = rschema.Schema{
	Attributes: map[string]rschema.Attribute{
		"listener_id": rschema.StringAttribute{Optional: true, Computed: true},
		"listeners": rschema.ListNestedAttribute{
			Computed: true,
			NestedObject: rschema.NestedAttributeObject{
				Attributes: map[string]rschema.Attribute{"id": rschema.StringAttribute{Computed: true}},
			},
		},
	},
}

var memberOfElemType = tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.String}}
var memberOfObjType = tftypes.Object{AttributeTypes: map[string]tftypes.Type{
	"listener_id": tftypes.String,
	"listeners":   tftypes.List{ElementType: memberOfElemType},
}}

// memberOfRaw builds a whole-resource value; listeners nil means a null list,
// listenerID nil a null id.
func memberOfRaw(listenerID any, listeners []string, listUnknown bool) tftypes.Value {
	list := tftypes.NewValue(tftypes.List{ElementType: memberOfElemType}, nil)
	if listUnknown {
		list = tftypes.NewValue(tftypes.List{ElementType: memberOfElemType}, tftypes.UnknownValue)
	} else if listeners != nil {
		var elems []tftypes.Value
		for _, id := range listeners {
			elems = append(elems, tftypes.NewValue(memberOfElemType, map[string]tftypes.Value{"id": tftypes.NewValue(tftypes.String, id)}))
		}
		list = tftypes.NewValue(tftypes.List{ElementType: memberOfElemType}, elems)
	}
	return tftypes.NewValue(memberOfObjType, map[string]tftypes.Value{
		"listener_id": tftypes.NewValue(tftypes.String, listenerID),
		"listeners":   list,
	})
}

func TestStringRequiresReplaceUnlessMemberOf(t *testing.T) {
	t.Parallel()

	l1, l2, l3 := types.StringValue("l1"), types.StringValue("l2"), types.StringValue("l3")
	null, unknown := types.StringNull(), types.StringUnknown()

	tests := []struct {
		name        string
		creating    bool
		destroying  bool
		state       types.String
		listeners   []string
		listUnknown bool
		config      types.String
		plan        types.String
		wantPlan    types.String
		wantReplace bool
	}{
		{name: "create never replaces", creating: true, state: null, config: l1, plan: l1, wantPlan: l1},
		{name: "destroy never replaces", destroying: true, state: l1, listeners: []string{"l1"}, config: null, plan: null, wantPlan: null},
		{name: "unconfigured keeps whatever the plan holds", state: l1, listeners: []string{"l1"}, config: null, plan: l1, wantPlan: l1},
		// The parent listener is being replaced: its id is unknown and the pool
		// must be replaced with it (the built-in's decision).
		{name: "unknown config replaces", state: l1, listeners: []string{"l1"}, config: unknown, plan: unknown, wantPlan: unknown, wantReplace: true},
		{name: "unknown config replaces even when an earlier modifier restored state into the plan", state: l1, listeners: []string{"l1"}, config: unknown, plan: l1, wantPlan: unknown, wantReplace: true},
		{name: "unchanged does not replace", state: l1, listeners: []string{"l1"}, config: l1, plan: l1, wantPlan: l1},
		// The multi-listener import: back-fill picked l1, config names l2 -
		// a true parent, so no replacement and the plan records it.
		{name: "configured member of the list is accepted in place", state: l1, listeners: []string{"l1", "l2"}, config: l2, plan: l2, wantPlan: l2},
		{name: "configured non-member replaces", state: l1, listeners: []string{"l1", "l2"}, config: l3, plan: l3, wantPlan: l3, wantReplace: true},
		// State records a parent the list does not contain.
		{name: "stale prior + configured member corrects in place", state: types.StringValue("stale"), listeners: []string{"l1"}, config: l1, plan: l1, wantPlan: l1},
		{name: "stale prior + configured non-member replaces", state: types.StringValue("stale"), listeners: []string{"l1"}, config: l2, plan: l2, wantPlan: l2, wantReplace: true},
		{name: "null prior + configured member is accepted", state: null, listeners: []string{"l1"}, config: l1, plan: l1, wantPlan: l1},
		{name: "null prior + configured non-member replaces", state: null, listeners: []string{"l1"}, config: l2, plan: l2, wantPlan: l2, wantReplace: true},
		{name: "null list: any change replaces", state: l1, listeners: nil, config: l2, plan: l2, wantPlan: l2, wantReplace: true},
		{name: "unknown list: any change replaces", state: l1, listUnknown: true, config: l2, plan: l2, wantPlan: l2, wantReplace: true},
		{name: "empty list: any change replaces", state: l1, listeners: []string{}, config: l2, plan: l2, wantPlan: l2, wantReplace: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			stateRaw := memberOfRaw(nil, tc.listeners, tc.listUnknown)
			if !tc.state.IsNull() {
				stateRaw = memberOfRaw(tc.state.ValueString(), tc.listeners, tc.listUnknown)
			}
			if tc.creating {
				stateRaw = tftypes.NewValue(memberOfObjType, nil)
			}
			planRaw := memberOfRaw(nil, tc.listeners, tc.listUnknown)
			if tc.destroying {
				planRaw = tftypes.NewValue(memberOfObjType, nil)
			}
			req := planmodifier.StringRequest{
				StateValue: tc.state, ConfigValue: tc.config, PlanValue: tc.plan,
				State: tfsdk.State{Schema: memberOfSchema, Raw: stateRaw},
				Plan:  tfsdk.Plan{Schema: memberOfSchema, Raw: planRaw},
			}
			resp := &planmodifier.StringResponse{PlanValue: tc.plan}
			planmodifiers.StringRequiresReplaceUnlessMemberOf("listeners", "id").PlanModifyString(context.Background(), req, resp)
			if resp.Diagnostics.HasError() {
				t.Fatalf("unexpected error: %v", resp.Diagnostics)
			}
			if !resp.PlanValue.Equal(tc.wantPlan) {
				t.Errorf("plan = %v, want %v", resp.PlanValue, tc.wantPlan)
			}
			if resp.RequiresReplace != tc.wantReplace {
				t.Errorf("RequiresReplace = %v, want %v", resp.RequiresReplace, tc.wantReplace)
			}
		})
	}

	m := planmodifiers.StringRequiresReplaceUnlessMemberOf("listeners", "id")
	if m.Description(context.Background()) == "" || m.MarkdownDescription(context.Background()) == "" {
		t.Error("empty description")
	}
}

func TestStringRequiresReplaceUnlessMemberOfReturnsRelationshipDiagnostics(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	s := rschema.Schema{Attributes: map[string]rschema.Attribute{
		"listener_id": rschema.StringAttribute{Optional: true, Computed: true},
		"listeners": rschema.ListNestedAttribute{
			Computed: true,
			NestedObject: rschema.NestedAttributeObject{Attributes: map[string]rschema.Attribute{
				"id": rschema.BoolAttribute{Computed: true},
			}},
		},
	}}
	elemType := tftypes.Object{AttributeTypes: map[string]tftypes.Type{"id": tftypes.Bool}}
	objType := tftypes.Object{AttributeTypes: map[string]tftypes.Type{
		"listener_id": tftypes.String,
		"listeners":   tftypes.List{ElementType: elemType},
	}}
	raw := tftypes.NewValue(objType, map[string]tftypes.Value{
		"listener_id": tftypes.NewValue(tftypes.String, "l1"),
		"listeners": tftypes.NewValue(tftypes.List{ElementType: elemType}, []tftypes.Value{
			tftypes.NewValue(elemType, map[string]tftypes.Value{"id": tftypes.NewValue(tftypes.Bool, true)}),
		}),
	})
	req := planmodifier.StringRequest{
		StateValue:  types.StringValue("l1"),
		ConfigValue: types.StringValue("l2"),
		PlanValue:   types.StringValue("l2"),
		State:       tfsdk.State{Schema: s, Raw: raw},
		Plan:        tfsdk.Plan{Schema: s, Raw: raw},
	}
	resp := &planmodifier.StringResponse{PlanValue: req.PlanValue}

	planmodifiers.StringRequiresReplaceUnlessMemberOf("listeners", "id").PlanModifyString(ctx, req, resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected invalid relationship ID diagnostics")
	}
	if resp.RequiresReplace {
		t.Fatal("invalid relationship data must stop the replacement decision")
	}
}
