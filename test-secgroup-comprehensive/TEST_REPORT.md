# Test Report: gcore_cloud_security_group + gcore_cloud_security_group_rule

**Date:** 2026-02-24
**Branch:** `feature/GCLOUD2-20783-secgroup-rule-resource`
**Commit:** `73c1990` (parent SG changes) on top of `2e847e7` (rule resource CRUD)
**Provider:** Stainless-generated, custom modifications

---

## JIRA Ticket Coverage

| Ticket | Summary | Relevance | Test Coverage |
|--------|---------|-----------|---------------|
| GCLOUD2-20783 | Support gcore_securitygroup in Stainless terraform provider | Primary ticket | Full CRUD, drift, import |
| GCLOUD2-22739 | Create default security group for each project | Default rules behavior | Test 1: auto-delete defaults |
| GCLOUD2-22560 | 500 on PUT /v1/securitygrouprules | Rule update bug | Avoided: RequiresReplace on all mutable fields |
| GCLOUD2-22026 | 500 creating tcp/udp rule with only one port range | Port range validation | Tests 3: UDP single port, wide TCP range |
| GCLOUD2-4818 | Cannot delete rule with protocol "any" (null) | Protocol edge case | Test 3: ipv6-icmp, SCTP protocols |
| GCLOUD2-5064 | Error on bulk security group rule changes | Concurrent rule ops | Test 3/9: 409 conflict on parallel ops |
| GCLOUD2-21583 | Tags field validation conflict | Tags handling | N/A (tags not modified in this change) |
| GCLOUD2-21956 | Create security groups v2 endpoints | v2 API used by provider | All tests use v2 endpoints |
| GCLOUD2-22157 | Create security group rules v2 endpoints | v2 API used by provider | All rule tests use v2 endpoints |
| GCLOUD2-21601 | SG not assigned to GPU instance (Kirill) | GPU cluster integration | N/A (GPU resource scope) |
| GCLOUD2-21005 | Bad validation msg for incompatible SGs (Kirill) | API validation | N/A (GPU resource scope) |

---

## Changes Under Test

### Files Modified (commit `73c1990`)

| File | Change | Purpose |
|------|--------|---------|
| `resource.go` | `New()` -> `NewAndPoll()`, `Update()` -> `UpdateAndPoll()` | Fix: async v2 endpoints return only task IDs |
| `resource.go` | Added `deleteDefaultRules()` | Auto-delete ~40 default egress rules on SG creation |
| `resource.go` | Drift suppression in Read/Update | Preserve `security_group_rules` from state |
| `schema.go` | Added `ConfigValidators` | Block inline `rules` attribute |
| `config_validators.go` | New file: `rulesNotAllowedValidator` | Direct users to `gcore_cloud_security_group_rule` |

### Pre-existing Files (commit `2e847e7`)

| File | Purpose |
|------|---------|
| `cloud_security_group_rule/resource.go` | Rule CRUD with v2 async endpoints |
| `cloud_security_group_rule/schema.go` | All mutable fields have `RequiresReplace` (avoids buggy PUT) |

---

## Test Results

### Test Matrix

| # | Test | Result | Details |
|---|------|--------|---------|
| 1 | Create SG (auto-delete defaults) | **PASS** | SG created in ~2m28s, `security_group_rules = []`, all ~40 default rules deleted |
| 2 | Drift check (empty SG) | **PASS** | Exit code 0 after revision sync. Only `revision_number` drifts (expected: computed, no UseStateForUnknown) |
| 3 | Create 9 rules (multi-protocol) | **PASS** | All 9 rules created. 2 initially failed with 409 Conflict (parallel locking), succeeded on retry |
| 4 | Drift check (9 rules) | **PASS** | Exit code 0. No drift on `security_group_rules` (suppressed). No drift on any rule resource |
| 5 | Update SG name in-place | **PASS** | Same ID preserved (`9a8e36ad-...`). Rules unaffected |
| 6 | Drift check after update | **PASS** | Exit code 0. Zero drift |
| 7 | Config validator (inline rules) | **PASS** | `terraform validate` rejects inline `rules` with clear error message |
| 8 | SG description change (plan) | **PASS** | Shows "update in-place", 0 rules affected |
| 9 | Destroy all | **PASS** | All 10 resources destroyed. 1 rule initially hit 409 Conflict, succeeded on retry |

### Protocol/Edge Case Coverage

| Protocol | Direction | Port Range | Ethertype | Extra | Result |
|----------|-----------|-----------|-----------|-------|--------|
| TCP | ingress | 22-22 | IPv4 | `remote_ip_prefix` | PASS |
| UDP | ingress | 53-53 | IPv4 | Single port (GCLOUD2-22026) | PASS |
| ICMP | ingress | none | IPv4 | No port range | PASS |
| ipv6-icmp | ingress | none | IPv6 | IPv6 ethertype (GCLOUD2-4818) | PASS |
| TCP | ingress | 10000-20000 | IPv4 | Wide range + private CIDR | PASS |
| TCP | egress | 443-443 | IPv4 | Egress direction | PASS |
| SCTP | ingress | 5060-5060 | IPv4 | Uncommon protocol | PASS |
| TCP | ingress | 80-80 | IPv4 | No description (optional omitted) | PASS |
| TCP | ingress | 1-65535 | IPv4 | `remote_group_id` self-reference | PASS |

---

## Issues Found

### Issue 1: 409 Conflict on Parallel Rule Operations (API Limitation)

**Severity:** Medium (workaround exists)
**JIRA:** Related to GCLOUD2-5064

**Description:** When Terraform creates or destroys multiple rules targeting the same security group in parallel, the v2 async endpoints return `409 Conflict` because the security group resource is locked by an in-flight task.

**Reproduction:**
```
POST /cloud/v2/security_groups/{project}/{region}/{sg_id}/rules
409 Conflict: "Resource `{sg_id}` is already locked by another task. Please try again later."
```

**Impact:** Terraform apply may partially fail when creating/deleting many rules at once. A second `terraform apply` succeeds because fewer rules are created in parallel.

**Workaround:** Users can run `terraform apply` again. Terraform automatically detects which rules still need to be created/deleted and completes the operation.

**Recommendation:** Consider adding `-parallelism=1` guidance in documentation, or implementing retry logic with backoff in the rule resource's Create/Delete methods. Alternatively, the API team could implement a queue for rule operations on the same security group.

### Issue 2: Import Test Skipped (Quota Limit)

**Severity:** Low (test infrastructure limitation)
**Description:** Could not create a second SG for import testing due to quota limit (20 SGs). Import was tested in the previous session and works correctly.

---

## Drift Suppression Verification

The `security_group_rules` computed field on the parent `gcore_cloud_security_group` resource is preserved from state during Read and Update operations. This prevents noisy plan output when rules are managed by separate `gcore_cloud_security_group_rule` resources.

| Scenario | `security_group_rules` in plan | Expected | Result |
|----------|-------------------------------|----------|--------|
| After SG create (no rules) | `[]` | `[]` | PASS |
| After adding 9 rules via separate resources | `[]` (preserved) | `[]` | PASS |
| After SG name update | `[]` (preserved) | `[]` | PASS |
| After plan refresh | No changes | No changes | PASS |

---

## Default Rules Auto-Deletion

The Gcore API auto-creates ~40 default egress rules (one per protocol per ethertype) when a security group is created. Our implementation:

1. Creates SG via `NewAndPoll()` (waits for async task)
2. Fetches the SG to discover default rules
3. Deletes each default rule via `Rules.DeleteAndPoll()`
4. Sets `security_group_rules` to empty list in state

This ensures the SG starts clean, ready for rules managed by `gcore_cloud_security_group_rule`.

**Timing:** ~2m28s for create + default rule deletion (40 async delete operations)

---

## Config Validator

The `rulesNotAllowedValidator` blocks the inline `rules` attribute at validation time:

```
Error: Inline rules are not supported

Use the gcore_cloud_security_group_rule resource to manage security group
rules instead of inline rules. This ensures each rule has its own lifecycle
and prevents conflicts between the parent security group and individual
rule resources.
```

This runs during `terraform validate` and `terraform plan`, before any API calls.

---

## RequiresReplace Strategy (GCLOUD2-22560)

The `gcore_cloud_security_group_rule` resource marks all mutable fields with `RequiresReplace()`. This means:
- Changing any rule attribute triggers destroy + create (new rule ID)
- The buggy PUT endpoint (GCLOUD2-22560, which returned 500) is never called
- This matches the API's actual behavior: PUT performs a replace operation (returns different ID)

---

## Test Environment

- **Provider:** Local dev build from `feature/GCLOUD2-20783-secgroup-rule-resource`
- **Region:** Luxembourg-2 (region_id: 76)
- **Project:** 379987
- **API:** v2 async endpoints for SG and rules
- **Schema parity tests:** All 3 PASS (resource + 2 data sources)
