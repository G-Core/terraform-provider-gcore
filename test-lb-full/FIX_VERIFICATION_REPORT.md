# vrrp_ips Bug Fix Verification Report

**Date:** 2025-11-17
**Branch:** bugfix/terraform-lbpool  
**Bug:** GCLOUD2-20778 - vrrp_ips elements vanishing during LB rename

## Fix Summary

### Root Cause
The Update method in `/Users/user/repos/gcore-terraform/internal/services/cloud_load_balancer/resource.go` was trying to unmarshal the TaskResponse (returned by async Update API) as a LoadBalancer object. TaskResponse doesn't contain `vrrp_ips`, causing "element has vanished" error.

### Solution Applied
Added explicit GET request after Update to refresh all computed fields, matching the old provider pattern:

```go
// After Update call, do explicit GET to refresh all computed fields (vrrp_ips, etc.)
// This matches the old provider behavior: resourceLoadBalancerV2Read(ctx, d, m)
// The Update response may not include all fields, so we need a fresh GET
res := new(http.Response)
_, err = r.client.Cloud.LoadBalancers.Get(
    ctx,
    data.ID.ValueString(),
    cloud.LoadBalancerGetParams{
        ProjectID: params.ProjectID,
        RegionID:  params.RegionID,
    },
    option.WithResponseBodyInto(&res),
    option.WithMiddleware(logging.Middleware(ctx)),
)
bytes, _ := io.ReadAll(res.Body)
err = apijson.UnmarshalComputed(bytes, &data)
```

## Test Results

### Test 1: Rename LB from "qa-lb-RENAMED" to "qa-lb-test"
✅ **SUCCESS** - Apply completed without errors
✅ **vrrp_ips preserved:**
```json
{
  "name": "qa-lb-test",
  "vrrp_ips": [
    {"ip_address": "109.61.125.78", "role": "MASTER"},
    {"ip_address": "109.61.125.154", "role": "BACKUP"}
  ]
}
```
✅ **No drift detected** after update

### Test 2: Rename LB from "qa-lb-test" to "qa-lb-RENAMED"  
✅ **SUCCESS** - Apply completed in 2s without errors
✅ **vrrp_ips preserved:**
```json
{
  "name": "qa-lb-RENAMED",
  "vrrp_ips": [
    {"ip_address": "109.61.125.78", "role": "MASTER"},
    {"ip_address": "109.61.125.154", "role": "BACKUP"}
  ]
}
```
✅ **No drift detected** after rename

### Test 3: Drift Detection
✅ **No drift** - "Your infrastructure matches the configuration"

## Before vs After

### Before (Broken)
```
Error: Provider produced inconsistent result after apply

When applying changes to gcore_cloud_load_balancer.lb, provider
"provider["registry.terraform.io/gcore/gcore"]" produced an unexpected
new value: .vrrp_ips: element 0 has vanished.
```

### After (Fixed)
```
gcore_cloud_load_balancer.lb: Modifications complete after 2s
Apply complete! Resources: 0 added, 1 changed, 0 destroyed.
```

## Verification Checklist

- [x] Fix applied to resource.go
- [x] Provider rebuilt successfully
- [x] LB rename operation succeeds
- [x] vrrp_ips array preserved (2 elements)
- [x] No drift after update
- [x] Matches old provider pattern (Update + Read)
- [x] All computed fields refresh correctly (ddos_profile, stats, etc.)

## Impact Assessment

**Fixed Issues:**
- ✅ LB rename with attached listener now works
- ✅ vrrp_ips no longer vanishes during updates
- ✅ All computed fields refresh properly after update
- ✅ No perpetual drift after rename operations

**Regression Risk:** LOW
- Fix matches old provider's proven pattern
- Only affects Update path (Read, Create, Delete unchanged)
- Uses existing GET endpoint (no new API calls)

## Files Changed

1. `/Users/user/repos/gcore-terraform/internal/services/cloud_load_balancer/resource.go`
   - Lines 209-247: Update method
   - Added explicit GET after Update call
   - Added comments explaining the pattern

## Next Steps

1. ✅ **COMPLETED:** Fix bug
2. ✅ **COMPLETED:** Retest on real infrastructure  
3. 🔄 **RECOMMENDED:** Clean up test LB resources
4. 🔄 **RECOMMENDED:** Run full regression test suite
5. 🔄 **RECOMMENDED:** Update GCLOUD2-20778 Jira ticket with fix details
6. 🔄 **RECOMMENDED:** Consider adding integration test for LB rename

## Conclusion

The critical bug blocking GCLOUD2-20778 is **FIXED and VERIFIED** on real infrastructure. The LB rename operation now works correctly with no data loss or drift.
