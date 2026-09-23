package adoptiontest_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"

	"github.com/G-Core/terraform-provider-gcore/internal/importadoption/adoptiontest"
)

// TestNewPrivate proves the helper hands back the request's own private state,
// so a key written through it is what the code under test reads.
func TestNewPrivate(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	var plan resource.ModifyPlanRequest
	if d := adoptiontest.NewPrivate(t, &plan).SetKey(ctx, "k", []byte(`{"phase":"armed"}`)); d.HasError() {
		t.Fatalf("SetKey: %v", d)
	}
	if v, d := plan.Private.GetKey(ctx, "k"); d.HasError() || string(v) != `{"phase":"armed"}` {
		t.Errorf("ModifyPlanRequest.Private.GetKey = %s, %v", v, d)
	}

	var str planmodifier.StringRequest
	adoptiontest.NewPrivate(t, &str)
	if v, _ := str.Private.GetKey(ctx, "k"); v != nil {
		t.Errorf("a fresh private state must be empty, got %s", v)
	}

	// Anything that is not a pointer to a request with a Private field is
	// fatal, not silently a no-op that leaves the adopting path untested.
	for name, req := range map[string]any{
		"value, not pointer": plan,
		"nil pointer":        (*resource.ModifyPlanRequest)(nil),
		"no Private field":   &struct{ Other int }{},
	} {
		rec := &recorder{TB: t}
		if p := adoptiontest.NewPrivate(rec, req); p != nil || !rec.fatal {
			t.Errorf("%s: want Fatalf and nil, got fatal=%v private=%v", name, rec.fatal, p)
		}
	}
}
