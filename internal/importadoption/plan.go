package importadoption

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Table maps a create-only attribute's tfsdk name (printed by the warning, so it
// must equal the model tag) to "the model has no value for it".
type Table[M any] map[string]func(M) bool

// Adopted lists the attributes this plan would record without applying them:
// unset in prior state, set in plan. Sorted so the warning text is stable.
func Adopted[M any](plan, state M, t Table[M]) []string {
	var out []string
	for name, unset := range t {
		if unset(state) && !unset(plan) {
			out = append(out, name)
		}
	}
	slices.Sort(out)
	return out
}

// Warning is the per-resource half of the plan-time warning. Transport is the
// one paragraph that differs between resources: what, if anything, reaches the
// API for the adopted values. It must be truthful for that resource - never a
// blanket "nothing is sent" where a PATCH body can carry the value.
type Warning struct {
	Noun      string
	Transport func(adopted []string) string
}

// DefaultTransport is the paragraph for a resource whose adopted values never
// reach the API at all.
func DefaultTransport([]string) string {
	return "Nothing is sent to the API for them, so the configured values are trusted as-is."
}

// Summary is the warning's fixed summary line.
const Summary = "Adopting create-only attributes into state after import"

// Diagnostic builds the plan-time warning naming the adopted attributes.
func (w Warning) Diagnostic(adopted []string) diag.Diagnostic {
	transport := w.Transport
	if transport == nil {
		transport = DefaultTransport
	}
	detail := fmt.Sprintf(
		"If the plan above shows this %[1]s being replaced, none of the following applies: "+
			"replacing it creates the %[1]s again and every configured value reaches the API.\n\n"+
			"Otherwise Terraform is recording %[2]s from configuration instead of replacing the %[1]s, "+
			"because this %[1]s was imported and the API does not return these attributes.\n\n"+
			"%[3]s\n\n"+
			"If a recorded value does not match the running %[1]s, state will be wrong and no future "+
			"refresh can correct it. Check the values before applying. To force them to actually take "+
			"effect instead, run `terraform apply -replace=` for this resource.\n\n"+
			"This adoption window closes once this apply has run and the resource has been refreshed "+
			"(one apply later when it was imported with an `import` block); afterwards, setting these "+
			"attributes plans a replacement.",
		w.Noun, strings.Join(quoteAll(adopted), ", "), transport(adopted),
	)
	return diag.NewWarningDiagnostic(Summary, detail)
}

func quoteAll(names []string) []string {
	out := make([]string, 0, len(names))
	for _, n := range names {
		out = append(out, fmt.Sprintf("%q", n))
	}
	return out
}

// PlanHook is the whole ModifyPlan body for a resource whose only plan-time
// adoption concern is the warning. It never writes private state: a plan whose
// only effect would be a marker write is a NoOp plan and Terraform drops its
// PlannedPrivate.
//
// It is advisory, so a plan or state that does not decode into M skips the
// warning instead of failing the plan. At plan time a value can still be
// unknown, and a model field that cannot hold unknown (a plain pointer to a
// nested struct, as generated models use for some objects) fails to decode it. The
// dropped diagnostics are logged at debug level. The apply is unaffected:
// Update decodes the same model from known values and reports its own errors.
type PlanHook[M any] struct {
	Marker  Marker
	Table   Table[M]
	Warning Warning
}

// ModifyPlan warns when the plan would adopt attributes from the Table.
func (h PlanHook[M]) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() { // create / destroy
		return
	}
	adopting, diags := h.Marker.Pending(ctx, req.Private)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() || !adopting {
		return
	}
	var plan, state *M
	decode := req.Plan.Get(ctx, &plan)
	decode.Append(req.State.Get(ctx, &state)...)
	if decode.HasError() {
		tflog.Debug(ctx, "Skipping the import adoption warning: plan or state does not decode", map[string]interface{}{
			"errors": errorTexts(decode),
		})
		return
	}
	if plan == nil || state == nil {
		return
	}
	if adopted := Adopted(*plan, *state, h.Table); len(adopted) > 0 {
		resp.Diagnostics.Append(h.Warning.Diagnostic(adopted))
	}
}

func errorTexts(diags diag.Diagnostics) []string {
	out := make([]string, 0, diags.ErrorsCount())
	for _, d := range diags.Errors() {
		out = append(out, d.Summary()+": "+d.Detail())
	}
	return out
}
