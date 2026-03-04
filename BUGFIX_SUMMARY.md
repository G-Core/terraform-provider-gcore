# GCLOUD2-20778 Critical Bug Fix - COMPLETED ✅

## Summary
**Fixed critical bug preventing Load Balancer rename operations**

### The Problem
When renaming a Load Balancer with an attached listener, the operation failed with:
```
Error: Provider produced inconsistent result after apply
.vrrp_ips: element 0 has vanished.
```

### Root Cause
File: `/Users/user/repos/gcore-terraform/internal/services/cloud_load_balancer/resource.go`

The Update method was attempting to unmarshal the async TaskResponse as a LoadBalancer object. Since TaskResponse doesn't contain the `vrrp_ips` field, Terraform detected data loss and aborted.

### The Fix
**Changed:** Lines 209-247 in `resource.go`

**Pattern:** Added explicit GET request after Update to refresh all computed fields, matching the proven pattern from old provider:

```go
// Call Update - SDK handles task polling internally
_, err = r.client.Cloud.LoadBalancers.Update(ctx, data.ID.ValueString(), params, ...)

// After update, do explicit GET to refresh all computed fields
res := new(http.Response)
_, err = r.client.Cloud.LoadBalancers.Get(ctx, data.ID.ValueString(), getParams,
    option.WithResponseBodyInto(&res),
    option.WithMiddleware(logging.Middleware(ctx)),
)
bytes, _ := io.ReadAll(res.Body)
err = apijson.UnmarshalComputed(bytes, &data)
```

### Test Results (Real Infrastructure)
✅ **Test 1:** Renamed LB from "qa-lb-RENAMED" → "qa-lb-test" - SUCCESS  
✅ **Test 2:** Renamed LB from "qa-lb-test" → "qa-lb-RENAMED" - SUCCESS  
✅ **Test 3:** Drift detection - NO DRIFT  
✅ **Test 4:** vrrp_ips preserved (2 elements: MASTER + BACKUP)  
✅ **Test 5:** All computed fields refresh correctly  

### Issues from GCLOUD2-20778 (Status Update)

1. ✅ **FIXED:** Pool creation timing - NewAndPoll works correctly
2. ✅ **FIXED:** Listener drift with timeout fields - UseStateForUnknown modifiers present
3. ✅ **FIXED:** LB rename drift - **THIS BUG** (was worse than drift - now completely fixed)
4. ⚠️ **SCHEMA ISSUE:** Tag removal - tags_v2 is read-only (separate issue, not critical)

### Files Changed
- `/Users/user/repos/gcore-terraform/internal/services/cloud_load_balancer/resource.go` (Update method)

### Reports Generated
- `/Users/user/repos/gcore-terraform/test-lb-full/FIX_VERIFICATION_REPORT.md` - Detailed fix verification
- `/Users/user/repos/gcore-terraform/LB_COMPREHENSIVE_TEST_REPORT.md` - Initial bug discovery report

### Regression Risk
**LOW** - Fix uses proven pattern from old provider, only affects Update path

### Ready for
- ✅ Code review
- ✅ Merge to main
- ✅ Jira ticket update
- ✅ Release

---
**Tested on:** Real Gcore infrastructure (Region: Luxembourg-2, Project: 379987)  
**Date:** 2025-11-17  
**Branch:** bugfix/terraform-lbpool
