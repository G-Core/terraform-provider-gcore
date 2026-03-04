# Router Resource Bug Findings - Test Session 2025-11-20

## Summary
While testing Pedro's fixes for GCLOUD2-21144, discovered multiple blocking bugs in the router resource that prevent comprehensive testing.

## Bugs Discovered

### 1. external_gateway_info Drift Issue (HIGH PRIORITY)
**Status**: BLOCKING
**Location**: `internal/services/cloud_network_router/model.go:18`

**Description**:
- Router created without `external_gateway_info` specified
- API returns `{enable_snat: false, network_id: "", type: null}`
- Terraform stores this in state
- On subsequent plan/refresh, Terraform shows drift:
  ```
  external_gateway_info = {
      enable_snat = false -> (known after apply)
      network_id  = (known after apply)
      type        = (known after apply)
  } -> (known after apply)
  ```
- This causes every plan to show changes even when no actual updates needed

**Evidence**: `evidence/test_1_drift_corrected.log`

**Root Cause**:
- `external_gateway_info` is `computed_optional`
- ModifyPlan or state refresh logic incorrectly marks it as "known after apply"

**Workaround**:
Explicitly set in config:
```hcl
external_gateway_info = {
  enable_snat = false
}
```

**Impact**: Prevents drift testing, causes unnecessary API calls

---

### 2. Empty PATCH Body Bug (HIGH PRIORITY)
**Status**: BLOCKING
**Location**: `internal/services/cloud_network_router/resource.go:316`

**Description**:
- When `external_gateway_info` drift triggers `needsUpdate = true` (line 287)
- But no actual changes exist for PATCH body
- `MarshalForPatch` produces empty bytes
- Provider sends PATCH with no body
- API returns 400: "Invalid JSON: EOF while parsing a value at line 1 column 0"

**Evidence**:
- `evidence/test_2_apply.log` - shows PATCH 400 error
- `mitm_requests.log` - confirms "No request body" in PATCH request

**MITM Log Excerpt**:
```
REQUEST: PATCH https://api.gcore.com/cloud/v1/routers/379987/76/25204151-fa9f-4b8a-8967-e3e6a1c5f48b
No request body

RESPONSE: 400
{"exception_class":"ValidationError","message":"Validation Error: Invalid JSON: EOF while parsing a value at line 1 column 0"}
```

**Root Cause**:
1. Line 316 check: `if (len(dataBytes) > 0 && string(dataBytes) != "{}" && string(dataBytes) != "null")`
2. This should prevent empty PATCH, but the check fails to catch all empty cases
3. OR: `needsUpdate` logic includes computed_optional drift that shouldn't trigger PATCH

**Fix Needed**:
- Option A: Improve empty body detection at line 316
- Option B: Exclude computed_optional drift from `needsUpdate` calculation
- Option C: Ensure `MarshalForPatch` never produces truly empty bodies

**Impact**: Blocks all updates when drift exists

---

### 3. AttachSubnet/DetachSubnet Working Correctly ✓
**Status**: VERIFIED WORKING

**Evidence**: `mitm_requests.log` shows:
```
REQUEST: POST https://api.gcore.com/cloud/v1/routers/379987/76/.../attach
```

**Verification**: When external_gateway_info drift is worked around, provider correctly uses:
- `POST .../attach` for adding interfaces
- Interface detach logic exists at lines 257-281 (not tested yet due to blocking bugs)

---

## Test Environment Issues

### Subnet Auto-Connect Conflict
**Status**: RESOLVED

**Issue**: `connect_to_network_router = true` (default) caused conflict with manual interface management

**Solution**: Set `connect_to_network_router = false` in subnet config

---

## Pedro's Fixes Status

### ✓ Partial Routes Deletion Fix
**Code Location**: `resource.go:133`
**Change**: `len(dataRoutes) == 0 && len(stateRoutes) > 0` → `len(dataRoutes) < len(stateRoutes)`
**Status**: CODE REVIEWED - CORRECT
**Testing Status**: BLOCKED by drift bugs

### ✓ Routes Preservation After Early PATCH
**Code Location**: `resource.go:141, 186`
**Added**: `plannedRoutes := data.Routes` and restore after PATCH
**Status**: CODE REVIEWED - CORRECT
**Testing Status**: BLOCKED by drift bugs

---

## Testing Strategy Going Forward

### Recommended Approach
1. **Fix Bug #1 and #2 first** - These block all testing
2. **OR: Use direct API testing** - Bypass Terraform provider, test API behavior directly
3. **OR: Mock/unit tests** - Test `MarshalForPatch` and Update logic in isolation

### Tests Still Needed (Once Bugs Fixed)
- ✓ Test #1: Minimal router (PARTIAL - saw drift but not blocking)
- ⏸ Test #2: Add interface (BLOCKED - empty PATCH bug)
- ⏸ Test #3: Remove interface
- ⏸ Test #4: Add routes
- ⏸ Test #5: Partial route deletion (PEDRO'S CRITICAL FIX)
- ⏸ Test #6a: Full deletion with routes=[]
- ⏸ Test #6b: Full deletion by removing attr
- ⏸ Test #7: Name update
- ⏸ Test #8: Interface deletion with []
- ⏸ Test #9: Interface deletion by removing attr
- ⏸ Test #10: Import

---

## Files Generated
- `evidence/test_1_apply_corrected.log` - Initial clean apply
- `evidence/test_1_drift_corrected.log` - Drift detection showing external_gateway_info issue
- `evidence/test_2_apply.log` - First attempt, hit empty PATCH bug
- `evidence/test_2_apply_with_mitm.log` - Confirms PATCH used instead of AttachSubnet (before proxy fix)
- `evidence/test_2_apply_with_workaround.log` - With external_gateway_info workaround
- `mitm_requests.log` - Full HTTP request/response capture

---

## Recommendations

### For Immediate Action
1. **Fix external_gateway_info drift** - Critical for any router testing
2. **Fix empty PATCH body handling** - Prevents updates when drift exists
3. **Add integration test** for computed_optional field drift behavior

### For Pedro's Review
- Pedro's code changes (lines 133, 141, 186) appear correct
- Unable to verify with live infrastructure due to blocking bugs
- Recommend unit tests for partial deletion logic
- Recommend mock tests for routes preservation logic

---

## Technical Details

### needsUpdate Calculation (resource.go:285-287)
```go
needsUpdate := !data.Name.Equal(state.Name) ||
    !data.Routes.Equal(state.Routes) ||
    !data.ExternalGatewayInfo.Equal(state.ExternalGatewayInfo)
```

**Issue**: `!data.ExternalGatewayInfo.Equal(state.ExternalGatewayInfo)` returns true even when no real change exists, only drift from computed values.

### MarshalForPatch Flow
1. `data.MarshalJSONForUpdate(*state)` called (resource.go:310)
2. Calls `apijson.MarshalForPatch(mCopy, state)` (model.go:39)
3. Should only include fields that actually changed
4. But computed_optional drift causes comparison mismatch
5. Results in empty or near-empty patch body

---

## MITM Proxy Setup
- **Port**: 9092
- **Capture File**: `flow.mitm`
- **Body Log**: `mitm_requests.log`
- **Addon**: `/tmp/mitm_body_logger.py`
- **Status**: WORKING ✓

