package cloud_baremetal_server

import (
	"context"
	"maps"
	"reflect"
	"runtime"
	"slices"
	"strings"
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
var rulesByModifierDescription = map[string]replacementRule{
	// stringplanmodifier.RequiresReplace and its typed siblings.
	"If the value of this attribute changes, Terraform will destroy and recreate the resource.": replacesOnAnyChange,
	// stringplanmodifier.RequiresReplaceIfConfigured and its typed siblings.
	"If the value of this attribute is configured and changes, Terraform will destroy and recreate the resource.": replacesWhenConfigured,
	// planmodifiers.ListRequiresReplaceIfNotNull.
	"Requires replacement when the value changes, but skips replacement when prior state is null (e.g., after terraform import with a no_refresh field).": replacesOverKnownPrior,
	// planmodifiers.*RequiresReplaceUnlessAdopting. Note this is NOT the rule
	// above: it carves out the adoption, not every null prior value.
	"Requires replacement when the value changes, except while adopting a value into state immediately after terraform import.": replacesUnlessAdopting,
}

// TestReplacementForcingAttributesMatchSchema keeps the table that suppresses the
// adoption warning honest against the schema it mirrors. It checks both halves:
// which attributes force replacement, and which rule each one is bound to. A
// regenerated schema that adds, drops or re-modifies an attribute fails here
// rather than making the warning lie again on a plan that replaces the server.
func TestReplacementForcingAttributesMatchSchema(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	expected := map[string]replacementRule{}

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
func ruleName(rule replacementRule) string {
	return runtime.FuncForPC(reflect.ValueOf(rule).Pointer()).Name()
}

func TestReplacementRules(t *testing.T) {
	t.Parallel()

	null := types.StringNull()
	before := types.StringValue("a")
	after := types.StringValue("b")

	testCases := []struct {
		name                string
		rule                replacementRule
		plan, state, config attr.Value
		adopting            bool
		expectedReplacement bool
	}{
		{name: "known prior value changed", rule: replacesOverKnownPrior, plan: after, state: before, config: after, expectedReplacement: true},
		{name: "null prior value never replaces under this rule", rule: replacesOverKnownPrior, plan: after, state: null, config: after},
		{name: "unchanged value", rule: replacesOverKnownPrior, plan: before, state: before, config: before},
		{name: "configured change", rule: replacesWhenConfigured, plan: after, state: before, config: after, expectedReplacement: true},
		{name: "unconfigured drift", rule: replacesWhenConfigured, plan: after, state: before, config: null},
		{name: "any change, including over a null prior value", rule: replacesOnAnyChange, plan: after, state: null, config: after, expectedReplacement: true},
		{name: "no change", rule: replacesOnAnyChange, plan: before, state: before, config: before},
		{name: "adopting a null prior value is not a replacement", rule: replacesUnlessAdopting, plan: after, state: null, config: after, adopting: true},
		{name: "a first value outside the adoption window does replace", rule: replacesUnlessAdopting, plan: after, state: null, config: after, expectedReplacement: true},
		{name: "a change away from a known value replaces even while adopting", rule: replacesUnlessAdopting, plan: after, state: before, config: after, adopting: true, expectedReplacement: true},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			if got := testCase.rule(testCase.plan, testCase.state, testCase.config, testCase.adopting); got != testCase.expectedReplacement {
				t.Errorf("expected %v, got %v", testCase.expectedReplacement, got)
			}
		})
	}
}

// TestRebuildsOnApply covers the decision that both warns at plan time and
// rebuilds at apply time. Update and ModifyPlan call this same function, so the
// cases below are the whole contract for both.
func TestRebuildsOnApply(t *testing.T) {
	t.Parallel()

	model := func(imageID, userData types.String) CloudBaremetalServerModel {
		return CloudBaremetalServerModel{ImageID: imageID, UserData: userData}
	}

	null := types.StringNull()
	unknown := types.StringUnknown()
	imageA, imageB := types.StringValue("image-a"), types.StringValue("image-b")
	dataA, dataB := types.StringValue("data-a"), types.StringValue("data-b")

	empty := types.StringValue("")

	testCases := []struct {
		name        string
		plan, state CloudBaremetalServerModel
		adopting    bool
		want        rebuildDecision
	}{
		{
			name:  "nothing changed",
			plan:  model(imageA, dataA),
			state: model(imageA, dataA),
			want:  rebuildNo,
		},
		{
			name:  "image changed",
			plan:  model(imageB, dataA),
			state: model(imageA, dataA),
			want:  rebuildYes,
		},
		{
			name:  "user data changed",
			plan:  model(imageA, dataB),
			state: model(imageA, dataA),
			want:  rebuildYes,
		},
		{
			name:     "adopting a null prior value writes state and touches no disk",
			plan:     model(imageA, dataA),
			state:    model(null, null),
			adopting: true,
			want:     rebuildNo,
		},
		{
			name:  "a first value outside the adoption window does rebuild",
			plan:  model(imageA, dataA),
			state: model(null, null),
			want:  rebuildYes,
		},
		{
			name:     "a change away from a known value rebuilds even while adopting",
			plan:     model(imageB, dataA),
			state:    model(imageA, dataA),
			adopting: true,
			want:     rebuildYes,
		},
		{
			name:  "removal is not acted on",
			plan:  model(null, null),
			state: model(imageA, dataA),
			want:  rebuildNo,
		},
		{
			name:  "an unknown planned value might change, so it counts",
			plan:  model(unknown, dataA),
			state: model(imageA, dataA),
			want:  rebuildUnknown,
		},
		{
			// The case the naive comparison got wrong: ValueString reports ""
			// for an unknown value, so an unknown against a null prior used to
			// compare equal and neither warn nor rebuild.
			name:  "an unknown planned value over a null prior value counts",
			plan:  model(imageA, unknown),
			state: model(imageA, null),
			want:  rebuildUnknown,
		},
		{
			// The same trap from the other side: the prior value is known and
			// is the empty string, which is exactly what an unknown stringifies
			// to.
			name:  "an unknown planned value over an empty prior value counts",
			plan:  model(imageA, unknown),
			state: model(imageA, empty),
			want:  rebuildUnknown,
		},
		{
			name:     "an unknown planned value over a null prior value is still an adoption",
			plan:     model(imageA, unknown),
			state:    model(null, null),
			adopting: true,
			want:     rebuildNo,
		},
		{
			name:     "an unknown planned value over a known prior value counts even while adopting",
			plan:     model(imageA, unknown),
			state:    model(imageA, dataA),
			adopting: true,
			want:     rebuildUnknown,
		},
		{
			// The mixed case, and the one that matters most. Update must not
			// read this as a settled rebuild: acting on it would send the
			// unresolved user_data to the API as an empty string while the
			// disks were erased for the image change beside it.
			name:  "an unresolved value outranks a definite change",
			plan:  model(imageB, unknown),
			state: model(imageA, dataA),
			want:  rebuildUnknown,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			if got := rebuildsOnApply(testCase.plan, testCase.state, testCase.adopting); got != testCase.want {
				t.Errorf("expected %v, got %v", testCase.want, got)
			}
		})
	}
}

// TestLeastResolved pins the precedence between the two attributes. An
// unresolved answer has to win, because the caller that can tell the two apart
// is Update, and what it needs to know is whether the plan can be acted on at
// all - which one unresolved attribute is enough to settle, however definite
// the other one is.
func TestLeastResolved(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		a, b, want rebuildDecision
	}{
		{a: rebuildNo, b: rebuildNo, want: rebuildNo},
		{a: rebuildNo, b: rebuildUnknown, want: rebuildUnknown},
		{a: rebuildUnknown, b: rebuildNo, want: rebuildUnknown},
		{a: rebuildUnknown, b: rebuildUnknown, want: rebuildUnknown},
		{a: rebuildNo, b: rebuildYes, want: rebuildYes},
		{a: rebuildYes, b: rebuildNo, want: rebuildYes},
		{a: rebuildYes, b: rebuildUnknown, want: rebuildUnknown},
		{a: rebuildUnknown, b: rebuildYes, want: rebuildUnknown},
		{a: rebuildYes, b: rebuildYes, want: rebuildYes},
	}

	for _, testCase := range testCases {
		t.Run(testCase.a.String()+" with "+testCase.b.String(), func(t *testing.T) {
			t.Parallel()

			if got := leastResolved(testCase.a, testCase.b); got != testCase.want {
				t.Errorf("expected %v, got %v", testCase.want, got)
			}
		})
	}
}

// TestWarningsNeverClaimAPlanAction guards the wording of both plan-time
// warnings, because with no way to detect a replacement the wording is the only
// thing keeping them correct. ModifyPlan is handed the same request for a
// `-replace` as for an ordinary update, so any sentence asserting one action or
// the other will eventually print directly above a plan showing the opposite -
// which is the bug this test exists to stop coming back.
//
// It checks three things, because checking only the first is theatre: the exact
// claims that were removed stay gone, the hedges that replaced them stay
// present, and - the part a banned-substring list cannot express - the warning
// still says the disks are destroyed. Deleting the consequence while keeping
// the hedges would otherwise pass.
//
// The banned lists name whole claims rather than words. Banning "replacement"
// outright would fail a truthful summary like "Disks are erased by an in-place
// update or a replacement", so only phrasings that assert what Terraform
// planned are listed.
func TestWarningsNeverClaimAPlanAction(t *testing.T) {
	t.Parallel()

	adoption := importAdoptionWarning([]string{"image_id"})

	for _, tc := range []struct {
		name     string
		text     string
		banned   []string
		required []string
	}{
		{
			name: "rebuild summary",
			text: rebuildWarningSummary,
			// The summary is the line the CLI prints in bold, so it is the one
			// most likely to be read on its own.
			banned:   []string{"This update", "rather than a replacement", "Terraform plans this"},
			required: []string{"erases the server's disks"},
		},
		{
			name:   "rebuild detail",
			text:   rebuildWarningDetail,
			banned: []string{"rather than a replacement", "Terraform plans this as"},
			required: []string{
				// The consequence, unconditional and before any hedge.
				"everything on the disks is lost",
				"the disks go either way",
				// Both branches, so neither action is described as the only one.
				"If the plan above shows an ordinary in-place update",
				"If the plan above shows the resource being replaced",
				// Reverting alone does not cancel a forced replacement.
				"reverting is not enough",
			},
		},
		{
			name:     "adoption summary",
			text:     adoption.Summary(),
			banned:   []string{"instead of replacing"},
			required: []string{"after import"},
		},
		{
			name:   "adoption detail",
			text:   adoption.Detail(),
			banned: []string{"instead of replacing the server", "no rebuild is performed"},
			required: []string{
				"If the plan above shows this resource being replaced, none of the following applies",
				// A concurrent rebuild does send the adopted values.
				"that rebuild sends the adopted values too",
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			for _, phrase := range tc.banned {
				if strings.Contains(tc.text, phrase) {
					t.Errorf("claims a plan action: contains %q", phrase)
				}
			}

			for _, phrase := range tc.required {
				if !strings.Contains(tc.text, phrase) {
					t.Errorf("lost a clause the warning needs to stay true: missing %q", phrase)
				}
			}
		})
	}
}
