package cloud_load_balancer_pool

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/G-Core/terraform-provider-gcore/internal/apijson"
)

// poolBody renders a pool response whose healthmonitor member is the given raw
// JSON. An empty string omits the key entirely.
func poolBody(healthmonitor string) []byte {
	body := `{"id":"pool-1","name":"p","lb_algorithm":"ROUND_ROBIN","protocol":"HTTP"`
	if healthmonitor != "" {
		body += `,"healthmonitor":` + healthmonitor
	}
	return []byte(body + `}`)
}

// configuredHealthmonitor is the value Plan.Get produces for a pool whose
// configuration carries a healthmonitor block.
func configuredHealthmonitor() *CloudLoadBalancerPoolHealthmonitorModel {
	return &CloudLoadBalancerPoolHealthmonitorModel{
		Delay:      types.Int64Value(10),
		MaxRetries: types.Int64Value(3),
		Timeout:    types.Int64Value(5),
		Type:       types.StringValue("TCP"),
	}
}

// TestRemovedHealthmonitorDecodesToNull is the regression test for the apply
// failure QA reported: removing the healthmonitor block plans
// healthmonitor = null, so every response shape the API can return while the
// deletion settles must still leave state null. Without
// dropPhantomHealthmonitor the last three rows decode into an object of nulls
// and the apply dies with "Provider produced inconsistent result after apply
// (.healthmonitor: was null, but now cty.ObjectVal({...}))".
func TestRemovedHealthmonitorDecodesToNull(t *testing.T) {
	t.Parallel()

	for name, healthmonitor := range map[string]string{
		"key absent":        "",
		"explicit null":     `null`,
		"empty object":      `{}`,
		"computed only":     `{"id":"hm-1","operating_status":"ONLINE","provisioning_status":"ACTIVE"}`,
		"monitor still set": `{"delay":10,"max_retries":3,"timeout":5,"type":"TCP","max_retries_down":3,"admin_state_up":true}`,
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// The plan removed the block, so the decode target starts nil.
			data := &CloudLoadBalancerPoolModel{}
			if err := apijson.UnmarshalComputed(poolBody(healthmonitor), &data); err != nil {
				t.Fatalf("decode failed: %s", err)
			}
			resolveHealthmonitorUnknowns(data)
			dropPhantomHealthmonitor(data)

			if data.Healthmonitor != nil {
				t.Fatalf("expected a null healthmonitor, got %+v", *data.Healthmonitor)
			}
		})
	}
}

// TestConfiguredHealthmonitorSurvivesDecode guards the other direction: the fix
// must not reach into a pool that genuinely has a monitor. A configured block
// is a non-nil pointer, which UnmarshalComputed skips outright, so the planned
// values have to come through the normalization untouched.
func TestConfiguredHealthmonitorSurvivesDecode(t *testing.T) {
	t.Parallel()

	data := &CloudLoadBalancerPoolModel{Healthmonitor: configuredHealthmonitor()}
	if err := apijson.UnmarshalComputed(poolBody(`{"delay":10,"max_retries":3,"timeout":5,"type":"TCP"}`), &data); err != nil {
		t.Fatalf("decode failed: %s", err)
	}
	resolveHealthmonitorUnknowns(data)
	dropPhantomHealthmonitor(data)

	if data.Healthmonitor == nil {
		t.Fatal("a configured healthmonitor was dropped")
	}
	if got := data.Healthmonitor.Type.ValueString(); got != "TCP" {
		t.Errorf("type = %q, want TCP", got)
	}
	if got := data.Healthmonitor.Delay.ValueInt64(); got != 10 {
		t.Errorf("delay = %d, want 10", got)
	}
}

// TestRefreshKeepsLiveHealthmonitor is the safety rail on the Read path: a pool
// whose monitor is still there must survive refresh untouched. Clearing it here
// would make plan propose re-creating a monitor that already exists.
func TestRefreshKeepsLiveHealthmonitor(t *testing.T) {
	t.Parallel()

	data := &CloudLoadBalancerPoolModel{Healthmonitor: configuredHealthmonitor()}
	body := poolBody(`{"delay":10,"max_retries":3,"timeout":5,"type":"TCP","operating_status":"ONLINE"}`)
	if err := apijson.UnmarshalComputed(body, &data); err != nil {
		t.Fatalf("decode failed: %s", err)
	}

	refreshHealthmonitor(body, data)

	if data.Healthmonitor == nil {
		t.Fatal("refresh dropped a health monitor the API still reports")
	}
	if got := data.Healthmonitor.Type.ValueString(); got != "TCP" {
		t.Errorf("type = %q, want TCP", got)
	}
}

// TestReadClearsHealthmonitorDeletedOnBackend covers the half of the QA report
// that never converged: once state holds a monitor, UnmarshalComputed skips the
// pointer entirely, so an API response saying there is no monitor left it in
// place forever. Refresh reported "No changes" while plan kept proposing
// healthmonitor -> null, and each apply re-issued a DELETE that returned 400.
func TestReadClearsHealthmonitorDeletedOnBackend(t *testing.T) {
	t.Parallel()

	for name, healthmonitor := range map[string]string{
		"explicit null": `null`,
		"key absent":    "",
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Read starts from prior state, which still describes the monitor.
			data := &CloudLoadBalancerPoolModel{Healthmonitor: configuredHealthmonitor()}
			body := poolBody(healthmonitor)
			if err := apijson.UnmarshalComputed(body, &data); err != nil {
				t.Fatalf("decode failed: %s", err)
			}
			if data.Healthmonitor == nil {
				t.Fatal("the decoder cleared the healthmonitor on its own; refreshHealthmonitor may be obsolete — re-evaluate rather than deleting it")
			}

			refreshHealthmonitor(body, data)

			if data.Healthmonitor != nil {
				t.Fatalf("refresh left a deleted healthmonitor in state: %+v", *data.Healthmonitor)
			}
		})
	}
}

// TestHealthmonitorAbsent pins the predicate Read relies on. A present monitor
// must not be read as absent, or refresh would wipe a live monitor out of
// state; an unparseable body must not be either, since guessing from a broken
// response is worse than leaving state alone.
func TestHealthmonitorAbsent(t *testing.T) {
	t.Parallel()

	for name, tc := range map[string]struct {
		body []byte
		want bool
	}{
		"null":             {poolBody(`null`), true},
		"null padded":      {[]byte(`{"healthmonitor":  null }`), true},
		"absent":           {poolBody(""), true},
		"empty object":     {poolBody(`{}`), false},
		"populated":        {poolBody(`{"delay":10,"max_retries":3,"timeout":5,"type":"TCP"}`), false},
		"not json":         {[]byte(`<html>502</html>`), false},
		"empty body":       {nil, false},
		"null as string":   {poolBody(`"null"`), false},
		"nested elsewhere": {[]byte(`{"healthmonitor":{"delay":1,"max_retries":1,"timeout":1,"type":"TCP"},"other":null}`), false},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if got := healthmonitorAbsent(tc.body); got != tc.want {
				t.Errorf("healthmonitorAbsent = %v, want %v", got, tc.want)
			}
		})
	}
}
