package cloud_baremetal_server

import (
	"context"
	"maps"
	"reflect"
	"runtime"
	"slices"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// fakePrivate stands in for the framework's private state, whose concrete type
// lives in an internal package and cannot be constructed from here.
type fakePrivate map[string][]byte

func (f fakePrivate) GetKey(_ context.Context, key string) ([]byte, diag.Diagnostics) {
	return f[key], nil
}

func TestImportAdoptionPending(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	pending, diags := importAdoptionPending(ctx, nil)
	if diags.HasError() {
		t.Fatalf("importAdoptionPending: %v", diags)
	}
	if pending {
		t.Error("expected no adoption pending without private state")
	}

	if pending, _ = importAdoptionPending(ctx, fakePrivate{}); pending {
		t.Error("expected no adoption pending when the key is absent")
	}

	marked := fakePrivate{importAdoptionPrivateKey: importAdoptionPrivateValue}
	if pending, _ = importAdoptionPending(ctx, marked); !pending {
		t.Error("expected adoption pending when the marker is set")
	}

	// Clearing writes an empty value rather than removing the entry, so a
	// zero-length read must count as absent or the window would never close.
	cleared := fakePrivate{importAdoptionPrivateKey: {}}
	if pending, _ = importAdoptionPending(ctx, cleared); pending {
		t.Error("expected a zero-length value to read as absent")
	}
}

// TestChangeSurvivesAdoption guards the rebuild gate. RebuildAndPoll reinstalls
// the operating system, so the only row that may suppress it is the one where
// an import genuinely explains the difference; every other row must still
// rebuild, or a real image_id or user_data change would be silently dropped.
func TestChangeSurvivesAdoption(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		changed        bool
		adopting       bool
		priorValueNull bool
		want           bool
	}{
		{
			name: "no change, nothing to do",
		},
		{
			name:    "ordinary change rebuilds",
			changed: true,
			want:    true,
		},
		{
			name:           "null prior value outside an adoption still rebuilds",
			changed:        true,
			priorValueNull: true,
			want:           true,
		},
		{
			name:     "known prior value during an adoption is a real change",
			changed:  true,
			adopting: true,
			want:     true,
		},
		{
			name:           "null prior value during an adoption is an import artefact",
			changed:        true,
			adopting:       true,
			priorValueNull: true,
			want:           false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := changeSurvivesAdoption(tt.changed, tt.adopting, tt.priorValueNull)
			if got != tt.want {
				t.Errorf("changeSurvivesAdoption(%v, %v, %v) = %v, want %v",
					tt.changed, tt.adopting, tt.priorValueNull, got, tt.want)
			}
		})
	}
}

func appConfig(values map[string]string) *map[string]jsontypes.Normalized {
	m := map[string]jsontypes.Normalized{}
	for k, v := range values {
		m[k] = jsontypes.NewNormalizedValue(v)
	}

	return &m
}

// TestAdoptedCreateOnlyAttributes pins which attributes the plan-time warning
// names. Only a transition from "absent in prior state" to "set in config"
// counts: that is the one shape where the value is recorded into state without
// ever reaching the API.
func TestAdoptedCreateOnlyAttributes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		plan  CloudBaremetalServerModel
		state CloudBaremetalServerModel
		want  []string
	}{
		{
			name: "nothing configured, nothing adopted",
		},
		{
			name:  "single attribute adopted over null state",
			plan:  CloudBaremetalServerModel{Username: types.StringValue("admin")},
			state: CloudBaremetalServerModel{Username: types.StringNull()},
			want:  []string{"username"},
		},
		{
			// image_id and user_data belong here because Update suppresses the
			// rebuild while adopting, so they too land in state unsent.
			name: "every adoptable attribute at once, reported in a stable order",
			plan: CloudBaremetalServerModel{
				ApptemplateID: types.StringValue("app-1"),
				Flavor:        types.StringValue("bm1-hf-medium"),
				ImageID:       types.StringValue("img-1"),
				NameTemplate:  types.StringValue("{ip_octets}"),
				SSHKeyName:    types.StringValue("key"),
				UserData:      types.StringValue("Zm9v"),
				Username:      types.StringValue("admin"),
				AppConfig:     appConfig(map[string]string{"k": `"v"`}),
			},
			want: []string{
				"app_config", "apptemplate_id", "flavor", "image_id",
				"name_template", "ssh_key_name", "user_data", "username",
			},
		},
		{
			name:  "ImportState back-fills flavor, so it is not adopted",
			plan:  CloudBaremetalServerModel{Flavor: types.StringValue("bm1-hf-medium")},
			state: CloudBaremetalServerModel{Flavor: types.StringValue("bm1-hf-medium")},
			want:  nil,
		},
		{
			name:  "a known prior value is a real change, not an adoption",
			plan:  CloudBaremetalServerModel{UserData: types.StringValue("bmV3")},
			state: CloudBaremetalServerModel{UserData: types.StringValue("b2xk")},
			want:  nil,
		},
		{
			name:  "removal is not an adoption",
			plan:  CloudBaremetalServerModel{Username: types.StringNull()},
			state: CloudBaremetalServerModel{Username: types.StringValue("admin")},
			want:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := adoptedCreateOnlyAttributes(tt.plan, tt.state)
			if !slices.Equal(got, tt.want) {
				t.Errorf("adoptedCreateOnlyAttributes = %v, want %v", got, tt.want)
			}
		})
	}
}

// planModifierDescriptions returns the Description of every plan modifier on a
// top-level attribute. An attribute type this does not know about is a failure,
// not a skip: a silent nil would let a regenerated set or object attribute carry
// a replace-forcing modifier past the drift test below.
func planModifierDescriptions(ctx context.Context, t *testing.T, name string, attribute schema.Attribute) []string {
	t.Helper()

	var descriptions []string

	switch typed := attribute.(type) {
	case schema.StringAttribute:
		for _, modifier := range typed.PlanModifiers {
			descriptions = append(descriptions, modifier.Description(ctx))
		}
	case schema.Int64Attribute:
		for _, modifier := range typed.PlanModifiers {
			descriptions = append(descriptions, modifier.Description(ctx))
		}
	case schema.BoolAttribute:
		for _, modifier := range typed.PlanModifiers {
			descriptions = append(descriptions, modifier.Description(ctx))
		}
	case schema.ListNestedAttribute:
		for _, modifier := range typed.PlanModifiers {
			descriptions = append(descriptions, modifier.Description(ctx))
		}
	case schema.ListAttribute:
		for _, modifier := range typed.PlanModifiers {
			descriptions = append(descriptions, modifier.Description(ctx))
		}
	case schema.MapAttribute:
		for _, modifier := range typed.PlanModifiers {
			descriptions = append(descriptions, modifier.Description(ctx))
		}
	case schema.SingleNestedAttribute:
		for _, modifier := range typed.PlanModifiers {
			descriptions = append(descriptions, modifier.Description(ctx))
		}
	default:
		t.Fatalf("attribute %q has type %T, which planModifierDescriptions cannot read plan modifiers from", name, attribute)
	}

	return descriptions
}

// rulesByModifierDescription ties each replace-forcing plan modifier to the rule
// replacementForcingAttributes has to bind for it. Descriptions are the only
// thing a modifier exposes about its behaviour, and these four are distinct.
var rulesByModifierDescription = map[string]func(plan, state, config attr.Value) bool{
	// stringplanmodifier.RequiresReplace and its typed siblings.
	"If the value of this attribute changes, Terraform will destroy and recreate the resource.": replacesOnAnyChange,
	// stringplanmodifier.RequiresReplaceIfConfigured and its typed siblings.
	"If the value of this attribute is configured and changes, Terraform will destroy and recreate the resource.": replacesWhenConfigured,
	// planmodifiers.ListRequiresReplaceIfNotNull.
	"Requires replacement when the value changes, but skips replacement when prior state is null (e.g., after terraform import with a no_refresh field).": replacesOverKnownPrior,
	// planmodifiers.*RequiresReplaceUnlessAdopting.
	"Requires replacement when the value changes, except while adopting a value into state immediately after terraform import.": replacesOverKnownPrior,
}

// TestReplacementForcingAttributesMatchSchema keeps the table that suppresses the
// adoption warning honest against the schema it mirrors. It checks both halves:
// which attributes force replacement, and which rule each one is bound to. A
// regenerated schema that adds, drops or re-modifies an attribute fails here
// rather than making the warning lie again on a plan that replaces the server.
func TestReplacementForcingAttributesMatchSchema(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	expected := map[string]func(plan, state, config attr.Value) bool{}

	for name, attribute := range ResourceSchema(ctx).Attributes {
		// password_wo is write-only: null in both plan and prior state, so its
		// RequiresReplace can never fire and it stays out of the table.
		if name == "password_wo" {
			continue
		}

		for _, description := range planModifierDescriptions(ctx, t, name, attribute) {
			rule, forcesReplacement := rulesByModifierDescription[description]
			if !forcesReplacement {
				continue
			}

			if _, already := expected[name]; already {
				t.Errorf("attribute %q carries more than one replace-forcing modifier, which this table cannot express", name)
			}

			expected[name] = rule
		}
	}

	if listed, want := slices.Sorted(maps.Keys(replacementForcingAttributes)), slices.Sorted(maps.Keys(expected)); !slices.Equal(listed, want) {
		t.Fatalf("schema forces replacement on %v, replacementForcingAttributes lists %v", want, listed)
	}

	for name, rule := range expected {
		if got, want := ruleName(replacementForcingAttributes[name]), ruleName(rule); got != want {
			t.Errorf("attribute %q is bound to %s, its schema modifier calls for %s", name, got, want)
		}
	}
}

// ruleName identifies a rule by function pointer, so the comparison above is
// about behaviour rather than about two functions that happen to look alike.
func ruleName(rule func(plan, state, config attr.Value) bool) string {
	return runtime.FuncForPC(reflect.ValueOf(rule).Pointer()).Name()
}

func TestReplacementRules(t *testing.T) {
	t.Parallel()

	null := types.StringNull()
	before := types.StringValue("a")
	after := types.StringValue("b")

	testCases := []struct {
		name                string
		rule                func(plan, state, config attr.Value) bool
		plan, state, config attr.Value
		expectedReplacement bool
	}{
		{"known prior value changed", replacesOverKnownPrior, after, before, after, true},
		{"null prior value is a creation or an adoption", replacesOverKnownPrior, after, null, after, false},
		{"unchanged value", replacesOverKnownPrior, before, before, before, false},
		{"configured change", replacesWhenConfigured, after, before, after, true},
		{"unconfigured drift", replacesWhenConfigured, after, before, null, false},
		{"any change, including over a null prior value", replacesOnAnyChange, after, null, after, true},
		{"no change", replacesOnAnyChange, before, before, before, false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			if got := testCase.rule(testCase.plan, testCase.state, testCase.config); got != testCase.expectedReplacement {
				t.Errorf("expected %v, got %v", testCase.expectedReplacement, got)
			}
		})
	}
}
