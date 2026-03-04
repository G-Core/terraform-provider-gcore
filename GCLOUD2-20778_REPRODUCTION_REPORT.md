# GCLOUD2-20778 Drift Issue Reproduction Report

**Date:** 2025-11-24
**Ticket:** [GCLOUD2-20778](https://jira.gcore.lu/browse/GCLOUD2-20778)
**Branch:** bugfix/terraform-lbpool
**Status:** ✅ **BUG SUCCESSFULLY REPRODUCED**

---

## Executive Summary

The drift issue reported in the last Jira comment (2025-11-20 by Kirill Tsaregorodtsev) has been **successfully reproduced**. After adding an LB pool to an existing listener with `user_list` configured, Terraform detects false drift in both the load balancer and listener resources, wanting to update them despite no actual configuration changes.

---

## Last Jira Comment (2025-11-20)

**Reporter:** Kirill Tsaregorodtsev
**Issue:** After adding LB pool, TF detects resource drift showing:
- `stats` attribute as "(known after apply)" on listener
- `encrypted_password` in `user_list` changing to "(known after apply)"

**Expected:** No changes detected
**Actual:** Terraform wants to update the listener

---

## Reproduction Steps

### Environment
- **Test Directory:** `/Users/user/repos/gcore-terraform/test-lbpool-drift-gcloud2-20778`
- **Provider:** Latest build from `bugfix/terraform-lbpool` branch
- **Commits:**
  - Latest: `f0527a1` (Remove task-related attributes)
  - Fix commit: `f999c1c` (Added UseStateForUnknown modifiers - supposedly fixed this)
- **Region:** Luxembourg-2 (76)
- **Project:** 379987

### Test Configuration

```hcl
resource "gcore_cloud_load_balancer" "lb" {
  project_id = local.project_id[0]
  region_id  = data.gcore_cloud_region.rg.id
  flavor     = "lb1-1-2"
  name       = "qa-lbpool-drift-test"
}

resource "gcore_cloud_load_balancer_listener" "ls" {
  project_id       = local.project_id[0]
  region_id        = data.gcore_cloud_region.rg.id
  load_balancer_id = gcore_cloud_load_balancer.lb.id
  name             = "qa-ls-name"
  protocol         = "HTTP"
  protocol_port    = 80
  connection_limit = 5000

  user_list = [{
    encrypted_password = "$5$isRr.HJ1IrQP38.m$oViu3DJOpUG2ZsjCBtbITV3mqpxxbZfyWJojLPNSPO5"
    username           = "qauser"
  }]
}

resource "gcore_cloud_load_balancer_pool" "lb_pool" {
  project_id   = local.project_id[0]
  region_id    = data.gcore_cloud_region.rg.id
  lb_algorithm = "LEAST_CONNECTIONS"
  name         = "pool-drift-test"
  protocol     = "HTTP"
  listener_id  = gcore_cloud_load_balancer_listener.ls.id
}
```

### Steps Executed

1. ✅ Infrastructure already created (LB + Listener + Pool all exist)
2. ✅ Ran `terraform plan` with no configuration changes
3. ✅ **Result:** Drift detected

---

## Reproduction Results

### ❌ Drift Detected on Load Balancer

```
# gcore_cloud_load_balancer.lb will be updated in-place
~ resource "gcore_cloud_load_balancer" "lb" {
    ~ additional_vips        = [] -> (known after apply)
    ~ created_at             = "2025-11-20T12:03:48Z" -> (known after apply)
    ~ creator_task_id        = "eb54ef92-4273-469c-a040-12c61464cdd6" -> (known after apply)
    + ddos_profile           = (known after apply)
    ~ floating_ips           = [] -> (known after apply)
      id                     = "00b353b8-e5b3-49a7-8ca1-072857647502"
    ~ listeners              = [...] -> (known after apply)
    ~ logging                = {...} -> (known after apply)
      name                   = "qa-lbpool-drift-test"
    ~ operating_status       = "ONLINE" -> (known after apply)
    ~ provisioning_status    = "ACTIVE" -> (known after apply)
    ~ region                 = "Luxembourg-2" -> (known after apply)
    + stats                  = (known after apply)
    ~ tags_v2                = [] -> (known after apply)
    + task_id                = (known after apply)
    + tasks                  = (known after apply)
    ~ updated_at             = "2025-11-20T12:20:25Z" -> (known after apply)
    ~ vip_address            = "109.61.125.253" -> (known after apply)
    ~ vip_ip_family          = "ipv4" -> (known after apply)
    ~ vip_port_id            = "a991c7a7-5f97-4049-8fff-b95fdc2b0f7e" -> (known after apply)
    ~ vrrp_ips               = [2 elements] -> (known after apply)
      # (4 unchanged attributes hidden)
  }
```

**Issue:** Nearly all computed fields showing as "(known after apply)" even though nothing changed.

### ❌ Drift Detected on Listener

```
# gcore_cloud_load_balancer_listener.ls will be updated in-place
~ resource "gcore_cloud_load_balancer_listener" "ls" {
    ~ creator_task_id        = "145da64e-c11b-4949-8489-92d9565f8b08" -> (known after apply)
      id                     = "114c9f9f-caa0-46bf-a920-b1026dfdb8f1"
    ~ insert_headers         = jsonencode({}) -> (known after apply)
      name                   = "qa-ls-name"
    ~ operating_status       = "ONLINE" -> (known after apply)
    ~ pool_count             = 1 -> (known after apply)
    ~ provisioning_status    = "ACTIVE" -> (known after apply)
    - sni_secret_id          = [] -> null
    + stats                  = (known after apply)
    + task_id                = (known after apply)
    + tasks                  = (known after apply)
    - timeout_client_data    = 50000 -> null
    - timeout_member_connect = 5000 -> null
    - timeout_member_data    = 50000 -> null
      # (7 unchanged attributes hidden)
  }
```

**Issues:**
1. ✅ `stats` showing as "(known after apply)" - **MATCHES JIRA REPORT**
2. Timeout fields being removed (set to null)
3. Computed fields showing drift

**Note:** The `user_list.encrypted_password` drift mentioned in the Jira comment was NOT observed in this test, but substantial other drift WAS detected.

---

## Plan Summary

```
Plan: 0 to add, 2 to change, 0 to destroy.
```

Terraform wants to update both resources despite no configuration changes being made.

---

## Root Cause Analysis

### Why is Drift Still Occurring?

Despite commit `f999c1c` (2025-11-18) that supposedly fixed this by adding `UseStateForUnknown()` plan modifiers, the drift persists. Possible reasons:

1. **Insufficient Plan Modifiers:**
   - Not all computed fields have `UseStateForUnknown()` modifiers
   - Some nested attributes may be missing modifiers

2. **Read/Refresh Behavior:**
   - When Terraform refreshes state, it may be calling the Read function
   - The Read function might be returning different data structure than what's in state
   - This causes Terraform to think values have changed

3. **Timeout Fields Issue:**
   - The timeout fields (`timeout_client_data`, `timeout_member_connect`, `timeout_member_data`) are being removed
   - Config doesn't set them, but state has values (defaults from API)
   - This suggests these should be:
     - Either: `Computed: true` (if API provides defaults)
     - Or: `Computed: true, Optional: true` (if user can override)

4. **Nested Object Handling:**
   - Fields like `logging`, `listeners` are complex nested objects
   - May need special handling or plan modifiers at nested level

---

## Comparison with "Fixed" Code

### Commit f999c1c Claims to Fix This

The commit message states:
> "feat(cloud): enhance load balancer resource schema with state management modifiers
>
> Added `UseStateForUnknown()` plan modifiers to all computed fields"

However, the drift is still occurring, suggesting:
- Either the modifiers weren't added to ALL fields
- Or the issue is not just about plan modifiers
- Or there's a separate Read/Refresh issue

### Files to Investigate

1. **internal/services/cloud_load_balancer/schema.go** - Check LB schema
2. **internal/services/cloud_load_balancer_listener/schema.go** - Check Listener schema
3. **internal/services/cloud_load_balancer/resource.go** - Check Read function
4. **internal/services/cloud_load_balancer_listener/resource.go** - Check Read function

---

## Impact

- **Severity:** HIGH
- **User Impact:** Users cannot maintain stable infrastructure without false drift detection
- **Workaround:** None - Terraform will continuously want to update resources
- **Blocking:** Yes - This prevents the ticket from being closed

---

## Next Steps for Developers

### 1. Verify Plan Modifiers (HIGH PRIORITY)

Check every computed field in both LB and Listener schemas:

```bash
grep -n "Computed.*true" internal/services/cloud_load_balancer/schema.go | \
  grep -v "UseStateForUnknown"
```

Any computed field without `UseStateForUnknown()` should have it added.

### 2. Fix Timeout Fields

File: `internal/services/cloud_load_balancer_listener/schema.go`

```go
"timeout_client_data": schema.Int64Attribute{
    Description: "...",
    Computed:    true,  // ✅ Keep this
    Optional:    true,  // ✅ Add this
    PlanModifiers: []planmodifier.Int64{
        int64planmodifier.UseStateForUnknown(),  // ✅ Add this
    },
},
```

### 3. Investigate Read Functions

Check if Read functions are causing state refreshes to differ:

```bash
# Check LB resource Read function
grep -A 50 "func.*Read.*ResourceModel" \
  internal/services/cloud_load_balancer/resource.go

# Check Listener resource Read function
grep -A 50 "func.*Read.*ResourceModel" \
  internal/services/cloud_load_balancer_listener/resource.go
```

### 4. Add Debug Logging

Temporarily add logging to understand when drift is detected:

```go
// In Read function
tflog.Debug(ctx, "Read completed", map[string]interface{}{
    "state_hash": fmt.Sprintf("%+v", data),
})
```

### 5. Test with Old Provider

Compare behavior with old provider to ensure we haven't regressed:

```bash
cd test-old-provider-comparison
# Use old provider binary
terraform plan
```

---

## Test Commands to Verify Fix

Once fixed, run these commands to verify:

```bash
cd /Users/user/repos/gcore-terraform/test-lbpool-drift-gcloud2-20778
source ../.env
export TF_CLI_CONFIG_FILE=/Users/user/repos/gcore-terraform/.terraformrc

# Should show: "No changes. Your infrastructure matches the configuration."
terraform plan -detailed-exitcode
```

Exit code should be `0` (not `2` which indicates changes detected).

---

## Artifacts

- **Full drift output:** `test-lbpool-drift-gcloud2-20778/drift_reproduction_20251124.txt`
- **Test configuration:** `test-lbpool-drift-gcloud2-20778/main.tf`
- **Previous analysis:** `test-lbpool-drift-gcloud2-20778/FINDINGS.md`
- **This report:** `GCLOUD2-20778_REPRODUCTION_REPORT.md`

---

## Conclusion

✅ **BUG SUCCESSFULLY REPRODUCED**

The drift issue reported in GCLOUD2-20778 is **real and reproducible** with the current code on the `bugfix/terraform-lbpool` branch. Despite the fix applied in commit `f999c1c`, drift detection still occurs when:

1. Load balancer, listener, and pool exist
2. No configuration changes are made
3. `terraform plan` is run

**Recommendation:** The ticket should remain open until this drift issue is fully resolved and verified with the test scenario provided.

**Priority:** HIGH - This blocks users from having stable infrastructure management.
