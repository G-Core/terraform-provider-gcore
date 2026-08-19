package cloud_load_balancer_pool

import (
	"encoding/json"
	"strings"
)

// dropPhantomHealthmonitor nulls a decoded healthmonitor that describes no
// health monitor at all.
//
// healthmonitor is a plain *CloudLoadBalancerPoolHealthmonitorModel tagged
// computed_optional, so apijson.UnmarshalComputed decodes it with IfUnset
// semantics (internal/apijson/decoder.go): a nil pointer is allocated for *any*
// JSON object, but every child attribute is required or optional and is
// therefore skipped, so the struct that lands in state has delay, max_retries,
// timeout and type all null — even when the API returned a fully populated
// monitor. Removing the block from configuration plans healthmonitor = null, so
// that phantom fails the apply with "Provider produced inconsistent result
// after apply (.healthmonitor: was null, but now cty.ObjectVal({...}))" and
// leaves the pool holding an object the next apply tries to delete again,
// which the API rejects with 400 "Health Monitor ... not found".
//
// All four of those attributes are required in the health monitor schema, so
// the API can never legitimately describe a monitor with every one of them
// null: an object of that shape is always a decoder artifact, never server
// state. A monitor the plan really carries is untouched — its pointer is
// non-nil, so IfUnset skips the response outright and the configured values
// survive.
func dropPhantomHealthmonitor(data *CloudLoadBalancerPoolModel) {
	if data == nil || data.Healthmonitor == nil {
		return
	}
	hm := data.Healthmonitor
	if hm.Delay.IsNull() && hm.MaxRetries.IsNull() && hm.Timeout.IsNull() && hm.Type.IsNull() {
		data.Healthmonitor = nil
	}
}

// refreshHealthmonitor reconciles a decoded healthmonitor with what the refresh
// response actually said, and is the Read-only counterpart to
// dropPhantomHealthmonitor.
//
// apijson.UnmarshalComputed skips a non-nil healthmonitor pointer wholesale, so
// a monitor that disappeared from the backend can never leave state on its own:
// refresh keeps reporting no drift while plan keeps proposing
// healthmonitor -> null, and every apply re-issues a DELETE the API answers with
// 400 "Health Monitor ... not found". Refresh has no planned value to honour, so
// it can take the response at its word. Create and Update cannot, because a
// response that lags behind the write would contradict the plan.
func refreshHealthmonitor(body []byte, data *CloudLoadBalancerPoolModel) {
	if data == nil {
		return
	}
	dropPhantomHealthmonitor(data)
	if healthmonitorAbsent(body) {
		data.Healthmonitor = nil
	}
}

// healthmonitorAbsent reports whether an API response body describes no health
// monitor — the key missing, or present and JSON null.
func healthmonitorAbsent(body []byte) bool {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(body, &fields); err != nil {
		return false
	}
	raw, ok := fields["healthmonitor"]
	if !ok {
		return true
	}
	return strings.TrimSpace(string(raw)) == "null"
}
