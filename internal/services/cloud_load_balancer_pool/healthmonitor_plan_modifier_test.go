package cloud_load_balancer_pool

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// TestHealthmonitorUnknownGuardWired guards the custom-code fix for the
// "Value Conversion Error ... Path: healthmonitor" apply crash: the model
// stores healthmonitor as a plain Go pointer, which cannot hold the unknown
// value the framework injects for Computed+Optional attributes omitted from
// config. The ObjectUseStateForUnknownWhenConfigNull plan modifier on the
// schema attribute resolves that unknown at plan time. If a regeneration or
// custom-code reseal drops the modifier, creating a pool without a
// healthmonitor block crashes again — this test turns that silent revert
// into a test failure.
func TestHealthmonitorUnknownGuardWired(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	attrRaw, ok := ResourceSchema(ctx).Attributes["healthmonitor"]
	if !ok {
		t.Fatal("healthmonitor attribute missing from resource schema")
	}
	hm, ok := attrRaw.(schema.SingleNestedAttribute)
	if !ok {
		t.Fatalf("healthmonitor is %T, expected schema.SingleNestedAttribute", attrRaw)
	}
	if !hm.Computed || !hm.Optional {
		// If the attribute is no longer Computed+Optional the unknown can no
		// longer be injected and the modifier may be obsolete — re-evaluate
		// rather than blindly re-adding it.
		t.Fatalf("healthmonitor configurability changed (Computed=%v Optional=%v); re-evaluate the unknown-guard plan modifier", hm.Computed, hm.Optional)
	}
	for _, pm := range hm.PlanModifiers {
		if strings.Contains(pm.Description(ctx), "omitted from configuration") {
			return
		}
	}
	t.Fatal("healthmonitor lost the ObjectUseStateForUnknownWhenConfigNull plan modifier; creating a pool without a healthmonitor block will fail with 'Received unknown value, however the target type cannot handle unknown values'")
}
