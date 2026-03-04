# Pedro's Router Review - Test Report
**Date**: 2025-11-20
**JIRA**: GCLOUD2-21144
**PR**: https://github.com/stainless-sdks/gcore-terraform/pull/38

## Executive Summary

Attempted comprehensive testing of Pedro's review comments and code fixes. **Discovered blocking bugs preventing full infrastructure testing**. Pedro's code changes appear correct upon code review, but live infrastructure verification blocked.

---

## Pedro's Review Items - Status

### 1. ✅ Remove Excessive Comments
**Pedro's Feedback**: "comments about comments in code - just do suggested edit"

**Action Taken**:
- Removed IMPORTANT/CRITICAL labels from comments
- Removed verbose comment block at model.go:40-45
- Cleaned up excessive documentation

**Status**: COMPLETED

**Files Modified**:
- `internal/services/cloud_network_router/model.go`
- `internal/services/cloud_network_router/resource.go`

---

### 2. ⚠️  Partial Routes Deletion Bug Fix
**Pedro's Feedback**:
> "I think there is an issue here. You currently only delete all routes if there are no routes in the plan and routes in the state. However, if you have 3 routes in the state and 1 route in the plan, you also need to delete the routes because the single PATCH will not delete routes, it will just replace them."

**Pedro's Suggested Fix**:
```go
routesDeletionNeeded := len(dataRoutes) < len(stateRoutes)
```

**Action Taken**:
- Changed condition from `len(dataRoutes) == 0 && len(stateRoutes) > 0`
- To `len(dataRoutes) < len(stateRoutes)`
- Location: `resource.go:133`

**Code Review Assessment**: ✅ CORRECT
- Old code only caught full deletion (n→0)
- New code catches partial deletion (n→m where m<n)
- Logic is sound

**Infrastructure Testing**: ❌ BLOCKED
- Cannot test due to external_gateway_info drift bug
- Test #5 (partial deletion 3→1) designed but not executed
- Recommend unit/mock tests as alternative

---

### 3. ⚠️ Routes Preservation After Early PATCH
**Related to Item #2**

**Issue Identified**: When routes deleted before interface detachment, planned routes value was lost

**Fix Applied**:
```go
// Line 141: Save planned routes
plannedRoutes := data.Routes

// ... early PATCH to delete routes ...

// Line 186: Restore planned routes
data.Routes = plannedRoutes
```

**Code Review Assessment**: ✅ CORRECT
- Preserves user's planned routes after early deletion PATCH
- Prevents routes from being lost when only interfaces updating
- Matches pattern used for interfaces preservation

**Infrastructure Testing**: ❌ BLOCKED

---

### 4. ❓ Verify Empty Slice Serialization Without WithJSONSet
**Pedro's Question**:
> "It looks like you're adding a with json set with null to delete the routes. I just wanted to make sure that if this is needed. If you have routes on the state, and a user removes the routes field, the provider would populate the routes plan value with an empty slice due to computed_optional. To serialize the patch, this empty slice would normally not be serialized and not appear in the patch json due to the omit zero tag on the attr. Check this behavior on practice using real calls on infra."

**Testing Plan Created**:
- Test #6a: Set `routes = []` explicitly in config
- Test #6b: Remove `routes` attr completely from config
- Both should verify if PATCH contains `{"routes": []}`

**Status**: ❌ NOT TESTED - BLOCKED
- Requires working provider to test
- MITM logging setup and working
- Would capture PATCH body to verify serialization

**Recommendation**:
- Test with SDK unit tests
- Check `apijson.MarshalForPatch` behavior with empty computed_optional slices
- May need Stainless framework investigation

---

### 5. ❓ Simplify Update Logic - Single PATCH?
**Pedro's Comment**:
> "I wonder if you could simplify your implementation. Looking at your update flow, you have a PATCH, an attach operation, a detach operation and then another PATCH. Double check it, can we make our implementation simpler? Verify with real infra example."

**Current Flow Analysis** (resource.go:111-344):
1. **Optional early PATCH** (lines 139-187): Only if routes deleted + interfaces changing
2. **Attach/Detach operations** (lines 190-283): Individual POST calls per interface
3. **Main PATCH** (lines 292-343): For name, routes, external_gateway_info changes
4. **Final GET** (lines 346-371): Refresh state

**Assessment**:
- ✅ Attach/Detach cannot be in PATCH (API doesn't support it)
- ✅ Early PATCH necessary due to API constraint: "routes using an interface must be deleted first"
- ✅ Multiple PATCHes needed only in specific case (routes deletion + interface detachment)
- ✅ Most common updates use single PATCH

**Conclusion**: Implementation is already optimized. Cannot be simplified further due to API constraints.

**Verification**: ⚠️ PARTIAL
- MITM logs confirm AttachSubnet uses `POST .../attach` ✓
- Cannot verify full flow due to blocking bugs

---

### 6. ❓ Is Final GET Necessary?
**Pedro's Comment**:
> "If you're doing a get in read, why do you need to do a GET after update/attach/detach?"

**Current Behavior**: Final GET at lines 346-371

**Assessment**:
- GET is necessary because:
  1. AttachSubnet/DetachSubnet responses don't include full router state
  2. Need to populate computed fields (updated_at, task_id, etc.)
  3. Ensures state matches actual API state after multiple operations

**Recommendation**: KEEP the final GET
- Standard Terraform provider pattern
- Ensures state consistency
- Minimal performance impact (single GET vs potential drift)

**Could Optimize**:
- Skip GET if only name updated (PATCH response sufficient)
- But adds complexity for minor gain

---

## Blocking Bugs Discovered

### Bug #1: external_gateway_info Drift
**Severity**: HIGH - BLOCKS ALL TESTING

**Symptom**: Every `terraform plan` shows changes to router even when no actual changes needed

**Evidence**: `evidence/test_1_drift_corrected.log`

**Workaround**: Explicitly set `external_gateway_info = { enable_snat = false }` in config

**Needs**: Provider fix for computed_optional field handling

---

### Bug #2: Empty PATCH Body
**Severity**: HIGH - BLOCKS ALL UPDATES

**Symptom**: When drift triggers update, provider sends PATCH with no body, gets 400 error

**Evidence**:
- `evidence/test_2_apply.log`
- `mitm_requests.log` shows "No request body"

**Root Cause**: `needsUpdate` includes computed_optional drift, but `MarshalForPatch` produces no actual changes

**Needs**: Fix at resource.go:316 or needsUpdate calculation

See `BUG_FINDINGS.md` for detailed analysis.

---

## What Was Successfully Tested

### ✅ Provider Build
- Built with Pedro's fixes
- No compilation errors
- Provider loads correctly

### ✅ MITM Logging
- Setup working on port 9092
- Capturing request/response bodies
- Confirmed AttachSubnet uses correct endpoint

### ✅ Code Review
- All Pedro's suggested changes reviewed
- Logic verified correct
- No obvious issues in code

### ✅ Basic Resource Creation
- Can create network, subnet, router
- Resources created successfully (when no drift present)

---

## What Could Not Be Tested

### ❌ Drift Detection
- external_gateway_info bug prevents clean drift testing

### ❌ Interface Operations
- Add interface: blocked by empty PATCH bug
- Remove interface: not reached
- Both attachment methods: not tested

### ❌ Routes Operations
- Add routes: blocked
- Partial deletion (Pedro's fix): blocked
- Full deletion with routes=[]: blocked
- Full deletion by removing attr: blocked

### ❌ Combined Operations
- Routes + Interfaces changes: blocked
- All complex scenarios: blocked

---

## Recommendations

### Immediate Actions Needed
1. **Fix external_gateway_info drift bug**
   - Root cause in computed_optional handling
   - May need Stainless framework investigation

2. **Fix empty PATCH body bug**
   - Improve empty body detection OR
   - Exclude computed drift from needsUpdate

3. **Add Unit Tests**
   - Test partial deletion logic in isolation
   - Test MarshalForPatch with various state transitions
   - Test empty computed_optional slice serialization

### For Pedro's Review
**Approval Recommendation**: ✅ APPROVE WITH CAVEATS

Pedro's code changes are correct:
- ✅ Partial deletion fix is the right solution
- ✅ Routes preservation fix prevents data loss
- ✅ Comment cleanup improves readability
- ✅ Update flow is already optimized for API constraints
- ✅ Final GET is necessary

**Caveats**:
- Cannot verify with live infrastructure due to pre-existing bugs
- Need unit tests for routes deletion logic
- Need investigation into empty slice serialization (Test #6a/6b)
- Provider has other bugs that should be fixed

### Testing Strategy Forward

**Option A: Fix Bugs First** (Recommended)
1. Fix external_gateway_info drift
2. Fix empty PATCH body
3. Re-run full test suite

**Option B: Unit Testing**
1. Mock Terraform state objects
2. Test Update function logic directly
3. Verify MarshalForPatch output

**Option C: Direct API Testing**
1. Use curl/SDK to test API behavior
2. Verify routes deletion requirements
3. Bypass provider bugs entirely

---

## Test Artifacts Generated

### Evidence Logs
- `evidence/test_1_apply_corrected.log` - Clean resource creation
- `evidence/test_1_drift_corrected.log` - Drift detection showing Bug #1
- `evidence/test_2_apply.log` - Shows Bug #2 (empty PATCH)
- `evidence/test_2_apply_with_mitm.log` - Early testing
- `evidence/test_2_apply_with_workaround.log` - With workaround
- `evidence/test_cleanup_for_restart.log` - Resource cleanup

### Analysis Documents
- `BUG_FINDINGS.md` - Detailed bug analysis
- `TESTING_PLAN.md` - Original 10-test comprehensive plan
- `PEDRO_REVIEW_TEST_REPORT.md` - This report

### MITM Captures
- `mitm_requests.log` - HTTP request/response bodies
- `flow.mitm` - Full mitmproxy capture
- Confirms AttachSubnet endpoint usage ✓

---

## Conclusion

Pedro's review identified a real bug (partial deletion) and his suggested fix is correct. The code changes are sound and ready to merge. However, the provider has pre-existing bugs that prevent comprehensive infrastructure testing.

**Recommendation**:
- ✅ Approve Pedro's changes
- ⚠️  Add unit tests before merge
- 🐛 Create separate tickets for Bug #1 and Bug #2
- 📋 Follow up with tests once bugs fixed

The partial deletion bug fix addresses an edge case that users would encounter, and the fix is technically correct per code review.

