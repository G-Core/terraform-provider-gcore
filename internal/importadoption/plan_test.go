package importadoption_test

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/G-Core/terraform-provider-gcore/internal/importadoption"
	"github.com/G-Core/terraform-provider-gcore/internal/importadoption/adoptiontest"
)

type model struct {
	A types.String `tfsdk:"a"`
	B types.String `tfsdk:"b"`
	C types.String `tfsdk:"c"`
}

var table = importadoption.Table[model]{
	"c": func(m model) bool { return m.C.IsNull() },
	"a": func(m model) bool { return m.A.IsNull() },
	"b": func(m model) bool { return m.B.IsNull() },
}

func TestAdoptedSorted(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		plan, state model
		want        []string
	}{
		{name: "nothing configured"},
		{
			name:  "two adopted, sorted",
			plan:  model{A: types.StringValue("x"), C: types.StringValue("y")},
			state: model{},
			want:  []string{"a", "c"},
		},
		{
			name:  "known state is a change, not an adoption",
			plan:  model{A: types.StringValue("x")},
			state: model{A: types.StringValue("old")},
			want:  nil,
		},
		{
			name:  "removal is not an adoption",
			plan:  model{},
			state: model{A: types.StringValue("old")},
			want:  nil,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := importadoption.Adopted(tc.plan, tc.state, table)
			if strings.Join(got, ",") != strings.Join(tc.want, ",") {
				t.Errorf("Adopted = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestWarningDiagnostic(t *testing.T) {
	t.Parallel()

	d := importadoption.Warning{Noun: "volume"}.Diagnostic([]string{"image_id", "source"})
	if d.Severity().String() != "Warning" {
		t.Fatalf("severity %s", d.Severity())
	}
	if d.Summary() != importadoption.Summary {
		t.Errorf("summary %q", d.Summary())
	}
	detail := d.Detail()
	for _, want := range []string{
		"If the plan above shows this volume being replaced, none of the following applies",
		`recording "image_id", "source" from configuration instead of replacing the volume`,
		"Nothing is sent to the API for them, so the configured values are trusted as-is.",
		"terraform apply -replace=",
		"one apply later when it was imported with an `import` block",
		"afterwards, setting these attributes plans a replacement",
	} {
		if !strings.Contains(detail, want) {
			t.Errorf("detail missing %q:\n%s", want, detail)
		}
	}

	custom := importadoption.Warning{Noun: "CDN resource", Transport: func(adopted []string) string {
		return "CUSTOM " + strings.Join(adopted, "+")
	}}.Diagnostic([]string{"origin"})
	if !strings.Contains(custom.Detail(), "CUSTOM origin") {
		t.Errorf("custom transport not used:\n%s", custom.Detail())
	}
	if strings.Contains(custom.Detail(), "Nothing is sent to the API") {
		t.Error("default transport leaked into a custom one")
	}
}

var (
	hookSchema = schema.Schema{Attributes: map[string]schema.Attribute{
		"a": schema.StringAttribute{Optional: true},
		"b": schema.StringAttribute{Optional: true},
		"c": schema.StringAttribute{Optional: true},
	}}
	hookType = tftypes.Object{AttributeTypes: map[string]tftypes.Type{
		"a": tftypes.String, "b": tftypes.String, "c": tftypes.String,
	}}
)

// hookRaw is a whole object with b and c null and a set to a: a string, nil
// for null, or tftypes.UnknownValue.
func hookRaw(a any) tftypes.Value {
	return tftypes.NewValue(hookType, map[string]tftypes.Value{
		"a": tftypes.NewValue(tftypes.String, a),
		"b": tftypes.NewValue(tftypes.String, nil),
		"c": tftypes.NewValue(tftypes.String, nil),
	})
}

// TestPlanHookSilentOutsideTheWindow: create and destroy plans, and any plan
// without a private-state marker, emit nothing.
func TestPlanHookSilentOutsideTheWindow(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	s := hookSchema
	set := hookRaw("x")
	unset := hookRaw(nil)
	null := tftypes.NewValue(hookType, nil)

	hook := importadoption.PlanHook[model]{Marker: marker, Table: table, Warning: importadoption.Warning{Noun: "thing"}}

	for name, req := range map[string]resource.ModifyPlanRequest{
		"create":    {Plan: tfsdk.Plan{Raw: set, Schema: s}, State: tfsdk.State{Raw: null, Schema: s}},
		"destroy":   {Plan: tfsdk.Plan{Raw: null, Schema: s}, State: tfsdk.State{Raw: set, Schema: s}},
		"no marker": {Plan: tfsdk.Plan{Raw: set, Schema: s}, State: tfsdk.State{Raw: unset, Schema: s}},
	} {
		resp := &resource.ModifyPlanResponse{Plan: req.Plan}
		hook.ModifyPlan(ctx, req, resp)
		if len(resp.Diagnostics) != 0 {
			t.Errorf("%s: unexpected diagnostics %v", name, resp.Diagnostics)
		}
	}
}

// pointerModel has a plain pointer field, like a generated model's nested
// object (the pool's healthmonitor): it decodes null, but not unknown.
type pointerModel struct {
	A *string      `tfsdk:"a"`
	B types.String `tfsdk:"b"`
	C types.String `tfsdk:"c"`
}

// TestPlanHookInsideTheWindow drives the hook with a real marker in private
// state. The warning names every adopted value, an unknown one included,
// because the plan records it. A plan that does not decode into the model
// skips the warning instead of failing: the hook is advisory.
func TestPlanHookInsideTheWindow(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	modelHook := importadoption.PlanHook[model]{Marker: marker, Table: table, Warning: importadoption.Warning{Noun: "thing"}}
	pointerHook := importadoption.PlanHook[pointerModel]{
		Marker:  marker,
		Table:   importadoption.Table[pointerModel]{"a": func(m pointerModel) bool { return m.A == nil }},
		Warning: importadoption.Warning{Noun: "thing"},
	}

	// The last row is only meaningful if its plan really does not decode.
	var decoded *pointerModel
	if d := (tfsdk.Plan{Raw: hookRaw(tftypes.UnknownValue), Schema: hookSchema}).Get(ctx, &decoded); !d.HasError() {
		t.Fatal("an unknown value must not decode into a plain pointer field")
	}

	tests := []struct {
		name string
		hook interface {
			ModifyPlan(context.Context, resource.ModifyPlanRequest, *resource.ModifyPlanResponse)
		}
		plan, state tftypes.Value
		wantWarning bool
	}{
		{name: "known adopted value warns", hook: modelHook, plan: hookRaw("x"), state: hookRaw(nil), wantWarning: true},
		{name: "unknown adopted value warns", hook: modelHook, plan: hookRaw(tftypes.UnknownValue), state: hookRaw(nil), wantWarning: true},
		{name: "nothing adopted is silent", hook: modelHook, plan: hookRaw(nil), state: hookRaw(nil)},
		{name: "a change from a known value is not an adoption", hook: modelHook, plan: hookRaw("y"), state: hookRaw("x")},
		{name: "decodable pointer model warns", hook: pointerHook, plan: hookRaw("x"), state: hookRaw(nil), wantWarning: true},
		{name: "undecodable plan skips the warning", hook: pointerHook, plan: hookRaw(tftypes.UnknownValue), state: hookRaw(nil)},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req := resource.ModifyPlanRequest{
				Plan:  tfsdk.Plan{Raw: tc.plan, Schema: hookSchema},
				State: tfsdk.State{Raw: tc.state, Schema: hookSchema},
			}
			if d := marker.Mark(ctx, adoptiontest.NewPrivate(t, &req)); d.HasError() {
				t.Fatalf("Mark: %v", d)
			}
			resp := &resource.ModifyPlanResponse{Plan: req.Plan}
			tc.hook.ModifyPlan(ctx, req, resp)

			if resp.Diagnostics.HasError() {
				t.Fatalf("the hook must never fail a plan, got %v", resp.Diagnostics)
			}
			if !tc.wantWarning {
				if len(resp.Diagnostics) != 0 {
					t.Errorf("unexpected diagnostics %v", resp.Diagnostics)
				}
				return
			}
			if len(resp.Diagnostics) != 1 || resp.Diagnostics[0].Summary() != importadoption.Summary {
				t.Fatalf("want exactly the adoption warning, got %v", resp.Diagnostics)
			}
			if detail := resp.Diagnostics[0].Detail(); !strings.Contains(detail, `recording "a" from configuration`) {
				t.Errorf("warning does not name the adopted attribute:\n%s", detail)
			}
		})
	}
}
