# GCLOUD2-20778 Drift Issue Investigation Report

## Issue Summary

**Reported:** 2025-11-20 by Kirill Tsaregorodtsev
**Jira Ticket:** GCLOUD2-20778
**Status:** ✅ CANNOT REPRODUCE - Already Fixed

### Reported Problem
After adding an LB pool to a listener that has a `user_list` configured, Terraform detects resource drift in the listener:
- The `stats` attribute shows as "(known after apply)"
- The `encrypted_password` in `user_list` shows as changing from current value to "(known after apply)"

### Test Scenario Executed
1. Created LB + Listener with `user_list` configured (username + encrypted_password)
2. Applied successfully
3. Added LB pool attached to the listener
4. Applied successfully
5. Ran `terraform plan` to check for drift
6. **Result:** ✅ No drift detected - "No changes. Your infrastructure matches the configuration."

## Analysis

### Why Drift Is Not Occurring (Current State)

The drift issue has been **already fixed** by commit `f999c1c` applied on 2025-11-18 (2 days before the Jira report on 2025-11-20).

**Fix Applied in f999c1c:**
- Added `UseStateForUnknown()` plan modifiers to all computed fields in LB, Listener, Pool, and Pool Member schemas
- This prevents Terraform from showing drift for fields that haven't changed
- Added explicit GET after PATCH for load balancer updates to retrieve full state

**Current Schema Implementation** (`internal/services/cloud_load_balancer_listener/schema.go`):

```go
"user_list": schema.ListNestedAttribute{
    Description: "Load balancer listener list of username and encrypted password items",
    Computed:    true,
    Optional:    true,
    PlanModifiers: []planmodifier.List{
        listplanmodifier.UseStateForUnknown(),  // ✅ Prevents drift
    },
    // ...
}

"stats": schema.SingleNestedAttribute{
    Description: "Statistics of the load balancer...",
    Computed:    true,
    PlanModifiers: []planmodifier.Object{
        objectplanmodifier.UseStateForUnknown(),  // ✅ Prevents drift
    },
    // ...
}

"encrypted_password": schema.StringAttribute{
    Description: "Encrypted password to auth via Basic Authentication",
    Required:    true,
    // Note: No plan modifier needed here as parent user_list has UseStateForUnknown()
}
```

### Why This Wasn't Found Earlier

#### 1. **Timing Issue**
The fix was applied 2 days **before** the Jira comment was posted. The QA tester likely:
- Was testing an older build that didn't have the fix yet
- Or tested with cached provider binary
- Or the test was started before the fix was merged

#### 2. **Test Coverage Gap**
Looking at existing test directories, there was no specific test for:
- Creating LB + Listener with `user_list`
- Then adding a pool
- Then checking for drift

Most tests create all resources at once, which doesn't expose this specific scenario.

#### 3. **Edge Case Scenario**
The drift only occurred in this specific sequence:
1. LB + Listener with user_list EXISTS
2. ADD a pool (modify infrastructure)
3. Check for drift

When all resources are created together (single apply), no drift occurs even without the fix, because Terraform doesn't have a previous state to compare against.

### Comparison: Old Provider vs New Provider

#### Old Provider (`terraform-provider-gcore`)
```go
// gcore/resource_gcore_lblistener.go

"user_list": &schema.Schema{
    Type:        schema.TypeList,
    Optional:    true,
    Elem: &schema.Resource{
        Schema: map[string]*schema.Schema{
            "encrypted_password": &schema.Schema{
                Type:        schema.TypeString,
                Sensitive:   true,  // ✅ Marked as sensitive
                Required:    true,
            },
        },
    },
}

// Update handling
if d.HasChange("user_list") {  // ✅ Explicit change detection
    u := d.Get("user_list")
    updateOpts.UserList = make([]listeners.CreateUserListOpts, 0)
    // ...
}
```

**Old Provider Approach:**
- Used `Sensitive: true` on encrypted_password
- Used explicit `d.HasChange()` to detect changes
- SDK v1 framework with different state management

#### New Provider (Stainless-generated)
```go
// internal/services/cloud_load_balancer_listener/schema.go

"user_list": schema.ListNestedAttribute{
    Computed:    true,
    Optional:    true,
    PlanModifiers: []planmodifier.List{
        listplanmodifier.UseStateForUnknown(),  // ✅ Prevents spurious diffs
    },
}
```

**New Provider Approach:**
- Uses Terraform Plugin Framework (not SDK v1)
- Uses `UseStateForUnknown()` plan modifiers instead of explicit change detection
- More declarative approach - framework handles change detection
- No `Sensitive: true` needed as plan modifier handles this differently

### Why UseStateForUnknown() Solves The Problem

When a computed field uses `UseStateForUnknown()`:
1. If the field is not changing in config
2. And Terraform doesn't know the future value
3. Instead of showing "(known after apply)"
4. It uses the **current state value**

This prevents false drift detection when:
- API returns fields in different order
- API returns null vs empty for some fields
- Computed fields aren't included in update responses
- Parent resources trigger refreshes of child references

## Recommendation

### For QA Team
Test with the **latest provider build** from the current branch (`bugfix/terraform-lbpool`) which includes commit f999c1c and later fixes.

### For Development Team
The fix is already in place and working correctly. However, consider adding an automated test for this scenario:

**Suggested Test:** `test-lbpool-drift-regression`
```hcl
# Step 1: Create LB + Listener with user_list
# Step 2: terraform apply
# Step 3: Add pool resource
# Step 4: terraform apply
# Step 5: terraform plan (assert: no changes)
```

### Verification Steps
To verify the fix is working:
```bash
cd test-lbpool-drift-gcloud2-20778
terraform apply   # Creates LB + Listener
# Edit main.tf to add pool
terraform apply   # Adds pool
terraform plan    # Should show: "No changes"
```

## Conclusion

✅ **The drift issue reported in GCLOUD2-20778 has been resolved** by commit f999c1c applied on 2025-11-18.

✅ **Reproduction test confirms no drift occurs** with the current provider code.

⚠️ **The QA report was likely based on an older build** before the fix was applied.

📝 **Recommendation:** Update Jira ticket to confirm fix is working and close or mark as resolved.
