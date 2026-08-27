package planmodifiers_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/planmodifiers"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// adoptionKey stands for the private-state key an owning resource would define.
// Every request below leaves Private nil, which reads as "no adoption pending" —
// the state these modifiers are in for every plan except the one right after an
// import. That is deliberate: the framework's private state is an internal type
// that cannot be constructed from here, so the adopting half of the matrix is
// covered by TestRequiresReplaceUnlessAdopting against the pure decision
// function, and these tests pin the typed wrappers to the non-adopting half.
const adoptionKey = "adoption_pending"

var adoptionProbeSchema = rschema.Schema{
	Attributes: map[string]rschema.Attribute{
		"a": rschema.StringAttribute{Optional: true},
	},
}

var adoptionProbeType = tftypes.Object{AttributeTypes: map[string]tftypes.Type{"a": tftypes.String}}

func adoptionLiveRaw() tftypes.Value {
	return tftypes.NewValue(adoptionProbeType, map[string]tftypes.Value{
		"a": tftypes.NewValue(tftypes.String, "a"),
	})
}

func adoptionLiveState() tfsdk.State {
	return tfsdk.State{Schema: adoptionProbeSchema, Raw: adoptionLiveRaw()}
}

func adoptionLivePlan() tfsdk.Plan {
	return tfsdk.Plan{Schema: adoptionProbeSchema, Raw: adoptionLiveRaw()}
}

// adoptionCreatingState is the whole-state-null shape terraform passes while
// planning a create; adoptionDestroyingPlan is its mirror for a destroy.
func adoptionCreatingState() tfsdk.State {
	return tfsdk.State{Schema: adoptionProbeSchema, Raw: tftypes.NewValue(adoptionProbeType, nil)}
}

func adoptionDestroyingPlan() tfsdk.Plan {
	return tfsdk.Plan{Schema: adoptionProbeSchema, Raw: tftypes.NewValue(adoptionProbeType, nil)}
}

// adoptionCase is one row of the transition matrix, expressed in terms the
// typed helpers below turn into a concrete request for their attribute type.
type adoptionCase struct {
	name string
	// stateKind and planKind select which of the type's sample values to use:
	// "null", "one" or "two". A raw-null state or plan is requested separately.
	stateKind  string
	planKind   string
	creating   bool
	destroying bool
	want       bool
}

// The rows are identical for every attribute type, which is the point: the
// typed wrappers must not drift apart from one another.
var adoptionCases = []adoptionCase{
	{name: "creation never replaces", stateKind: "null", planKind: "one", creating: true, want: false},
	{name: "destruction never replaces", stateKind: "one", planKind: "null", destroying: true, want: false},
	{name: "unchanged value does not replace", stateKind: "one", planKind: "one", want: false},
	{name: "unset on both sides does not replace", stateKind: "null", planKind: "null", want: false},
	{name: "changed known value replaces", stateKind: "one", planKind: "two", want: true},
	{name: "removal replaces", stateKind: "one", planKind: "null", want: true},
	// The case the null-prior-state heuristic gets wrong: an attribute left
	// unset at create is null in state too, so setting it later must still
	// force replacement when no import adoption is pending.
	{name: "null prior value replaces without an adoption marker", stateKind: "null", planKind: "one", want: true},
}

func adoptionRequestShape(c adoptionCase) (tfsdk.State, tfsdk.Plan) {
	state := adoptionLiveState()
	plan := adoptionLivePlan()
	if c.creating {
		state = adoptionCreatingState()
	}
	if c.destroying {
		plan = adoptionDestroyingPlan()
	}
	return state, plan
}

func TestStringRequiresReplaceUnlessAdopting(t *testing.T) {
	t.Parallel()

	pick := map[string]types.String{
		"null": types.StringNull(),
		"one":  types.StringValue("one"),
		"two":  types.StringValue("two"),
	}

	for _, tt := range adoptionCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			state, plan := adoptionRequestShape(tt)
			resp := &planmodifier.StringResponse{PlanValue: pick[tt.planKind]}
			planmodifiers.StringRequiresReplaceUnlessAdopting(adoptionKey).PlanModifyString(
				context.Background(),
				planmodifier.StringRequest{
					StateValue: pick[tt.stateKind],
					PlanValue:  pick[tt.planKind],
					State:      state,
					Plan:       plan,
				}, resp)

			if resp.Diagnostics.HasError() {
				t.Fatalf("unexpected error diagnostics: %v", resp.Diagnostics)
			}
			if resp.RequiresReplace != tt.want {
				t.Errorf("RequiresReplace = %v, want %v", resp.RequiresReplace, tt.want)
			}
		})
	}
}

func TestBoolRequiresReplaceUnlessAdopting(t *testing.T) {
	t.Parallel()

	pick := map[string]types.Bool{
		"null": types.BoolNull(),
		"one":  types.BoolValue(true),
		"two":  types.BoolValue(false),
	}

	for _, tt := range adoptionCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			state, plan := adoptionRequestShape(tt)
			resp := &planmodifier.BoolResponse{PlanValue: pick[tt.planKind]}
			planmodifiers.BoolRequiresReplaceUnlessAdopting(adoptionKey).PlanModifyBool(
				context.Background(),
				planmodifier.BoolRequest{
					StateValue: pick[tt.stateKind],
					PlanValue:  pick[tt.planKind],
					State:      state,
					Plan:       plan,
				}, resp)

			if resp.Diagnostics.HasError() {
				t.Fatalf("unexpected error diagnostics: %v", resp.Diagnostics)
			}
			if resp.RequiresReplace != tt.want {
				t.Errorf("RequiresReplace = %v, want %v", resp.RequiresReplace, tt.want)
			}
		})
	}
}

func TestInt64RequiresReplaceUnlessAdopting(t *testing.T) {
	t.Parallel()

	pick := map[string]types.Int64{
		"null": types.Int64Null(),
		"one":  types.Int64Value(1),
		"two":  types.Int64Value(2),
	}

	for _, tt := range adoptionCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			state, plan := adoptionRequestShape(tt)
			resp := &planmodifier.Int64Response{PlanValue: pick[tt.planKind]}
			planmodifiers.Int64RequiresReplaceUnlessAdopting(adoptionKey).PlanModifyInt64(
				context.Background(),
				planmodifier.Int64Request{
					StateValue: pick[tt.stateKind],
					PlanValue:  pick[tt.planKind],
					State:      state,
					Plan:       plan,
				}, resp)

			if resp.Diagnostics.HasError() {
				t.Fatalf("unexpected error diagnostics: %v", resp.Diagnostics)
			}
			if resp.RequiresReplace != tt.want {
				t.Errorf("RequiresReplace = %v, want %v", resp.RequiresReplace, tt.want)
			}
		})
	}
}

func TestMapRequiresReplaceUnlessAdopting(t *testing.T) {
	t.Parallel()

	pick := map[string]types.Map{
		"null": types.MapNull(types.StringType),
		"one":  types.MapValueMust(types.StringType, map[string]attr.Value{"k": types.StringValue("1")}),
		"two":  types.MapValueMust(types.StringType, map[string]attr.Value{"k": types.StringValue("2")}),
	}

	for _, tt := range adoptionCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			state, plan := adoptionRequestShape(tt)
			resp := &planmodifier.MapResponse{PlanValue: pick[tt.planKind]}
			planmodifiers.MapRequiresReplaceUnlessAdopting(adoptionKey).PlanModifyMap(
				context.Background(),
				planmodifier.MapRequest{
					StateValue: pick[tt.stateKind],
					PlanValue:  pick[tt.planKind],
					State:      state,
					Plan:       plan,
				}, resp)

			if resp.Diagnostics.HasError() {
				t.Fatalf("unexpected error diagnostics: %v", resp.Diagnostics)
			}
			if resp.RequiresReplace != tt.want {
				t.Errorf("RequiresReplace = %v, want %v", resp.RequiresReplace, tt.want)
			}
		})
	}
}

func TestObjectRequiresReplaceUnlessAdopting(t *testing.T) {
	t.Parallel()

	attrTypes := map[string]attr.Type{"k": types.StringType}
	pick := map[string]types.Object{
		"null": types.ObjectNull(attrTypes),
		"one":  types.ObjectValueMust(attrTypes, map[string]attr.Value{"k": types.StringValue("1")}),
		"two":  types.ObjectValueMust(attrTypes, map[string]attr.Value{"k": types.StringValue("2")}),
	}

	for _, tt := range adoptionCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			state, plan := adoptionRequestShape(tt)
			resp := &planmodifier.ObjectResponse{PlanValue: pick[tt.planKind]}
			planmodifiers.ObjectRequiresReplaceUnlessAdopting(adoptionKey).PlanModifyObject(
				context.Background(),
				planmodifier.ObjectRequest{
					StateValue: pick[tt.stateKind],
					PlanValue:  pick[tt.planKind],
					State:      state,
					Plan:       plan,
				}, resp)

			if resp.Diagnostics.HasError() {
				t.Fatalf("unexpected error diagnostics: %v", resp.Diagnostics)
			}
			if resp.RequiresReplace != tt.want {
				t.Errorf("RequiresReplace = %v, want %v", resp.RequiresReplace, tt.want)
			}
		})
	}
}
