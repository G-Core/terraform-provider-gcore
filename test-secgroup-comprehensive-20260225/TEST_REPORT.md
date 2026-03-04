# Test Report: gcore_cloud_security_group + gcore_cloud_security_group_rule

**Date**: 2026-02-25
**Branch**: `feature/GCLOUD2-20783-secgroup-rule-resource`
**Commit**: `f43f722`
**Ticket**: GCLOUD2-20783

## Summary

| Test # | Description | Result | Notes |
|--------|-------------|--------|-------|
| 1 | Create SG (no default rules) | PASS | `rules:[]` via WithJSONSet |
| 2 | Create egress TCP 443 rule | PASS | Async NewAndPoll |
| 3 | Create ingress ICMP rule (no ports) | PASS | No port fields sent |
| 4 | Create ingress UDP 5000-6000 rule | PASS | Port range works |
| 5 | Update SG name (in-place) | PASS | Same ID preserved |
| 6 | Update SG description | PASS | Same ID preserved |
| 7 | Import SG | PASS* | *Known drift on project_id/region_id (GCLOUD2-23752) |
| 8 | Import rule (4-part path) | PASS* | *Same known drift issue |
| 9 | Delete single rule | PASS | Other rules unaffected |
| 10 | Destroy all | PASS | Clean teardown, API 404 |

**Overall: 10/10 PASS** (2 with known import drift issue affecting all resources)

## Resource IDs

### First Create Cycle
- SG: `e91abc80-c14b-4e13-92ac-8ab44fa2fb1d`
- Egress TCP rule: `88cdd6f2-0c38-4aba-b5db-48ae6ed147f5`
- Ingress ICMP rule: `773f0f19-030f-4953-8b26-53fa9472bd3c`
- Ingress UDP rule: `8cff8f10-9aea-4fe2-81b6-b41cf182ca4d`

### After Import Re-create
- SG: `9dfdc5d1-a603-4112-a271-58b0d78da5ae`
- Egress TCP rule: `d49db61b-7eca-40ab-84ee-3c3a34a95641`
- Ingress ICMP rule: `725f4046-34fd-4e37-b885-d445c9861bad`
- Ingress UDP rule: `f2c9d022-fffa-433a-b4c8-c5e2a5220b29`

## Drift Test Results

| Resource | After Create | After Update | After Delete | After Import |
|----------|-------------|-------------|-------------|-------------|
| SG | exit 0 | exit 0 | exit 0 | exit 2* |
| TCP rule | exit 0 | N/A | exit 0 | exit 2* |
| ICMP rule | exit 0 | N/A | N/A (deleted) | N/A |
| UDP rule | exit 0 | N/A | exit 0 | exit 2* |

*Import drift is due to provider-wide `project_id`/`region_id` issue (GCLOUD2-23752)

## Test Details

### Test 1-4: Create SG + 3 Rules
- SG created with `option.WithJSONSet("rules", []any{})` to suppress default egress rules
- All 3 rules created in parallel (~4-7s each via async polling)
- Immediate drift check: exit code 0 (no changes)

### Test 5: SG Name Update
- Changed name from `tf-test-secgroup-comprehensive` to `tf-test-secgroup-renamed`
- Plan showed `~ update in-place` (not replacement)
- Same SG ID before and after: `e91abc80-c14b-4e13-92ac-8ab44fa2fb1d`
- Post-update drift: exit code 0

### Test 6: SG Description Update
- Changed description to "Updated description for testing"
- In-place update, same ID preserved
- Post-update drift: exit code 0

### Test 7: SG Import
- Import path: `379987/76/<sg_id>`
- Import succeeded, state populated correctly
- **Known issue**: Import sets `project_id=379987, region_id=76` in state, but config has them as null, causing `forces replacement`

### Test 8: Rule Import
- Import path: `379987/76/<group_id>/<rule_id>` (4-part)
- Import succeeded, all rule fields populated correctly
- Same known import drift issue cascades to rules (group_id changes when SG is replaced)

### Test 9: Single Rule Deletion
- Removed ICMP rule from config
- Plan: `0 to add, 0 to change, 1 to destroy`
- Only ICMP rule destroyed, TCP and UDP rules kept identical IDs
- Post-delete drift: exit code 0

### Test 10: Full Destroy
- Clean destruction: rules deleted first (parallel), then SG
- API returns 404 for all deleted resources
- No orphaned resources (verified via API)

## Known Issues

### GCLOUD2-23752: Import Drift on project_id/region_id
- **Affects**: ALL importable resources in the provider (19 of 20)
- **Cause**: Import path populates `project_id`/`region_id` with explicit values, but config uses provider defaults (null in state)
- **Impact**: After import, plan shows `forces replacement` on these fields
- **Fix**: Apply `RequiresReplaceIfConfiguredPreservingState()` plan modifier (only `cloud_instance` has this fix today)

## Code Quality

### Schema Design
- SG: 8 attributes (3 required/optional, 5 computed)
- Rule: 12 attributes (2 required, 8 optional, 1 computed, 1 path-required)
- All rule fields use `RequiresReplace` (immutable after creation)
- Proper validators: direction, ethertype, protocol, port ranges

### API Methods
- SG: `NewAndPoll`, `UpdateAndPoll`, `Get`, `Delete` (sync)
- Rule: `NewAndPoll`, `DeleteAndPoll`
- Rule Read: Fetches parent SG, iterates rules list (no individual GET)
- Delete: Idempotent (404 silently ignored)

## Cleanup Verification
```
curl "https://api.gcore.com/cloud/v1/securitygroups/379987/76" | jq '[.results[] | select(.name | test("tf-test"))]'
# Result: [] (no test resources remain)
```
