package adoptiontest_test

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/G-Core/terraform-provider-gcore/internal/importadoption/adoptiontest"
	"github.com/G-Core/terraform-provider-gcore/internal/planmodifiers"
)

// recorder captures Errorf/Fatalf so the helpers can be proven to fail closed.
type recorder struct {
	testing.TB
	errors []string
	fatal  bool
}

func (r *recorder) Helper() {}
func (r *recorder) Errorf(format string, args ...any) {
	r.errors = append(r.errors, strings.ReplaceAll(format, "%", "")+" "+strings.Join(stringify(args), " "))
}
func (r *recorder) Fatalf(format string, args ...any) {
	r.fatal = true
	r.Errorf(format, args...)
}

func stringify(args []any) []string {
	out := make([]string, 0, len(args))
	for _, a := range args {
		out = append(out, strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(sprint(a), "\n", " "), "\t", " ")))
	}
	return out
}

func sprint(a any) string {
	switch v := a.(type) {
	case string:
		return v
	case []string:
		return strings.Join(v, ",")
	default:
		return "?"
	}
}

const fixture = `package x

func (r *ThingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	resp.Diagnostics.Append(importAdoption.Advance(ctx, resp.Private)...)
	backfill(data)
}

func (r *ThingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	adopting, settle := importAdoption.BeginUpdate(ctx, req.Private, resp.Private, &resp.Diagnostics)
	defer settle()
	backfill(data)
	backfill(data)
}

func (r *ThingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(importAdoption.Mark(ctx, resp.Private)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ThingResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	importAdoptionPlan.ModifyPlan(ctx, req, resp)
}
`

// TestPinHelpersSelfCheck proves AssertWiredSource fails when a needle is
// missing and passes when every needle is present - a pin that cannot fail
// guards nothing.
func TestPinHelpersSelfCheck(t *testing.T) {
	t.Parallel()

	present := adoptiontest.StandardWiring("ThingResource", "importAdoption", "importAdoptionPlan")
	rec := &recorder{TB: t}
	adoptiontest.AssertWiredSource(rec, fixture, present)
	if len(rec.errors) != 0 {
		t.Fatalf("expected no errors for the complete fixture, got %v", rec.errors)
	}

	missing := present.With("Update", "verifyAdoption(ctx,*data,*state)")
	rec = &recorder{TB: t}
	adoptiontest.AssertWiredSource(rec, fixture, missing)
	if len(rec.errors) != 1 || !strings.Contains(rec.errors[0], "verifyAdoption") {
		t.Fatalf("expected exactly one error naming the missing needle, got %v", rec.errors)
	}

	// A needle must be found in the named method, not anywhere in the file.
	wrongMethod := adoptiontest.Wiring{Receiver: "ThingResource", Needles: map[string][]string{
		"Read": {"defersettle()"},
	}}
	rec = &recorder{TB: t}
	adoptiontest.AssertWiredSource(rec, fixture, wrongMethod)
	if len(rec.errors) != 1 {
		t.Fatalf("expected the Update-only needle to be missing from Read, got %v", rec.errors)
	}

	// An unknown receiver/method is fatal, not silently green.
	rec = &recorder{TB: t}
	adoptiontest.AssertWiredSource(rec, fixture, adoptiontest.Wiring{Receiver: "Other", Needles: map[string][]string{"Read": {"x"}}})
	if !rec.fatal {
		t.Fatal("expected Fatalf for a missing method header")
	}

	// A needle that only survives inside a comment - a call commented out
	// while hand-resolving a seal conflict - must NOT satisfy the pin.
	standard := adoptiontest.StandardWiring("ThingResource", "importAdoption", "importAdoptionPlan")
	rec = &recorder{TB: t}
	adoptiontest.AssertWiredSource(rec, commentedFixture, standard)
	if len(rec.errors) != 5 {
		t.Fatalf("expected every standard needle (5) to be reported missing from the commented-out fixture, got %d: %v", len(rec.errors), rec.errors)
	}

	// Source that does not parse is fatal.
	rec = &recorder{TB: t}
	adoptiontest.AssertWiredSource(rec, "package x\nfunc (r *ThingResource) Read(", standard)
	if !rec.fatal {
		t.Fatal("expected Fatalf for unparsable source")
	}
}

// commentedFixture carries every StandardWiring needle, but only in comments.
const commentedFixture = `package x

func (r *ThingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// resp.Diagnostics.Append(importAdoption.Advance(ctx, resp.Private)...)
	backfill(data)
}

func (r *ThingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	/* adopting, settle := importAdoption.BeginUpdate(ctx, req.Private, resp.Private, &resp.Diagnostics)
	defer settle() */
	backfill(data)
}

func (r *ThingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// resp.Diagnostics.Append(importAdoption.Mark(ctx, resp.Private)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ThingResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// importAdoptionPlan.ModifyPlan(ctx, req, resp)
}
`

func TestSchemaPinHelpersSelfCheck(t *testing.T) {
	t.Parallel()

	attrs := map[string]schema.Attribute{
		"adopting": schema.StringAttribute{Optional: true, PlanModifiers: []planmodifier.String{
			planmodifiers.StringRequiresReplaceUnlessAdopting("k"),
		}},
		"builtin": schema.StringAttribute{Computed: true, Optional: true, PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseNonNullStateForUnknown(), stringplanmodifier.RequiresReplaceIfConfigured(),
		}},
		"plain": schema.StringAttribute{Optional: true},
		"nested": schema.SingleNestedAttribute{Optional: true, Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{Optional: true},
		}, PlanModifiers: []planmodifier.Object{
			planmodifiers.ObjectRequiresReplaceUnlessAdopting("k"),
		}},
	}

	rec := &recorder{TB: t}
	adoptiontest.AssertAdoptingModifiers(rec, attrs, "adopting", "nested")
	if len(rec.errors) != 0 {
		t.Errorf("adopting: %v", rec.errors)
	}
	rec = &recorder{TB: t}
	adoptiontest.AssertAdoptingModifiers(rec, attrs, "builtin", "plain")
	if len(rec.errors) != 2 {
		t.Errorf("expected two failures for non-adopting attributes, got %v", rec.errors)
	}

	rec = &recorder{TB: t}
	adoptiontest.AssertModifierTypes(rec, attrs, "builtin",
		"stringplanmodifier.useNonNullStateForUnknown", "stringplanmodifier.requiresReplaceIfModifier")
	if len(rec.errors) != 0 {
		t.Errorf("builtin exact: %v", rec.errors)
	}
	rec = &recorder{TB: t}
	adoptiontest.AssertModifierTypes(rec, attrs, "builtin",
		"stringplanmodifier.requiresReplaceIfModifier", "stringplanmodifier.useNonNullStateForUnknown")
	if len(rec.errors) != 1 {
		t.Errorf("order must matter, got %v", rec.errors)
	}

	rec = &recorder{TB: t}
	adoptiontest.AssertComputedOptional(rec, attrs, "builtin")
	if len(rec.errors) != 0 {
		t.Errorf("builtin computed_optional: %v", rec.errors)
	}
	rec = &recorder{TB: t}
	adoptiontest.AssertComputedOptional(rec, attrs, "plain")
	if len(rec.errors) != 1 {
		t.Errorf("plain must fail AssertComputedOptional, got %v", rec.errors)
	}
}

const schemaFixture = `package x

func ResourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"source": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{planmodifiers.StringRequiresReplaceUnlessAdopting(importAdoption.Key)},
			},
			"origin": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					originHandover(),
					planmodifiers.StringRequiresReplaceUnlessAdoptingAllowingRemoval(
						importAdoption.Key,
					),
				},
			},
			"wrong_key": schema.StringAttribute{
				Optional:      true,
				PlanModifiers: []planmodifier.String{planmodifiers.StringRequiresReplaceUnlessAdopting("volume_import_adoption_pending")},
			},
			"origin_group_name": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{planmodifiers.StringUseStateUnlessAttributesChange("origin_group")},
			},
			"plain": schema.StringAttribute{Optional: true},
			"project_id": schema.Int64Attribute{
				Optional:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplaceIfConfigured()},
			},
			"resealed_id": schema.Int64Attribute{
				Optional:      true,
				PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"nested": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"source": schema.StringAttribute{
						Required:      true,
						PlanModifiers: []planmodifier.String{planmodifiers.StringRequiresReplaceUnlessAdopting("nested_key")},
					},
				},
			},
		},
	}
}
`

// TestAttributeCallArgsSelfCheck proves the constructor pin fails closed: a
// wrong key, a missing call, a sibling constructor of the same type, the
// framework's RequiresReplace() where RequiresReplaceIfConfigured() belongs, a
// same-named nested attribute and unparsable source are all reported.
func TestAttributeCallArgsSelfCheck(t *testing.T) {
	t.Parallel()

	const adopting = "planmodifiers.StringRequiresReplaceUnlessAdopting"
	const adoptingAllowingRemoval = adopting + "AllowingRemoval"

	rec := &recorder{TB: t}
	adoptiontest.AssertAttributeCallArgsSource(rec, schemaFixture, "source", adopting, "importAdoption.Key")
	if len(rec.errors) != 0 {
		t.Errorf("source: %v", rec.errors)
	}
	rec = &recorder{TB: t}
	adoptiontest.AssertAttributeCallArgsSource(rec, schemaFixture, "origin", adoptingAllowingRemoval, "importAdoption.Key")
	if len(rec.errors) != 0 {
		t.Errorf("origin: %v", rec.errors)
	}
	rec = &recorder{TB: t}
	adoptiontest.AssertAttributeCallArgsSource(rec, schemaFixture, "origin_group_name", "planmodifiers.StringUseStateUnlessAttributesChange", `"origin_group"`)
	if len(rec.errors) != 0 {
		t.Errorf("origin_group_name: %v", rec.errors)
	}
	// "" pins a call without arguments.
	rec = &recorder{TB: t}
	adoptiontest.AssertAttributeCallArgsSource(rec, schemaFixture, "project_id", "int64planmodifier.RequiresReplaceIfConfigured", "")
	if len(rec.errors) != 0 {
		t.Errorf("project_id: %v", rec.errors)
	}

	// The callee is matched whole: the longer sibling does not satisfy the
	// shorter name (both constructors return the same %T), and the shorter
	// does not satisfy the longer.
	rec = &recorder{TB: t}
	adoptiontest.AssertAttributeCallArgsSource(rec, schemaFixture, "origin", adopting, "importAdoption.Key")
	if len(rec.errors) != 1 || !strings.Contains(rec.errors[0], "no call") {
		t.Errorf("...AllowingRemoval must not satisfy the pin for the plain constructor, got %v", rec.errors)
	}
	rec = &recorder{TB: t}
	adoptiontest.AssertAttributeCallArgsSource(rec, schemaFixture, "source", adoptingAllowingRemoval, "importAdoption.Key")
	if len(rec.errors) != 1 || !strings.Contains(rec.errors[0], "no call") {
		t.Errorf("the plain constructor must not satisfy the pin for ...AllowingRemoval, got %v", rec.errors)
	}
	// A reseal to the framework's RequiresReplace() keeps the %T and drops the
	// "if configured" rule; the exact callee catches it.
	rec = &recorder{TB: t}
	adoptiontest.AssertAttributeCallArgsSource(rec, schemaFixture, "resealed_id", "int64planmodifier.RequiresReplaceIfConfigured", "")
	if len(rec.errors) != 1 || !strings.Contains(rec.errors[0], "no call") {
		t.Errorf("RequiresReplace() must not satisfy the pin for RequiresReplaceIfConfigured(), got %v", rec.errors)
	}
	// Only the top-level Attributes map counts: "nested"."source" carries the
	// call with another key and must neither satisfy nor fail the pin for the
	// top-level "source"; asked for by its own name it is absent.
	rec = &recorder{TB: t}
	adoptiontest.AssertAttributeCallArgsSource(rec, schemaFixture, "nested", adopting, "importAdoption.Key")
	if len(rec.errors) != 1 || !strings.Contains(rec.errors[0], "nested_key") {
		t.Errorf("the nested attribute's own call must be seen under its parent, got %v", rec.errors)
	}
	if got := adoptiontest.AttributeCallArgs(t, schemaFixture, "source", adopting); len(got) != 1 || got[0] != "importAdoption.Key" {
		t.Errorf("top-level source must see only its own call, got %v", got)
	}

	rec = &recorder{TB: t}
	adoptiontest.AssertAttributeCallArgsSource(rec, schemaFixture, "wrong_key", adopting, "importAdoption.Key")
	if len(rec.errors) != 1 || !strings.Contains(rec.errors[0], "volume_import_adoption_pending") {
		t.Errorf("a literal key in place of importAdoption.Key must be reported, got %v", rec.errors)
	}
	rec = &recorder{TB: t}
	adoptiontest.AssertAttributeCallArgsSource(rec, schemaFixture, "origin_group_name", "planmodifiers.StringUseStateUnlessAttributesChange", `"origin"`)
	if len(rec.errors) != 1 {
		t.Errorf("a different watched attribute must be reported, got %v", rec.errors)
	}
	rec = &recorder{TB: t}
	adoptiontest.AssertAttributeCallArgsSource(rec, schemaFixture, "plain", adopting, "importAdoption.Key")
	if len(rec.errors) != 1 {
		t.Errorf("an attribute without the call must be reported, got %v", rec.errors)
	}
	rec = &recorder{TB: t}
	adoptiontest.AssertAttributeCallArgsSource(rec, schemaFixture, "missing", adopting, "importAdoption.Key")
	if len(rec.errors) != 1 {
		t.Errorf("an attribute absent from the schema must be reported, got %v", rec.errors)
	}
	rec = &recorder{TB: t}
	adoptiontest.AssertAttributeCallArgsSource(rec, "package x\nfunc ResourceSchema(", "source", adopting, "importAdoption.Key")
	if !rec.fatal {
		t.Error("expected Fatalf for unparsable source")
	}
	rec = &recorder{TB: t}
	adoptiontest.AssertAttributeCallArgsSource(rec, "package x\nvar x = 1\n", "source", adopting, "importAdoption.Key")
	if !rec.fatal {
		t.Error("expected Fatalf for source without a schema.Schema literal")
	}
}
