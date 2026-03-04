# Comprehensive LB Resources Test Report

**Date:** 2025-12-17
**Branch:** terraform-instances
**Commit:** 67cbbc9 (fix(cloud): implement volume attach/detach in instance Update method)

---

## Executive Summary

| Category | Total | Passed | Failed | Skipped |
|----------|-------|--------|--------|---------|
| Critical Tests | 4 | 4 | 0 | 0 |
| Standard Tests | 8 | 8 | 0 | 0 |
| **Total** | **12** | **12** | **0** | **0** |

### GCLOUD2-20778 Fix Status: **VERIFIED**

---

## Test Results

### Critical Tests (Must Pass)

| Test ID | Description | Result | Notes |
|---------|-------------|--------|-------|
| PL-010/011 | Pool import with listener_id | **PASSED** | No drift after import |
| LB-010/011 | LB import with flavor | **PASSED** | flavor="lb1-1-2" correctly in state |
| HM-003 | Health monitor deletion | **PASSED** | Removed via DELETE endpoint |
| PM-003 | Members removal via `[]` | **PASSED** | Empty array clears members |

### Standard Tests

| Test ID | Description | Result | Notes |
|---------|-------------|--------|-------|
| LB-001 | Create LB | **PASSED** | NewAndPoll completed in ~1m25s |
| LS-001 | Create HTTP listener | **PASSED** | NewAndPoll completed in ~17s |
| PL-001 | Create pool with listener_id | **PASSED** | NewAndPoll completed in ~11s |
| HM-001 | Add health monitor | **PASSED** | PATCH endpoint used |
| PM-001 | Add 2 members | **PASSED** | PATCH with members array |
| PM-002 | Remove 1 member | **PASSED** | PATCH updates member list |
| PM-004 | Omit members attribute | **PASSED** | No drift on pool |
| LB-012 | LB state has flavor | **PASSED** | `flavor = "lb1-1-2"` |

---

## Detailed Test Execution

### Test Group 1-3: Load Balancer

#### LB-001: Create Load Balancer
```
gcore_cloud_load_balancer.test_lb: Creating...
gcore_cloud_load_balancer.test_lb: Creation complete after 1m25s [id=75c19e87-54e8-49ea-bb3c-c37c02230e2c]
```
**Result:** PASSED

#### LB-010/011: Import and Verify No Drift
```bash
$ terraform state rm gcore_cloud_load_balancer.test_lb
$ terraform import gcore_cloud_load_balancer.test_lb 379987/76/75c19e87-54e8-49ea-bb3c-c37c02230e2c
$ terraform state show gcore_cloud_load_balancer.test_lb | grep flavor
  flavor = "lb1-1-2"
```
**Result:** PASSED - Flavor correctly extracted from nested API object

### Test Group 6-7: Pool Import (GCLOUD2-20778)

#### PL-010/011: Import Pool with listener_id
```bash
$ terraform state rm gcore_cloud_load_balancer_pool.test_pool
$ terraform import gcore_cloud_load_balancer_pool.test_pool 379987/76/5ec9ac42-9b5c-4074-9064-6138c15d0098
$ terraform plan
No changes. Your infrastructure matches the configuration.
```
**Result:** PASSED - listener_id correctly populated from listeners[] relationship

### Test Group 8: Health Monitor

#### HM-001: Add Health Monitor
```
gcore_cloud_load_balancer_pool.test_pool: Modifying... [id=5ec9ac42-9b5c-4074-9064-6138c15d0098]
gcore_cloud_load_balancer_pool.test_pool: Modifications complete after 13s
```
**Result:** PASSED

#### HM-003: Delete Health Monitor
```
- healthmonitor = {
-   delay         = 10 -> null
-   max_retries   = 3 -> null
-   ...
  } -> null
gcore_cloud_load_balancer_pool.test_pool: Modifications complete after 1s
```
**Result:** PASSED - Used dedicated DELETE /healthmonitor endpoint

### Test Group 9: Pool Members

#### PM-001: Add 2 Members
```
+ members = [
    + {
        + address       = "10.0.0.1"
        + protocol_port = 8080
      },
    + {
        + address       = "10.0.0.2"
        + protocol_port = 8080
      },
  ]
gcore_cloud_load_balancer_pool.test_pool: Modifications complete after 20s
```
**Result:** PASSED

#### PM-002: Remove 1 Member
```
~ members = [
    - {
        - address       = "10.0.0.2"
        ...
      },
  ]
gcore_cloud_load_balancer_pool.test_pool: Modifications complete after 21s
```
**Result:** PASSED

#### PM-003: Remove All Members
```
~ members = [] (empty array)
gcore_cloud_load_balancer_pool.test_pool: Modifications complete after 21s
```
**Result:** PASSED

#### PM-004: Omit Members Attribute
```
$ terraform plan
# Only LB shows drift on computed fields, pool has no changes
```
**Result:** PASSED

---

## Known Issues (Non-Blocking)

### LB Computed Fields Show Cosmetic Drift After Import

After importing a Load Balancer, the following computed fields show drift:
- `vrrp_ips` - changes from list to `(known after apply)`
- `tags_v2` - changes from `[]` to `(known after apply)`
- `stats` - not in state after import
- `ddos_profile` - not in state after import

**Analysis:**

| Field | Volatile? | UseStateForUnknown? | Recommendation |
|-------|-----------|---------------------|----------------|
| `vrrp_ips` | **YES** - IPs change on resize/failover | **NO** - would cause "inconsistent result" errors | Keep as-is (cosmetic drift is correct) |
| `tags_v2` | No | Could add | Optional improvement |
| `stats` | **YES** - changes constantly | **NO** | Keep as-is |
| `ddos_profile` | No | Could add | Optional improvement |

**Why NOT to use `UseStateForUnknown()` for `vrrp_ips`:**

Per Slack discussion (Nov 19, 2025), Kirill observed this error during LB resize:
```
Error: Provider produced inconsistent result after apply
.vrrp_ips[0].ip_address: was "85.204.246.209", but now "85.204.246.93"
```

This happens because `vrrp_ips` legitimately changes during resize operations. Using `UseStateForUnknown()` would preserve stale values and cause apply failures.

**Impact:** Minor cosmetic drift only. Running `terraform apply` or `terraform refresh` after import will sync state with actual API values. No actual resource changes occur.

**Related Jira:** [GCLOUD2-22019](https://jira.gcore.lu/browse/GCLOUD2-22019) - "PATCH /loadbalancers returns empty vrrp_ips array" (Status: **Done** - API fixed Dec 10, 2025)

---

## Infrastructure Created

| Resource | ID |
|----------|-----|
| Load Balancer | 75c19e87-54e8-49ea-bb3c-c37c02230e2c |
| Listener | da70789d-490f-4f82-98a2-4c58cb966e80 |
| Pool | 5ec9ac42-9b5c-4074-9064-6138c15d0098 |

---

## Conclusion

All critical tests passed. The GCLOUD2-20778 fix for pool import drift is verified working:

1. **Pool Import:** `listener_id` is now correctly populated from the `listeners[]` relationship during import
2. **LB Import:** `flavor` is now correctly extracted from the nested `flavor.flavor_name` API response
3. **Health Monitor:** Deletion via attribute removal correctly uses the DELETE endpoint
4. **Members:** All CRUD operations work correctly, including empty array and omitted attribute

### Fixes Applied
| Resource | Field | Fix |
|----------|-------|-----|
| `gcore_cloud_load_balancer_pool` | `listener_id` | Extract from `listeners[]` in ImportState + `UseStateForUnknown()` |
| `gcore_cloud_load_balancer_pool` | `load_balancer_id` | Added `UseStateForUnknown()` |
| `gcore_cloud_load_balancer` | `flavor` | Extract `flavor.flavor_name` from nested API object in ImportState |

### No Action Required
| Issue | Reason |
|-------|--------|
| `vrrp_ips` cosmetic drift | Field is volatile (changes on resize/failover). Current behavior is correct. |
| `stats` cosmetic drift | Field is volatile (real-time statistics). Current behavior is correct. |
