package planmodifiers_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/planmodifiers"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// inclNullCase is one row of the decision matrix shared by every typed
// UseStateForUnknownInclNull modifier. want names which of the type's sample
// values the plan must hold afterwards: "state" (copied from prior state,
// whatever it is) or "plan" (left as it came in).
type inclNullCase struct {
	name          string
	stateRawNull  bool
	stateKind     string // "null" or "one"
	planKind      string // "unknown", "null" or "one"
	configUnknown bool
	want          string
}

var inclNullCases = []inclNullCase{
	{name: "new resource leaves unknown", stateRawNull: true, stateKind: "null", planKind: "unknown", want: "plan"},
	{name: "known plan is untouched", stateKind: "one", planKind: "one", want: "plan"},
	{name: "null config + unknown plan -> prior value", stateKind: "one", planKind: "unknown", want: "state"},
	{name: "null config + unknown plan -> prior null", stateKind: "null", planKind: "unknown", want: "state"},
	// The guard the framework's UseStateForUnknown has: an unknown CONFIG value
	// (interpolation of a not-yet-known upstream) must stay unknown so a later
	// RequiresReplace* modifier sees a change instead of plan == state.
	{name: "unknown config + non-null prior stays unknown", stateKind: "one", planKind: "unknown", configUnknown: true, want: "plan"},
	{name: "unknown config + null prior stays unknown", stateKind: "null", planKind: "unknown", configUnknown: true, want: "plan"},
}

var inclNullProbeType = tftypes.Object{AttributeTypes: map[string]tftypes.Type{"a": tftypes.String}}

func inclNullState(rawNull bool) tfsdk.State {
	if rawNull {
		return tfsdk.State{Raw: tftypes.NewValue(inclNullProbeType, nil)}
	}
	return tfsdk.State{Raw: tftypes.NewValue(inclNullProbeType, map[string]tftypes.Value{"a": tftypes.NewValue(tftypes.String, "a")})}
}

func TestStringUseStateForUnknownInclNull(t *testing.T) {
	t.Parallel()

	pick := map[string]types.String{"null": types.StringNull(), "one": types.StringValue("one"), "unknown": types.StringUnknown()}
	for _, tc := range inclNullCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			config := pick["null"]
			if tc.configUnknown {
				config = types.StringUnknown()
			}
			resp := &planmodifier.StringResponse{PlanValue: pick[tc.planKind]}
			planmodifiers.StringUseStateForUnknownInclNull().PlanModifyString(context.Background(), planmodifier.StringRequest{
				StateValue: pick[tc.stateKind], PlanValue: pick[tc.planKind], ConfigValue: config, State: inclNullState(tc.stateRawNull),
			}, resp)
			want := pick[tc.planKind]
			if tc.want == "state" {
				want = pick[tc.stateKind]
			}
			if !resp.PlanValue.Equal(want) {
				t.Errorf("plan = %v, want %v", resp.PlanValue, want)
			}
		})
	}
}

func TestBoolUseStateForUnknownInclNull(t *testing.T) {
	t.Parallel()

	pick := map[string]types.Bool{"null": types.BoolNull(), "one": types.BoolValue(true), "unknown": types.BoolUnknown()}
	for _, tc := range inclNullCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			config := pick["null"]
			if tc.configUnknown {
				config = types.BoolUnknown()
			}
			resp := &planmodifier.BoolResponse{PlanValue: pick[tc.planKind]}
			planmodifiers.BoolUseStateForUnknownInclNull().PlanModifyBool(context.Background(), planmodifier.BoolRequest{
				StateValue: pick[tc.stateKind], PlanValue: pick[tc.planKind], ConfigValue: config, State: inclNullState(tc.stateRawNull),
			}, resp)
			want := pick[tc.planKind]
			if tc.want == "state" {
				want = pick[tc.stateKind]
			}
			if !resp.PlanValue.Equal(want) {
				t.Errorf("plan = %v, want %v", resp.PlanValue, want)
			}
		})
	}
}

func TestInt64UseStateForUnknownInclNull(t *testing.T) {
	t.Parallel()

	pick := map[string]types.Int64{"null": types.Int64Null(), "one": types.Int64Value(1), "unknown": types.Int64Unknown()}
	for _, tc := range inclNullCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			config := pick["null"]
			if tc.configUnknown {
				config = types.Int64Unknown()
			}
			resp := &planmodifier.Int64Response{PlanValue: pick[tc.planKind]}
			planmodifiers.Int64UseStateForUnknownInclNull().PlanModifyInt64(context.Background(), planmodifier.Int64Request{
				StateValue: pick[tc.stateKind], PlanValue: pick[tc.planKind], ConfigValue: config, State: inclNullState(tc.stateRawNull),
			}, resp)
			want := pick[tc.planKind]
			if tc.want == "state" {
				want = pick[tc.stateKind]
			}
			if !resp.PlanValue.Equal(want) {
				t.Errorf("plan = %v, want %v", resp.PlanValue, want)
			}
		})
	}
}

func TestSetUseStateForUnknownInclNull(t *testing.T) {
	t.Parallel()

	pick := map[string]types.Set{
		"null":    types.SetNull(types.StringType),
		"one":     types.SetValueMust(types.StringType, []attr.Value{types.StringValue("one")}),
		"unknown": types.SetUnknown(types.StringType),
	}
	for _, tc := range inclNullCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			config := pick["null"]
			if tc.configUnknown {
				config = types.SetUnknown(types.StringType)
			}
			resp := &planmodifier.SetResponse{PlanValue: pick[tc.planKind]}
			planmodifiers.SetUseStateForUnknownInclNull().PlanModifySet(context.Background(), planmodifier.SetRequest{
				StateValue: pick[tc.stateKind], PlanValue: pick[tc.planKind], ConfigValue: config, State: inclNullState(tc.stateRawNull),
			}, resp)
			want := pick[tc.planKind]
			if tc.want == "state" {
				want = pick[tc.stateKind]
			}
			if !resp.PlanValue.Equal(want) {
				t.Errorf("plan = %v, want %v", resp.PlanValue, want)
			}
		})
	}
}

func TestObjectUseStateForUnknownInclNull(t *testing.T) {
	t.Parallel()

	attrTypes := map[string]attr.Type{"k": types.StringType}
	pick := map[string]types.Object{
		"null":    types.ObjectNull(attrTypes),
		"one":     types.ObjectValueMust(attrTypes, map[string]attr.Value{"k": types.StringValue("1")}),
		"unknown": types.ObjectUnknown(attrTypes),
	}
	for _, tc := range inclNullCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			config := pick["null"]
			if tc.configUnknown {
				config = types.ObjectUnknown(attrTypes)
			}
			resp := &planmodifier.ObjectResponse{PlanValue: pick[tc.planKind]}
			planmodifiers.ObjectUseStateForUnknownInclNull().PlanModifyObject(context.Background(), planmodifier.ObjectRequest{
				StateValue: pick[tc.stateKind], PlanValue: pick[tc.planKind], ConfigValue: config, State: inclNullState(tc.stateRawNull),
			}, resp)
			want := pick[tc.planKind]
			if tc.want == "state" {
				want = pick[tc.stateKind]
			}
			if !resp.PlanValue.Equal(want) {
				t.Errorf("plan = %v, want %v", resp.PlanValue, want)
			}
		})
	}
}

// TestInclNullThenRequiresReplaceUnlessAdopting_UnknownConfig chains the two
// modifiers the way a create-only Computed+Optional attribute wires them
// (InclNull first): an unknown configured value must reach the replace
// modifier still unknown, so it plans the replacement in the same run instead
// of a no-op now and a replacement on the next plan.
func TestInclNullThenRequiresReplaceUnlessAdopting_UnknownConfig(t *testing.T) {
	t.Parallel()

	for _, prior := range []types.String{types.StringValue("old.example.com"), types.StringNull()} {
		state := tfsdk.State{Schema: adoptionProbeSchema, Raw: adoptionLiveRaw()}
		plan := tfsdk.Plan{Schema: adoptionProbeSchema, Raw: adoptionLiveRaw()}
		req := planmodifier.StringRequest{
			StateValue: prior, PlanValue: types.StringUnknown(), ConfigValue: types.StringUnknown(), State: state, Plan: plan,
		}
		first := &planmodifier.StringResponse{PlanValue: req.PlanValue}
		planmodifiers.StringUseStateForUnknownInclNull().PlanModifyString(context.Background(), req, first)
		if !first.PlanValue.IsUnknown() {
			t.Fatalf("prior %v: InclNull rewrote an unknown configured value to %v", prior, first.PlanValue)
		}
		req.PlanValue = first.PlanValue
		second := &planmodifier.StringResponse{PlanValue: req.PlanValue}
		planmodifiers.StringRequiresReplaceUnlessAdopting(adoptionKey).PlanModifyString(context.Background(), req, second)
		if !second.RequiresReplace {
			t.Errorf("prior %v: unknown configured value did not plan a replacement", prior)
		}
	}
}
