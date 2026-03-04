# API Discrepancy: LoadBalancer PATCH Returns Empty vrrp_ips

## Critical Discovery

**Date:** 2025-11-17
**Issue:** GCLOUD2-20778 - vrrp_ips vanishing during LB updates

### The Problem

When updating a LoadBalancer (e.g., renaming), the `vrrp_ips` field was showing as "element has vanished" error in Terraform.

### Root Cause Identified

**The API documentation is incorrect/incomplete.**

#### API Docs Say:
From https://gcore.com/docs/api-reference/cloud/load-balancers/update-load-balancer:

```json
{
  "vrrp_ips": [
    {
      "ip_address": "127.0.0.1",
      "role": "MASTER",
      "subnet_id": "00000000-0000-4000-8000-000000000000"
    }
  ],
  ...
}
```

#### Actual API Response:
Debug logging from real API call (PATCH /loadbalancers/{id}):

```json
{
  "project_id": 6,
  "region_id": 76,
  "id": "6ef4aa51-06fe-4627-8ca1-37dc4162ade5",
  "name": "qa-lb-DEBUG-TEST",
  "vrrp_ips": [],  // ❌ EMPTY ARRAY!
  "floating_ips": [],
  "operating_status": "ONLINE",
  "provisioning_status": "ACTIVE",
  ...
}
```

**The PATCH endpoint returns `vrrp_ips: []` (empty array) even though the documentation claims it returns the full array.**

### Impact

This affects:
1. **Terraform Provider** - Cannot use Update() response directly
2. **SDK Consumers** - Cannot rely on Update() for full resource state
3. **API Documentation** - Needs correction

### Solution

**Terraform Provider Fix:**
```go
// Call Update
_, err = r.client.Cloud.LoadBalancers.Update(ctx, data.ID.ValueString(), params, ...)

// IMPORTANT: The PATCH endpoint returns vrrp_ips: [] (empty array)
// We must do an explicit GET to retrieve all computed fields
res := new(http.Response)
_, err = r.client.Cloud.LoadBalancers.Get(ctx, data.ID.ValueString(), getParams,
    option.WithResponseBodyInto(&res),
    option.WithMiddleware(logging.Middleware(ctx)),
)
bytes, _ := io.ReadAll(res.Body)
err = apijson.UnmarshalComputed(bytes, &data)
```

**Pattern:** Update + GET (matches old provider behavior)

### Why UpdateAndPoll Won't Help

Even if we add `UpdateAndPoll` to the SDK:
1. Update doesn't return tasks (synchronous operation)
2. Update response has empty `vrrp_ips` regardless
3. We still need GET to retrieve full state

### Test Results

✅ **Rename Operation:** SUCCESS
✅ **vrrp_ips Preserved:** 2 elements (MASTER + BACKUP)
✅ **No Drift:** "Your infrastructure matches the configuration"
✅ **All Computed Fields:** Correctly refreshed

### Files Changed

`/Users/user/repos/gcore-terraform/internal/services/cloud_load_balancer/resource.go`
- Lines 210-248: Update method
- Added comment explaining API discrepancy
- Implements Update + GET pattern

### Comparison with Other Resources

**LoadBalancerListener and LoadBalancerPool:**
- Their Update() endpoints return full data (not just empty fields)
- UpdateAndPoll works correctly for them
- This is a LoadBalancer-specific API issue

### Recommendations for API Team

1. **Fix PATCH /loadbalancers response:**
   - Include `vrrp_ips` array in response (not empty)
   - Include other computed fields that might be missing

2. **Update API documentation:**
   - Clarify which fields are returned by PATCH
   - Mark fields that require subsequent GET

3. **Consider consistency:**
   - Listener/Pool PATCH returns full data
   - LoadBalancer PATCH should do the same

### Alternative: SDK Convenience Method

If API won't be fixed, add to SDK:

```go
// UpdateAndRefresh updates load balancer and refreshes state via GET
// Note: Unlike other *AndPoll methods, this doesn't poll tasks because
// the Update endpoint returns synchronously but with incomplete data.
func (r *LoadBalancerService) UpdateAndRefresh(ctx context.Context,
    loadBalancerID string, params LoadBalancerUpdateParams,
    opts ...option.RequestOption) (*LoadBalancer, error) {

    // Step 1: Update
    _, err := r.Update(ctx, loadBalancerID, params, opts...)
    if err != nil {
        return nil, err
    }

    // Step 2: GET to retrieve full state including vrrp_ips
    var getParams LoadBalancerGetParams
    getParams.ProjectID = params.ProjectID
    getParams.RegionID = params.RegionID

    return r.Get(ctx, loadBalancerID, getParams, opts...)
}
```

This would:
- Centralize the workaround in SDK
- Make it reusable for all SDK consumers
- Document the API behavior clearly
- Not require API changes

## Summary

The critical bug was caused by **API returning incomplete data**, not provider/SDK bugs. The fix (Update + GET) is correct and matches old provider behavior. The API documentation needs updating to reflect actual behavior.

---

**Status:** ✅ FIXED and TESTED
**Regression Risk:** LOW (matches proven old provider pattern)
**API Issue:** HIGH (documentation vs reality mismatch)
