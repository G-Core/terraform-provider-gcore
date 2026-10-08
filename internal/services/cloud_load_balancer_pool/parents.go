package cloud_load_balancer_pool

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// The API does not return listener_id / load_balancer_id. It reports a pool's
// parents through the listeners[] and loadbalancers[] relationship lists. A pool
// can list several listeners: those that use it as their default pool and
// those whose L7 policies redirect to it. The response does not say which is
// which, so any listed listener counts as a parent.

// backfillParentIDs fills a null or unknown listener_id / load_balancer_id
// with the first element of its relationship list, or null when the list is
// empty. Create and Update use it: a known planned value must reach state
// unchanged, and an unknown one must resolve to a known value.
func backfillParentIDs(ctx context.Context, data *CloudLoadBalancerPoolModel) diag.Diagnostics {
	return fillParentIDs(ctx, data, false)
}

// refreshParentIDs makes listener_id / load_balancer_id name a parent the API
// reports. A value in the relationship list is kept; any other value is
// replaced by the first element of the list, or null when the list is empty.
// Read and ImportState use it, so a configured parent the pool is not attached
// to differs from state and plans a replacement, and an omitted one follows
// the API.
func refreshParentIDs(ctx context.Context, data *CloudLoadBalancerPoolModel) diag.Diagnostics {
	return fillParentIDs(ctx, data, true)
}

func fillParentIDs(ctx context.Context, data *CloudLoadBalancerPoolModel, replaceUnlisted bool) (diags diag.Diagnostics) {
	if data == nil {
		return nil
	}
	listenerIDs, d := listenerIDs(ctx, data)
	diags.Append(d...)
	loadBalancerIDs, d := loadBalancerIDs(ctx, data)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}
	data.ListenerID = resolveParentID(data.ListenerID, listenerIDs, replaceUnlisted)
	data.LoadBalancerID = resolveParentID(data.LoadBalancerID, loadBalancerIDs, replaceUnlisted)
	return nil
}

// resolveParentID returns the parent ID to store. ids is nil when the list is
// null or unknown; then a known value is kept and an unknown one becomes null.
func resolveParentID(current types.String, ids []string, replaceUnlisted bool) types.String {
	known := !current.IsNull() && !current.IsUnknown()
	if known && (!replaceUnlisted || ids == nil || contains(ids, current.ValueString())) {
		return current
	}
	if len(ids) > 0 {
		return types.StringValue(ids[0])
	}
	return types.StringNull()
}

func contains(ids []string, id string) bool {
	for _, v := range ids {
		if v == id {
			return true
		}
	}
	return false
}

// listenerIDs returns the ids in m.listeners, or nil when the list is null or
// unknown.
func listenerIDs(ctx context.Context, m *CloudLoadBalancerPoolModel) ([]string, diag.Diagnostics) {
	if m.Listeners.IsNull() || m.Listeners.IsUnknown() {
		return nil, nil
	}
	ls, diags := m.Listeners.AsStructSliceT(ctx)
	if diags.HasError() {
		return nil, diags
	}
	out := []string{}
	for _, l := range ls {
		if !l.ID.IsNull() && !l.ID.IsUnknown() {
			out = append(out, l.ID.ValueString())
		}
	}
	return out, nil
}

// loadBalancerIDs returns the ids in m.loadbalancers, or nil when the list is
// null or unknown.
func loadBalancerIDs(ctx context.Context, m *CloudLoadBalancerPoolModel) ([]string, diag.Diagnostics) {
	if m.Loadbalancers.IsNull() || m.Loadbalancers.IsUnknown() {
		return nil, nil
	}
	lbs, diags := m.Loadbalancers.AsStructSliceT(ctx)
	if diags.HasError() {
		return nil, diags
	}
	out := []string{}
	for _, l := range lbs {
		if !l.ID.IsNull() && !l.ID.IsUnknown() {
			out = append(out, l.ID.ValueString())
		}
	}
	return out, nil
}
