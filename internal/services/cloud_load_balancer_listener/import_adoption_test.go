package cloud_load_balancer_listener

import (
	"context"
	"encoding/json"
	"maps"
	"os"
	"reflect"
	"runtime"
	"slices"
	"strings"
	"testing"
	"unicode"

	"github.com/G-Core/terraform-provider-gcore/internal/customfield"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// fakePrivate stands in for the framework's private state, whose concrete type
// lives in an internal package and cannot be constructed from here.
type fakePrivate map[string][]byte

func (f fakePrivate) GetKey(_ context.Context, key string) ([]byte, diag.Diagnostics) {
	return f[key], nil
}

// SetKey mirrors the framework: an empty value deletes the key. Retirement is
// a SetKey(key, nil), so the fake must not diverge from that.
func (f fakePrivate) SetKey(_ context.Context, key string, value []byte) diag.Diagnostics {
	if len(value) == 0 {
		delete(f, key)
		return nil
	}
	f[key] = slices.Clone(value)

	return nil
}

func TestImportAdoptionPhaseOf(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		value []byte
		want  string
	}{
		"absent":           {nil, ""},
		"empty":            {[]byte{}, ""},
		"pending":          {importAdoptionValuePending, importAdoptionPhasePending},
		"armed":            {importAdoptionValueArmed, importAdoptionPhaseArmed},
		"spent":            {importAdoptionValueSpent, importAdoptionPhaseSpent},
		"malformed json":   {[]byte(`not-json`), importAdoptionPhaseArmed},
		"unknown phase":    {[]byte(`{"phase":"bogus"}`), importAdoptionPhaseArmed},
		"missing phase":    {[]byte(`{"pending":true}`), importAdoptionPhaseArmed},
		"empty phase name": {[]byte(`{"phase":""}`), importAdoptionPhaseArmed},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if got := importAdoptionPhaseOf(tc.value); got != tc.want {
				t.Errorf("importAdoptionPhaseOf(%q) = %q, want %q", tc.value, got, tc.want)
			}
		})
	}
}

// TestImportAdoptionLifecycle walks the marker through every refresh after an
// import. The window must stay open through exactly two refreshes and be gone
// after the third, so it retires even when no Update ever runs.
func TestImportAdoptionLifecycle(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	private := fakePrivate{}

	if diags := markImportAdoptionPending(ctx, private); diags.HasError() {
		t.Fatalf("mark: %v", diags)
	}

	steps := []struct {
		wantPhase   string
		wantPending bool
	}{
		{importAdoptionPhaseArmed, true},
		{importAdoptionPhaseSpent, true},
		{"", false},
		{"", false}, // advancing an absent marker is a no-op
	}

	for i, step := range steps {
		if diags := advanceImportAdoption(ctx, private); diags.HasError() {
			t.Fatalf("refresh %d: advance: %v", i+1, diags)
		}
		if got := importAdoptionPhaseOf(private[importAdoptionPrivateKey]); got != step.wantPhase {
			t.Errorf("refresh %d: phase = %q, want %q", i+1, got, step.wantPhase)
		}
		if pending, _ := importAdoptionPending(ctx, private); pending != step.wantPending {
			t.Errorf("refresh %d: pending = %v, want %v", i+1, pending, step.wantPending)
		}
	}
}

// TestImportAdoptionUpdateOutcomes pins what Update does to the window: a
// successful adopting apply closes it, a failed one restores a full apply so a
// transient error does not turn the retry into a replacement.
func TestImportAdoptionUpdateOutcomes(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	spent := fakePrivate{importAdoptionPrivateKey: importAdoptionValueSpent}
	if diags := rearmImportAdoption(ctx, spent); diags.HasError() {
		t.Fatalf("rearm: %v", diags)
	}
	if got := importAdoptionPhaseOf(spent[importAdoptionPrivateKey]); got != importAdoptionPhaseArmed {
		t.Errorf("after rearm phase = %q, want %q", got, importAdoptionPhaseArmed)
	}

	armed := fakePrivate{importAdoptionPrivateKey: importAdoptionValueArmed}
	if diags := clearImportAdoptionPending(ctx, armed); diags.HasError() {
		t.Fatalf("clear: %v", diags)
	}
	if pending, _ := importAdoptionPending(ctx, armed); pending {
		t.Error("expected the window closed after a successful adopting apply")
	}
}

func TestImportAdoptionNilPrivateState(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	if diags := markImportAdoptionPending(ctx, nil); diags.HasError() {
		t.Errorf("mark: %v", diags)
	}
	if diags := advanceImportAdoption(ctx, nil); diags.HasError() {
		t.Errorf("advance: %v", diags)
	}
	if diags := rearmImportAdoption(ctx, nil); diags.HasError() {
		t.Errorf("rearm: %v", diags)
	}
	if diags := clearImportAdoptionPending(ctx, nil); diags.HasError() {
		t.Errorf("clear: %v", diags)
	}
	if pending, diags := importAdoptionPending(ctx, nil); pending || diags.HasError() {
		t.Errorf("pending = %v, diags = %v; want false, none", pending, diags)
	}
}

func TestAdoptedCreateOnlyAttributes(t *testing.T) {
	t.Parallel()

	unset := CloudLoadBalancerListenerModel{
		DefaultPoolID:    types.StringNull(),
		InsertXForwarded: types.BoolNull(),
	}
	set := CloudLoadBalancerListenerModel{
		DefaultPoolID:    types.StringValue("pool-a"),
		InsertXForwarded: types.BoolValue(true),
	}
	onlyPool := CloudLoadBalancerListenerModel{
		DefaultPoolID:    types.StringValue("pool-a"),
		InsertXForwarded: types.BoolNull(),
	}

	tests := map[string]struct {
		plan, state CloudLoadBalancerListenerModel
		want        []string
	}{
		"import with both configured":        {set, unset, []string{"default_pool_id", "insert_x_forwarded"}},
		"import with only the pool":          {onlyPool, unset, []string{"default_pool_id"}},
		"nothing configured":                 {unset, unset, nil},
		"already in state":                   {set, set, nil},
		"configured value removed from plan": {unset, set, nil},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if got := adoptedCreateOnlyAttributes(tc.plan, tc.state); !slices.Equal(got, tc.want) {
				t.Errorf("adoptedCreateOnlyAttributes() = %v, want %v", got, tc.want)
			}
		})
	}
}

// TestPatchBaselineFor pins that an adopting apply never sends the create-only
// attributes in the PATCH body, while any other change in the same apply is
// still sent.
func TestPatchBaselineFor(t *testing.T) {
	t.Parallel()

	state := CloudLoadBalancerListenerModel{
		ID:               types.StringValue("listener-1"),
		LoadBalancerID:   types.StringValue("lb-1"),
		Protocol:         types.StringValue("HTTP"),
		ProtocolPort:     types.Int64Value(80),
		Name:             types.StringValue("old-name"),
		DefaultPoolID:    types.StringNull(),
		InsertXForwarded: types.BoolNull(),
	}
	plan := state
	plan.Name = types.StringValue("new-name")
	plan.DefaultPoolID = types.StringValue("pool-a")
	plan.InsertXForwarded = types.BoolValue(true)

	patchKeys := func(t *testing.T, adopting bool) map[string]any {
		t.Helper()
		raw, err := plan.MarshalJSONForUpdate(patchBaselineFor(state, plan, adopting))
		if err != nil {
			t.Fatalf("MarshalJSONForUpdate: %v", err)
		}
		body := map[string]any{}
		if err := json.Unmarshal(raw, &body); err != nil {
			t.Fatalf("unmarshal %s: %v", raw, err)
		}
		return body
	}

	adopting := patchKeys(t, true)
	for _, key := range []string{"default_pool_id", "insert_x_forwarded"} {
		if _, ok := adopting[key]; ok {
			t.Errorf("adopting PATCH body contains %q: %v", key, adopting)
		}
	}
	if adopting["name"] != "new-name" {
		t.Errorf("adopting PATCH body lost the real change: %v", adopting)
	}

	if _, ok := patchKeys(t, false)["default_pool_id"]; !ok {
		t.Error("outside adoption the baseline must be prior state unchanged")
	}
}

func TestImportAdoptionWarning(t *testing.T) {
	t.Parallel()

	detail := importAdoptionWarning([]string{"default_pool_id", "insert_x_forwarded"}).Detail()

	for _, want := range []string{
		`"default_pool_id", "insert_x_forwarded"`,
		// A replacement planned for another reason sends everything to the API.
		"If the plan above shows this resource being replaced",
		// The documented failed-dependency limitation and its only recovery.
		"depends on",
		"terraform state rm",
	} {
		if !strings.Contains(detail, want) {
			t.Errorf("warning is missing %q:\n%s", want, detail)
		}
	}
}

// TestImportAdoptionLifecycleWiringPinned asserts the lifecycle hooks are
// literally present in resource.go. It is a generated file and these hooks are
// sealed custom code; small hunks in generated files are exactly what can
// silently disappear when the seal is rebuilt, and every behavioural test here
// calls the helpers directly, so none of them would notice.
func TestImportAdoptionLifecycleWiringPinned(t *testing.T) {
	t.Parallel()

	raw, err := os.ReadFile("resource.go")
	if err != nil {
		t.Fatalf("reading resource.go: %s", err)
	}
	source := string(raw)

	required := map[string][]string{
		"ImportState": {"markImportAdoptionPending(ctx,resp.Private)"},
		"Read":        {"advanceImportAdoption(ctx,resp.Private)"},
		"Update": {
			"importAdoptionPending(ctx,req.Private)",
			"rearmImportAdoption(ctx,resp.Private)",
			"clearImportAdoptionPending(ctx,resp.Private)",
			"patchBaselineFor(",
		},
		"ModifyPlan": {
			"ifplanForcesReplacement(ctx,req.Plan,req.State,req.Config,adopting){return}",
			"importAdoptionWarning(",
		},
	}
	for method, calls := range required {
		body := stripSpace(methodBody(t, source, method))
		for _, call := range calls {
			if !strings.Contains(body, call) {
				t.Errorf("%s no longer contains %s", method, call)
			}
		}
	}

	modifyPlan := stripSpace(methodBody(t, source, "ModifyPlan"))

	// The replacement check has to run before the warning, or a plan that
	// replaces the listener still claims to be adopting into it.
	if check, warn := strings.Index(modifyPlan, "planForcesReplacement("), strings.Index(modifyPlan, "importAdoptionWarning("); check < 0 || check > warn {
		t.Error("ModifyPlan must call planForcesReplacement before importAdoptionWarning")
	}

	// ModifyPlan's private-state writes are discarded on a no-op plan, so the
	// marker must never be written there.
	if strings.Contains(modifyPlan, "SetKey(") {
		t.Error("ModifyPlan writes private state; plan-time writes do not persist on no-op plans")
	}
}

// TestCreateOnlyAttributesUseAdoptingModifiers pins the schema side: a regen
// that restores the generated RequiresReplace() would bring back the
// replacement after every import.
func TestCreateOnlyAttributesUseAdoptingModifiers(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	attrs := ResourceSchema(ctx).Attributes

	pool := attrs["default_pool_id"].(schema.StringAttribute).PlanModifiers
	if len(pool) != 1 || !strings.Contains(pool[0].Description(ctx), "adopting") {
		t.Errorf("default_pool_id must use the adopting modifier, got %d modifiers", len(pool))
	}

	xff := attrs["insert_x_forwarded"].(schema.BoolAttribute).PlanModifiers
	if len(xff) != 1 || !strings.Contains(xff[0].Description(ctx), "adopting") {
		t.Errorf("insert_x_forwarded must use the adopting modifier, got %d modifiers", len(xff))
	}

	for name := range createOnlyAttributes {
		if _, ok := attrs[name]; !ok {
			t.Errorf("createOnlyAttributes names %q, which is not in the schema", name)
		}
	}
}

// rulesByModifierDescription maps each replace-forcing plan modifier, by its
// Description, to the rule replacementForcingAttributes must bind it to.
var rulesByModifierDescription = map[string]replacementRule{
	// stringplanmodifier.RequiresReplace and its typed siblings.
	"If the value of this attribute changes, Terraform will destroy and recreate the resource.": replacesOnAnyChange,
	// stringplanmodifier.RequiresReplaceIfConfigured and its typed siblings.
	"If the value of this attribute is configured and changes, Terraform will destroy and recreate the resource.": replacesWhenConfigured,
	// planmodifiers.*RequiresReplaceUnlessAdopting.
	"Requires replacement when the value changes, except while adopting a value into state immediately after terraform import.": replacesUnlessAdopting,
}

// TestReplacementForcingAttributesMatchSchema keeps the table that suppresses
// the adoption warning honest against the schema: which attributes force
// replacement, and which rule each one is bound to. A regen that adds, drops or
// re-modifies one fails here instead of making the warning lie again.
func TestReplacementForcingAttributesMatchSchema(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	expected := map[string]replacementRule{}

	for name, attribute := range ResourceSchema(ctx).Attributes {
		assertNoNestedReplaceModifiers(ctx, t, name, attribute)

		for _, description := range planModifierDescriptions(ctx, t, name, attribute) {
			rule, forcesReplacement := rulesByModifierDescription[description]
			if !forcesReplacement {
				if looksReplaceForcing(description) {
					t.Errorf("attribute %q has an unrecognised replace-related modifier: %q", name, description)
				}
				continue
			}
			if _, already := expected[name]; already {
				t.Errorf("attribute %q carries more than one replace-forcing modifier", name)
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

// planModifierDescriptions returns the Description of every plan modifier on an
// attribute. It reads the PlanModifiers field by reflection so every attribute
// kind is covered, and an attribute without one is a failure, not a skip: a
// silent nil would let a regenerated attribute carry a replace-forcing modifier
// past the drift test.
func planModifierDescriptions(ctx context.Context, t *testing.T, name string, attribute any) []string {
	t.Helper()

	field := reflect.ValueOf(attribute).FieldByName("PlanModifiers")
	if !field.IsValid() || field.Kind() != reflect.Slice {
		t.Fatalf("attribute %q has type %T, which has no PlanModifiers slice to read", name, attribute)
	}

	descriptions := make([]string, 0, field.Len())
	for i := range field.Len() {
		modifier, ok := field.Index(i).Interface().(interface {
			Description(context.Context) string
		})
		if !ok {
			t.Fatalf("attribute %q has a plan modifier of type %T without a Description", name, field.Index(i).Interface())
		}
		descriptions = append(descriptions, modifier.Description(ctx))
	}

	return descriptions
}

// looksReplaceForcing catches replace-forcing modifiers whose description is not
// in rulesByModifierDescription. The framework's own RequiresReplace says
// "destroy and recreate" and never "replace", so all three words count.
func looksReplaceForcing(description string) bool {
	d := strings.ToLower(description)

	return strings.Contains(d, "replace") || strings.Contains(d, "recreate") || strings.Contains(d, "destroy")
}

// assertNoNestedReplaceModifiers fails on a replace-forcing modifier inside a
// nested attribute: the table is keyed by top-level attributes and cannot
// express one, so planForcesReplacement would miss it.
func assertNoNestedReplaceModifiers(ctx context.Context, t *testing.T, name string, attribute schema.Attribute) {
	t.Helper()

	nested, ok := attribute.(schema.NestedAttribute)
	if !ok {
		return
	}

	for childName, child := range nested.GetNestedObject().GetAttributes() {
		path := name + "." + childName
		for _, description := range planModifierDescriptions(ctx, t, path, child) {
			if _, known := rulesByModifierDescription[description]; known || looksReplaceForcing(description) {
				t.Errorf("nested attribute %q forces replacement (%q), which replacementForcingAttributes cannot express", path, description)
			}
		}
		if childSchema, isSchemaAttr := child.(schema.Attribute); isSchemaAttr {
			assertNoNestedReplaceModifiers(ctx, t, path, childSchema)
		}
	}
}

// ruleName identifies a rule by function pointer, so the comparison is about
// behaviour rather than about two functions that happen to look alike.
func ruleName(rule replacementRule) string {
	return runtime.FuncForPC(reflect.ValueOf(rule).Pointer()).Name()
}

func TestReplacementRules(t *testing.T) {
	t.Parallel()

	null := types.StringNull()
	unknown := types.StringUnknown()
	before := types.StringValue("a")
	after := types.StringValue("b")

	tests := []struct {
		name                string
		rule                replacementRule
		plan, state, config attr.Value
		adopting            bool
		want                bool
	}{
		{"configured change", replacesWhenConfigured, after, before, after, false, true},
		{"unconfigured drift", replacesWhenConfigured, after, before, null, false, false},
		{"any change, including over a null prior value", replacesOnAnyChange, after, null, after, false, true},
		{"no change", replacesOnAnyChange, before, before, before, false, false},
		{"adopting a null prior value", replacesUnlessAdopting, after, null, after, true, false},
		{"a first value outside the adoption window", replacesUnlessAdopting, after, null, after, false, true},
		{"a change from a known value while adopting", replacesUnlessAdopting, after, before, after, true, true},
		// An unknown planned value, e.g. a reference to a resource being
		// replaced in the same apply, is a change under every rule.
		{"unknown value", replacesOnAnyChange, unknown, before, unknown, false, true},
		{"unknown configured value", replacesWhenConfigured, unknown, before, unknown, false, true},
		{"unknown value over a known prior while adopting", replacesUnlessAdopting, unknown, before, unknown, true, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := tc.rule(tc.plan, tc.state, tc.config, tc.adopting); got != tc.want {
				t.Errorf("rule(...) = %v, want %v", got, tc.want)
			}
		})
	}
}

// TestPlanForcesReplacement runs the check against real plan, state and config
// objects. The QA scenario is the "protocol_port" case: an import adoption in
// the same plan as a change that forces replacement must not warn.
func TestPlanForcesReplacement(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	imported := listenerModel(ctx) // what an import leaves in state
	configured := listenerModel(ctx)
	configured.DefaultPoolID = types.StringValue("pool-a")
	configured.InsertXForwarded = types.BoolValue(true)

	portChanged := configured
	portChanged.ProtocolPort = types.Int64Value(8080)

	projectUnset := configured
	projectUnset.ProjectID = types.Int64Null()

	lbReplaced := configured
	lbReplaced.LoadBalancerID = types.StringUnknown()

	tests := []struct {
		name        string
		plan, state CloudLoadBalancerListenerModel
		adopting    bool
		want        bool
	}{
		{"adoption only", configured, imported, true, false},
		{"adoption plus protocol_port change", portChanged, imported, true, true},
		{"first set outside the adoption window", configured, imported, false, true},
		{"nothing changes", configured, configured, false, false},
		{"unconfigured project_id", projectUnset, configured, false, false},
		{"adoption while the load balancer is replaced", lbReplaced, imported, true, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			plan, state, config := planStateConfig(ctx, t, tc.plan, tc.state)
			if got := planForcesReplacement(ctx, plan, state, config, tc.adopting); got != tc.want {
				t.Errorf("planForcesReplacement() = %v, want %v", got, tc.want)
			}
		})
	}
}

// listenerModel is an imported listener: every create-only attribute null.
func listenerModel(ctx context.Context) CloudLoadBalancerListenerModel {
	return CloudLoadBalancerListenerModel{
		ID:               types.StringValue("listener-1"),
		ProjectID:        types.Int64Value(1),
		RegionID:         types.Int64Value(2),
		LoadBalancerID:   types.StringValue("lb-1"),
		Protocol:         types.StringValue("HTTP"),
		ProtocolPort:     types.Int64Value(80),
		Name:             types.StringValue("listener"),
		DefaultPoolID:    types.StringNull(),
		InsertXForwarded: types.BoolNull(),
		SniSecretID:      customfield.NullList[types.String](ctx),
	}
}

func planStateConfig(ctx context.Context, t *testing.T, planModel, stateModel CloudLoadBalancerListenerModel) (tfsdk.Plan, tfsdk.State, tfsdk.Config) {
	t.Helper()

	s := ResourceSchema(ctx)
	plan := tfsdk.Plan{Schema: s}
	state := tfsdk.State{Schema: s}

	if diags := plan.Set(ctx, &planModel); diags.HasError() {
		t.Fatalf("plan.Set: %v", diags)
	}
	if diags := state.Set(ctx, &stateModel); diags.HasError() {
		t.Fatalf("state.Set: %v", diags)
	}

	return plan, state, tfsdk.Config{Schema: s, Raw: plan.Raw}
}

func methodBody(t *testing.T, source, method string) string {
	t.Helper()

	header := "func (r *CloudLoadBalancerListenerResource) " + method + "("
	start := strings.Index(source, header)
	if start < 0 {
		t.Fatalf("could not find %q in resource.go", header)
	}
	body := source[start:]
	if end := strings.Index(body[len(header):], "\nfunc "); end >= 0 {
		body = body[:len(header)+end]
	}

	return body
}

// stripSpace removes all whitespace so the assertions survive gofmt re-wrapping.
func stripSpace(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, s)
}
