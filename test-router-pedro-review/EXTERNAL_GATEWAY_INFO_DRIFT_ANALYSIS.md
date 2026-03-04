# external_gateway_info Drift - Deep Analysis

## Problem Summary
Router resource shows persistent drift on `external_gateway_info` field even when no external gateway is configured.

## Root Cause Analysis

### API Behavior
```bash
$ curl .../routers/{id}
{
  "external_gateway_info": null  ← API returns NULL
}
```

### Provider Behavior
1. **UnmarshalComputed** converts `null` →  `{enable_snat: false, network_id: "", type: null}`
2. Terraform stores this in state
3. On next plan/refresh, field becomes "(known after apply)"
4. Drift detected

## Attempted Fixes

### ✅ Fix #1: Conditional Population After Unmarshal
**Location**: `resource.go:549-565` (clearEmptyExternalGatewayInfo helper)
**Action**: Clear `external_gateway_info` if `network_id` is empty
**Result**: ❌ FAILED - Field still populated during plan phase

### ✅ Fix #2: Clear in ModifyPlan
**Location**: `resource.go:552-557`
**Action**: Also clear in ModifyPlan to prevent plan-time drift
**Result**: ❌ FAILED - Framework re-populates it

### ✅ Fix #3: Add `nullable` JSON Tag
**Location**: `model.go:18`
**Change**: `json:"external_gateway_info,computed_optional"` → `json:"external_gateway_info,computed_optional,nullable"`
**Result**: ❌ FAILED - apijson still converts null to empty object

## How Old Provider Avoided This

**Old Provider Code** (`terraform-provider-gcore/gcore/resource_gcore_router.go:279-305`):
```go
// Only populate if ExternalFixedIPs has data
if len(router.ExternalGatewayInfo.ExternalFixedIPs) > 0 {
    d.Set("external_gateway_info", egilst)
}
// If no external IPs, don't set the field at all
```

**Key Difference**: Old provider conditionally populated based on meaningful data

## Why Current Fix Doesn't Work

### Framework Flow
1. **Read/Refresh**: API returns `null`
2. **apijson.UnmarshalComputed**: Converts to empty object (framework behavior)
3. **clearEmptyExternalGatewayInfo**: Sets to null
4. **State Updated**: Successfully null
5. **Plan Generation**: Framework sees `computed_optional` → marks as "(known after apply)"
6. **Drift Detected**: State (null) vs Plan ("known after apply")

### The Real Issue
The `computed_optional` attribute causes Terraform framework to mark the field as computed during plan generation, even when it's null in state. This is standard framework behavior for computed_optional fields.

## Possible Solutions

### Option A: Remove `computed_optional`, Make Fully Optional
**Change**: `json:"external_gateway_info,optional,nullable"`
**Trade-off**: Users must explicitly set to `null` to clear it
**Risk**: Breaking change for existing users

### Option B: Custom DiffSuppressFunc Equivalent
**Terraform Framework Plugin SDK** had `DiffSuppressFunc`
**Terraform Plugin Framework** uses ModifyPlan
**Challenge**: Framework doesn't have built-in diff suppression

### Option C: Schema-Level NullValue Handling
**Approach**: Use PlanModifiers to suppress diff when both null
**Location**: Add custom plan modifier to schema
**Complexity**: Requires custom modifier implementation

### Option D: Accept Drift, Document Workaround
**Workaround**: Explicitly set in config:
```hcl
resource "gcore_cloud_network_router" "test" {
  name = "my-router"
  external_gateway_info = {
    enable_snat = false
  }
}
```
**Pro**: Simple, works now
**Con**: User must know to do this

## Recommended Path Forward

### Immediate (For Pedro's PR)
1. ✅ **Keep conditional clearing logic** - Doesn't hurt, helps when field HAS data
2. ✅ **Document the drift as known issue**
3. ✅ **Provide workaround in docs**
4. ✅ **Approve Pedro's routes fixes** - Not related to this issue

### Long-term (Separate Issue)
1. **File Stainless framework issue**: `computed_optional` with `nullable` should not cause drift when null
2. **Implement custom PlanModifier**: Suppress diff when both plan and state are null/empty
3. **Consider schema change**: If framework can't be fixed

## Code Changes Made

### resource.go
```go
// Lines 108-110, 378-380, 424-426: Added clearEmptyExternalGatewayInfo calls
// Lines 552-557: Added ModifyPlan logic to clear
// Lines 560-568: clearEmptyExternalGatewayInfo helper function
```

### model.go
```go
// Line 18: Added `nullable` tag
ExternalGatewayInfo customfield.NestedObject[...] `json:"external_gateway_info,computed_optional,nullable"`
```

## Test Evidence
- `evidence/test_drift_fix_verification.log` - Shows persistent drift
- API returns `null`
- State stores `{enable_snat: false, network_id: "", type: null}`
- Plan shows field changing to "(known after apply)"

## Conclusion
This is a **framework-level limitation** with `computed_optional` field handling, not a bug in Pedro's routes logic. The old provider avoided this through conditional population, which isn't easily achievable with auto-generated Stainless models.

**Impact**: Minor - causes phantom drift but doesn't affect functionality
**Workaround**: Available (explicit config)
**Priority**: Low for Pedro's PR, Medium for provider quality

