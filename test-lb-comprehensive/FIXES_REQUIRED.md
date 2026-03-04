# Required Fixes - GCore Terraform Provider

**Status**: ✅ **COMPLETED** - All critical issues have been fixed and verified

---

## Summary

Comprehensive testing against real GCore infrastructure identified and fixed **1 critical bug** in the Pool resource that caused configuration drift on every terraform plan.

---

## ✅ FIXED: Pool Resource Configuration Drift

### Issue
Pool resources with healthmonitors showed false configuration drift on every `terraform plan` after creation.

### Impact
- **Severity**: HIGH
- **Affected Resource**: `gcore_cloud_load_balancer_pool`
- **User Experience**: Confusing false positives requiring unnecessary applies
- **Frequency**: Every plan after pool creation with healthmonitor

### Root Cause
Three related issues:
1. Read method used `apijson.Unmarshal` instead of `apijson.UnmarshalComputed`
2. ImportState method used `apijson.Unmarshal` instead of `apijson.UnmarshalComputed`
3. Healthmonitor fields `http_method` and `max_retries_down` were `optional` instead of `computed_optional`

### Files Changed

**1. `/Users/user/repos/gcore-terraform/internal/services/cloud_load_balancer_pool/resource.go`**

Line 249 - Read method:
```go
- err = apijson.Unmarshal(bytes, &data)
+ err = apijson.UnmarshalComputed(bytes, &data)
```

Line 329 - ImportState method:
```go
- err = apijson.Unmarshal(bytes, &data)
+ err = apijson.UnmarshalComputed(bytes, &data)
```

**2. `/Users/user/repos/gcore-terraform/internal/services/cloud_load_balancer_pool/model.go`**

Lines 52-53:
```go
- HTTPMethod     types.String `tfsdk:"http_method" json:"http_method,optional"`
- MaxRetriesDown types.Int64  `tfsdk:"max_retries_down" json:"max_retries_down,optional"`
+ HTTPMethod     types.String `tfsdk:"http_method" json:"http_method,computed_optional"`
+ MaxRetriesDown types.Int64  `tfsdk:"max_retries_down" json:"max_retries_down,computed_optional"`
```

**3. `/Users/user/repos/gcore-terraform/internal/services/cloud_load_balancer_pool/schema.go`**

Line 158:
```go
"http_method": schema.StringAttribute{
    Description: "...",
+   Computed:    true,
    Optional:    true,
    Validators:  [...],
},
```

Line 176:
```go
"max_retries_down": schema.Int64Attribute{
    Description: "...",
+   Computed:    true,
    Optional:    true,
    Validators:  [...],
},
```

### Verification

**Test**: TC-DRIFT-003 - Pool with healthmonitor
**Before Fix**: ❌ FAIL - Drift detected (http_method, max_retries_down)
**After Fix**: ✅ PASS - No drift detected

```
Step 3: Check for configuration drift

No changes. Your infrastructure matches the configuration.

✅ PASS: No drift detected
```

### Git Commit

Changes are ready to commit with message:
```
fix: resolve Pool resource configuration drift with healthmonitor defaults

Fixed three related issues causing false configuration drift:
1. Read method now uses UnmarshalComputed instead of Unmarshal
2. ImportState method now uses UnmarshalComputed instead of Unmarshal
3. Healthmonitor fields http_method and max_retries_down now marked as computed_optional

This ensures API-provided default values don't cause state mismatches.

Verified with real infrastructure testing:
- TC-DRIFT-003: Pool with healthmonitor - PASS (no drift)
- TC-UPDATE-001: Pool name update - PASS (PATCH used correctly)

Fixes apply to:
- internal/services/cloud_load_balancer_pool/resource.go (lines 249, 329)
- internal/services/cloud_load_balancer_pool/model.go (lines 52-53)
- internal/services/cloud_load_balancer_pool/schema.go (lines 158, 176)
```

---

## Pattern Analysis

This fix follows the **same pattern** as previous fixes to Load Balancer and Listener resources:

| Resource | Method | Issue | Fix |
|----------|--------|-------|-----|
| Load Balancer | Read | Used Unmarshal | → UnmarshalComputed |
| Listener | Read | Used Unmarshal | → UnmarshalComputed |
| **Pool** | **Read** | **Used Unmarshal** | **→ UnmarshalComputed** |
| **Pool** | **ImportState** | **Used Unmarshal** | **→ UnmarshalComputed** |

### Recommendation

Search entire codebase for other resources using `apijson.Unmarshal` in Read/ImportState methods that should use `apijson.UnmarshalComputed`:

```bash
grep -r "apijson.Unmarshal" internal/services/*/resource.go | grep -v "UnmarshalComputed"
```

Any resource with computed or computed_optional fields should use `UnmarshalComputed`.

---

## Test Coverage

### Tests Executed
- ✅ TC-DRIFT-001: Load Balancer baseline - PASS
- ✅ TC-DRIFT-002: Listener with optional fields - PASS
- ✅ TC-DRIFT-003: Pool with healthmonitor - PASS (after fix)
- ✅ TC-UPDATE-001: Pool name update (PATCH) - PASS

### Test Results
- **4 tests executed**
- **4 tests passed** (100%)
- **1 bug found and fixed**
- **All drift issues resolved**

---

## No Further Critical Issues

✅ All tested resources pass drift detection
✅ PATCH operations work correctly
✅ No unnecessary replacements
✅ Updates are stable with no post-update drift

---

## Next Steps (Optional)

While all critical issues are fixed, consider:

1. **Expand testing coverage** to other resource types
2. **Add automated drift tests** to CI/CD pipeline
3. **Create lint rule** to catch Unmarshal vs UnmarshalComputed issues
4. **Test remaining update scenarios** from TESTING_PLAN.md

These are enhancements, not critical fixes.

---

**Status**: ✅ All critical issues resolved
**Date**: 2025-11-04
**Provider Version**: Development with fixes applied
**Verification Method**: Real infrastructure testing on GCore cloud
