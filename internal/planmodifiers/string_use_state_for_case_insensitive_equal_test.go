package planmodifiers_test

import (
	"context"
	"testing"

	"github.com/G-Core/terraform-provider-gcore/internal/planmodifiers"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestStringUseStateForCaseInsensitiveEqual(t *testing.T) {
	t.Parallel()

	std, Std, ssd := types.StringValue("standard"), types.StringValue("Standard"), types.StringValue("ssd_hiiops")
	null, unknown := types.StringNull(), types.StringUnknown()
	live := tfsdk.State{Schema: adoptionProbeSchema, Raw: adoptionLiveRaw()}
	creating := tfsdk.State{Schema: adoptionProbeSchema, Raw: tftypes.NewValue(adoptionProbeType, nil)}

	tests := []struct {
		name          string
		state         tfsdk.State
		stateV, conf  types.String
		wantPlan      types.String
		wantReplaceIf bool // what the framework's RequiresReplaceIfConfigured then decides
	}{
		{"create leaves the config", creating, null, Std, Std, false},
		{"null config leaves the plan", live, std, null, null, false},
		{"unknown config leaves the plan (the built-in then replaces)", live, std, unknown, unknown, true},
		{"null prior leaves the config", live, null, Std, Std, true},
		{"exact equal is a no-op", live, std, std, std, false},
		// No replacement and nothing sent; the whole-resource plan is still an
		// in-place update of known-after-apply computed attributes (the
		// framework marks them unknown before this modifier runs), so only the
		// canonical spelling plans clean.
		{"case-folded equal keeps the prior value, no replacement", live, std, Std, std, false},
		{"different value replaces", live, std, ssd, ssd, true},
		// rc.1-poisoned state is not case: config equal to the recorded value is a
		// no-op; a config naming another type still replaces.
		{"recorded wrong value restated is a no-op", live, ssd, ssd, ssd, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			req := planmodifier.StringRequest{
				StateValue: tc.stateV, ConfigValue: tc.conf, PlanValue: tc.conf,
				State: tc.state, Plan: tfsdk.Plan{Schema: adoptionProbeSchema, Raw: adoptionLiveRaw()},
			}
			resp := &planmodifier.StringResponse{PlanValue: tc.conf}
			planmodifiers.StringUseStateForCaseInsensitiveEqual().PlanModifyString(context.Background(), req, resp)
			if !resp.PlanValue.Equal(tc.wantPlan) {
				t.Errorf("plan = %v, want %v", resp.PlanValue, tc.wantPlan)
			}
			// Chain the framework built-in the way the schema does.
			req.PlanValue = resp.PlanValue
			stringplanmodifier.RequiresReplaceIfConfigured().PlanModifyString(context.Background(), req, resp)
			if resp.RequiresReplace != tc.wantReplaceIf {
				t.Errorf("RequiresReplace after RequiresReplaceIfConfigured = %v, want %v", resp.RequiresReplace, tc.wantReplaceIf)
			}
		})
	}

	m := planmodifiers.StringUseStateForCaseInsensitiveEqual()
	if m.Description(context.Background()) == "" || m.MarkdownDescription(context.Background()) == "" {
		t.Error("empty description")
	}
}
