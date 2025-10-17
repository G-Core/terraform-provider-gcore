# Reserved Fixed IP is_vip Update Fix - Complete

## Issue Summary

When setting `is_vip = true` for an existing `gcore_cloud_reserved_fixed_ip` resource, Terraform was forcing resource replacement instead of performing an in-place update using the patch method.

### Problem Details

**Symptom:**
```
# gcore_cloud_reserved_fixed_ip.external_reserved_fixed_ip must be replaced
-/+ resource "gcore_cloud_reserved_fixed_ip" "external_reserved_fixed_ip" {
      ~ port_id = "8c6bc5a5-d27a-458c-b370-f30a38ce7ad2" # forces replacement -> (known after apply)
      ~ is_vip  = false -> true
    }

Plan: 1 to add, 0 to change, 1 to destroy
```

**Expected Behavior:**
The resource should be updated in-place using the VIP toggle API, preserving the same resource ID and IP address.

## Root Cause Analysis

### Primary Issue: port_id Plan Modifier

The `port_id` field in the schema had `stringplanmodifier.RequiresReplace()` which was causing Terraform to force replacement when the field showed as "(known after apply)" during updates.

**Why this happened:**
1. `port_id` is defined as both `Computed: true` and `Optional: true`
2. It serves dual purposes:
   - Input parameter when `type = "port"` (user provides the port_id)
   - Resource identifier that equals `id` for other types (computed)
3. During an `is_vip` update, all computed fields show as "(known after apply)"
4. The `RequiresReplace()` modifier interprets this as a change, triggering replacement

### Secondary Issue: allowed_address_pairs Validation

The Update method was checking if `allowed_address_pairs` changed, but it's a computed-only field that shows as "(known after apply)" during updates, causing a false positive error.

## Solution Implemented

### 1. Modified Schema (schema.go)

Changed the `port_id` plan modifier from `RequiresReplace()` to `UseStateForUnknown()`:

```go
"port_id": schema.StringAttribute{
    Description:   "Port ID to make a reserved fixed IP (for example, `vip_port_id` of the Load Balancer entity).",
    Computed:      true,
    Optional:      true,
    PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
},
```

**Why this works:**
- `UseStateForUnknown()` tells Terraform to use the state value when the plan value is unknown
- This prevents `port_id` from showing as "(known after apply)" and triggering replacement
- The field can still be provided by users when creating resources with `type = "port"`

### 2. Enhanced Update Logic (resource.go)

Modified the `allowed_address_pairs` validation to only check when the value is known:

```go
// Only check if the plan value is known (not null/unknown) and different from state
if !plan.AllowedAddressPairs.IsNull() && !plan.AllowedAddressPairs.IsUnknown() && !plan.AllowedAddressPairs.Equal(state.AllowedAddressPairs) {
    resp.Diagnostics.AddError("Update Not Supported",
        "Updating 'allowed_address_pairs' is not supported yet. This feature requires Ports API integration which is not available in the current SDK version.")
    return
}
```

**Why this works:**
- Allows computed fields to show as "(known after apply)" without triggering validation errors
- Still prevents users from explicitly trying to modify the field

## Files Modified

1. **internal/services/cloud_reserved_fixed_ip/schema.go**
   - Changed line 81: `stringplanmodifier.RequiresReplace()` → `stringplanmodifier.UseStateForUnknown()`

2. **internal/services/cloud_reserved_fixed_ip/resource.go**
   - Updated line 170: Added `!plan.AllowedAddressPairs.IsNull() && !plan.AllowedAddressPairs.IsUnknown()` checks

## Testing Results

### Test 1: Update is_vip from false to true
```bash
# Before fix:
Plan: 1 to add, 0 to change, 1 to destroy (REPLACEMENT)

# After fix:
Plan: 0 to add, 1 to change, 0 to destroy (IN-PLACE UPDATE)
```

**Verification:**
- ✅ Resource ID remained unchanged: `296762b8-bac8-4b1f-a121-5b466d43662a`
- ✅ IP address preserved: `193.57.88.228`
- ✅ `is_vip` successfully updated from `false` to `true`
- ✅ Update completed using VIP toggle API

### Test 2: Update is_vip from true to false (reverse)
```bash
Plan: 0 to add, 1 to change, 0 to destroy (IN-PLACE UPDATE)
```

**Verification:**
- ✅ Resource ID remained unchanged: `296762b8-bac8-4b1f-a121-5b466d43662a`
- ✅ IP address preserved: `193.57.88.228`
- ✅ `is_vip` successfully updated from `true` to `false`
- ✅ Update completed using VIP toggle API

## Test Logs

All test logs are saved in `test-reserved-fixed-ip/`:
- `create_with_is_vip_false.log` - Initial resource creation
- `plan_is_vip_true.log` - Plan output showing the issue (before fix)
- `plan_after_schema_fix.log` - Plan output showing in-place update (after fix)
- `apply_is_vip_true_v2.log` - Successful apply changing is_vip to true
- `apply_is_vip_false.log` - Successful reverse update back to false

## API Integration

The fix leverages the existing `Vip.Toggle()` API implementation:

```go
apiRes, err := r.client.Cloud.ReservedFixedIPs.Vip.Toggle(
    ctx,
    resourceID,
    params,
    option.WithMiddleware(logging.Middleware(ctx)),
)
```

This API endpoint:
- Updates the VIP status without replacing the resource
- Returns the updated reserved fixed IP object
- Maintains the same resource ID and IP address

## Migration Notes

### For Generated Code
Both modified files (`schema.go` and `resource.go`) are marked as generated by Stainless. When regenerating:

1. **schema.go**: The `port_id` plan modifier change should be configured in the OpenAPI spec or applied as a post-generation patch
2. **resource.go**: The Update method logic is custom and should be preserved

### For Users
This fix is **backward compatible**:
- Existing resources will not be affected
- No state migration required
- Users can now update `is_vip` in-place without resource recreation

## Conclusion

The fix successfully resolves the `is_vip` update issue by:
1. Removing the incorrect `RequiresReplace()` modifier from `port_id`
2. Properly handling computed fields in the Update method
3. Preserving the existing VIP toggle API functionality

**Result:** Users can now toggle `is_vip` on reserved fixed IPs without forcing resource replacement, using the proper PATCH API method as intended.
