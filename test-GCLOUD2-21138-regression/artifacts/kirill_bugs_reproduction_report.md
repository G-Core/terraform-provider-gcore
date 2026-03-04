# Kirill's Bug Reproduction Report (GCLOUD2-21138)

**Date**: 2025-12-30
**Tester**: Claude Code
**Commit tested**: c5947a6 (terraform-instances branch) + additional fixes

## Summary

| Bug | Description | Status | Notes |
|-----|-------------|--------|-------|
| Bug 1 | New FIP on private interface | **FIXED** | Required 3 code changes |
| Bug 2 | Import with two volumes boot_index drift | **NOT REPRODUCED** | Fix in c5947a6 works correctly |

---

## Bug 1: New Floating IP on Private Interface

### Description
When adding `floating_ip { source = "new" }` to a private interface, Terraform produces "Provider produced inconsistent result after apply" error.

### Steps to Reproduce
1. Create instance with private interface (no FIP)
2. Update config to add `floating_ip { source = "new" }` to the interface
3. Run `terraform apply`

### Expected Result
FIP should be created, assigned to port, and state should show `source = "new"`

### Actual Result (Before Fix)
```
Error: Provider produced inconsistent result after apply

.interfaces[0].floating_ip.existing_floating_id: was null, but now
cty.StringVal("26369297-93b4-4095-8f04-eae0517d4f5b")

.interfaces[0].floating_ip.source: was cty.StringVal("new"), but now
cty.StringVal("existing")
```

### Root Cause Analysis
1. User specifies `floating_ip { source = "new" }` in config
2. Provider creates new FIP via POST /floatingips (works correctly)
3. Provider assigns FIP to port (works correctly)
4. When reading state after apply, provider returns:
   - `source = "existing"` (instead of preserving "new")
   - `existing_floating_id = <the new FIP ID>`
5. Terraform detects mismatch between planned (`source = "new"`) and actual (`source = "existing"`)

### Fixes Applied

**Fix 1: Update function (lines 812-823)**
Changed the code after creating new FIP to NOT change source to "existing":
```go
// FIP created and assigned to port successfully.
// Keep source="new" to match user config - don't change to "existing"
// Don't populate existing_floating_id - it's only for user-specified existing FIPs
tflog.Info(ctx, "Created and assigned new floating IP", map[string]interface{}{
    "floating_ip_id": newFIP.ID,
    "port_id":        portID,
    "interface":      i,
})
floatingIPChanged = true
continue
```

**Fix 2: Read function (lines 1332-1357)**
Preserve original source value from state when reading:
```go
// Preserve source from current state - if user specified "new", keep it
currentSource := "existing"
if (*data.Interfaces)[i].FloatingIP != nil && !(*data.Interfaces)[i].FloatingIP.Source.IsNull() {
    currentSource = (*data.Interfaces)[i].FloatingIP.Source.ValueString()
}
if currentSource == "new" {
    (*data.Interfaces)[i].FloatingIP = &CloudInstanceInterfacesFloatingIPModel{
        Source:             types.StringValue("new"),
        ExistingFloatingID: types.StringNull(),
    }
}
```

**Fix 3: FIP removal when source="new" (lines 826-897)**
Added logic to look up and unassign FIP when user removes `floating_ip { source = "new" }`:
```go
// Handle removal of FIP that was created with source="new"
if stateSource == "new" && planIface.FloatingIP == nil && portID != "" {
    // List floating IPs to find the one attached to this port
    fipPage, err := r.client.Cloud.FloatingIPs.List(ctx, listParams, ...)
    // Find FIP attached to our port
    for fipPage != nil && len(fipPage.Results) > 0 {
        for _, fip := range fipPage.Results {
            if fip.PortID == portID {
                fipToUnassign = fip.ID
                break
            }
        }
    }
    // Unassign the FIP
    r.client.Cloud.FloatingIPs.Unassign(ctx, fipToUnassign, ...)
}
```

### Test Results After Fixes

| Test | Result |
|------|--------|
| Create instance with no FIP | PASS |
| Add FIP with source="new" | PASS (no inconsistent state error) |
| Plan after adding FIP | PASS (exit code 0, no drift) |
| Remove FIP | PASS |
| Plan after removing FIP | PASS (exit code 0, no drift) |
| Re-add FIP | PASS |
| Plan after re-adding FIP | PASS (exit code 0, no drift) |

### Artifacts
- `artifacts/bug1_create_no_fip.log` - Instance creation with private interface
- `artifacts/bug1_add_fip.log` - Adding FIP with source="new"
- `artifacts/bug1_add_fip_drift.log` - Plan after adding FIP (no changes)
- `artifacts/bug1_remove_fip.log` - Removing FIP
- `artifacts/bug1_remove_fip_drift.log` - Plan after removing FIP (no changes)
- `artifacts/bug1_readd_fip.log` - Re-adding FIP

---

## Bug 2: Import with Two Volumes Shows boot_index Drift

### Description
After importing a VM with 2 volumes, terraform plan shows drift: `boot_index = -1 -> 0`

### Steps to Reproduce
1. Create instance with two volumes:
   - Boot volume: `boot_index = 0`
   - Data volume: `boot_index = -1`
2. `terraform state rm` the instance
3. `terraform import` the instance
4. `terraform plan -detailed-exitcode`

### Expected Result
Plan should show no changes (exit code 0)

### Actual Result (BUG NOT REPRODUCED)
```
No changes. Your infrastructure matches the configuration.
Exit code: 0
```

### Analysis
The fix in commit c5947a6 correctly handles boot_index during import:

1. Provider fetches instance from API - volumes array does NOT include boot_index
2. For each volume, provider calls Volume API to get `bootable` field
3. Sets `boot_index = 0` for bootable volumes, `-1` for non-bootable
4. State is correctly populated, plan shows no drift

**Fix working correctly** - resource.go lines 1457-1481:
```go
// Fetch volume details from Volume API to check bootable field
volResp, err := r.client.Cloud.Volumes.Get(...)
if volResp.Bootable {
    volModel.BootIndex = types.Int64Value(0)
} else {
    volModel.BootIndex = types.Int64Value(-1)
}
```

### Artifacts
- `artifacts/bug2_step1_create.log` - Instance creation with two volumes
- `artifacts/bug2_step3_import.log` - Import command output
- `artifacts/bug2_step4_plan.log` - Plan showing no changes

---

## Test Environment

- Instance ID: f014db96-ac8a-4297-8706-48117cd98aa4 (Bug 1 test)
- Instance ID: 38346859-1d7f-461a-a6c7-867e48b81b1e (Bug 2 test)
- Region: Luxembourg-2 (76)
- Project: 379987

## Files Changed in resource.go

1. **Lines 812-823**: Update function - preserve source="new" after FIP creation
2. **Lines 826-897**: Update function - handle FIP removal when source="new"
3. **Lines 1332-1357**: Read function - preserve source value from state

## Recommendations

1. **Bug 1 (FIP on private interface)**: FIXED - all test cases pass
2. **Bug 2 (boot_index drift)**: Already fixed in c5947a6

## Next Steps

1. Run full regression test suite
2. Create PR with the Bug 1 fixes
3. Close GCLOUD2-21138 ticket
