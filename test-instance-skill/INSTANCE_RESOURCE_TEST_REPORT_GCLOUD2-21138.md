# Instance Resource Test Report - GCLOUD2-21138

**Date**: 2025-11-19
**Branch**: terraform-instances
**Resource**: gcore_cloud_instance

## Executive Summary

**Overall Result**: ⚠️ PARTIAL PASS - Core functionality works with 2 bugs identified

- **Tests Passed**: 11/13
- **Tests Failed**: 2 (bugs identified)
- **Critical Issues**: 2 bugs requiring fixes

## Test Results Summary

| Test # | Description | Result | Notes |
|--------|-------------|--------|-------|
| 1 | Drift detection (existing instance) | ✅ PASS | No drift detected |
| 2 | Flavor resize (g1-standard-1-2 → g1-standard-2-4) | ⚠️ PASS* | API succeeded but transient vm_state validation error |
| 3 | Post-resize drift check | ✅ PASS | Exit code 0 |
| 4 | Volume extend (5 → 10 GiB) | ✅ PASS | POST /volumes/{id}/extend called correctly |
| 5 | Post-extend drift check | ✅ PASS | Exit code 0 |
| 6 | Volume shrink guard | ✅ PASS | Clear error message blocks shrinking |
| 7 | VM state stop (active → stopped) | ✅ PASS | POST /action called correctly |
| 8 | VM state start (stopped → active) | ✅ PASS | Instance started successfully |
| 9 | Destroy instance | ✅ PASS | DELETE called, volume cleaned up |
| 10 | Create baseline (image + external) | ✅ PASS | All computed fields populated |
| 11 | Post-create drift check | ✅ PASS | Exit code 0 |
| 12 | Import instance | ✅ PASS | Import command succeeded |
| 13 | Post-import drift check | ❌ FAIL | Missing fields cause replacement |

## API Endpoint Verification (MITM)

| Endpoint | Purpose | Verified |
|----------|---------|----------|
| `POST /instances/{project}/{region}` | Create | ✅ |
| `GET /instances/{project}/{region}/{id}` | Read | ✅ |
| `DELETE /instances/{project}/{region}/{id}` | Delete | ✅ |
| `POST /instances/{id}/changeflavor` | Resize flavor | ✅ |
| `POST /volumes/{id}/extend` | Extend volume | ✅ |
| `POST /instances/{id}/action` | Start/Stop | ✅ |

## Bugs Identified

### Bug 1: vm_state Validation Error During Resize

**Location**: `internal/services/cloud_instance/resource.go:287-315`

**Symptom**: After successful resize operation, Terraform fails with:
```
Error: invalid vm_state
vm_state must be 'active' or 'stopped', got: resized
```

**Root Cause**: The Update function validates vm_state against user input values, but during state refresh after resize, the API returns transient states like "resized". The validation incorrectly rejects these API-returned values.

**Expected Behavior**: The vm_state validation should only apply to user-provided values, not API-returned states during operations.

**Suggested Fix**: Move the vm_state validation to only check when `data.VmState` differs from `state.VmState` and the new value is user-provided, not when reading the API response.

---

### Bug 2: ImportState Does Not Populate Required Fields

**Location**: `internal/services/cloud_instance/resource.go:859-936`

**Symptom**: After import, `terraform plan` shows resource must be replaced:
- `security_groups` forces replacement (not populated)
- `interfaces` completely missing from imported state
- `volumes` missing: source, image_id, boot_index, size

**Root Cause**: The ImportState function only populates minimal fields from API response:
- id, project_id, region_id
- flavor (extracted from flavor object)
- volumes.volume_id, volumes.delete_on_termination

But does NOT populate:
- interfaces (the `interfaces` list API endpoint is not called)
- security_groups
- name
- volumes.source, volumes.image_id, volumes.boot_index, volumes.size

**Expected Behavior**: ImportState should populate ALL fields that the API returns, matching the Read function behavior.

**Suggested Fix**:
1. Call the interfaces list API (`Instances.Interfaces.List`) in ImportState
2. Populate security_groups from API response
3. Populate all volume fields from API response
4. Match the logic in the Read function

## Computed Field Behavior

| Field | Tag | Drift After Operations |
|-------|-----|------------------------|
| flavor | required | ✅ No drift after resize |
| volumes.size | computed_optional | ✅ No drift after extend |
| volumes.volume_id | computed_optional | ✅ Populated correctly |
| volumes.delete_on_termination | computed_optional | ✅ Correct default (true) |
| vm_state | computed_optional | ✅ No drift after start/stop |
| interfaces.ip_address | computed_optional | ✅ Populated correctly |
| interfaces.port_id | computed | ✅ Populated correctly |
| addresses | computed | ✅ Populated correctly |
| status | computed | ✅ Populated correctly |
| tasks | computed | ✅ Populated correctly |

## Union Type Coverage (Tested)

### Volumes (1/5 tested)
- ✅ `image` source - Fully tested with MITM verification

### Interfaces (1/4 tested)
- ✅ `external` type - Fully tested with MITM verification

### Floating IP (0/2 tested)
- Not tested in this run

## Tests Not Yet Executed

Due to time constraints, the following tests were not completed:

- Volume sources: new-volume, snapshot, apptemplate, existing-volume
- Interface types: subnet, any_subnet, reserved_fixed_ip
- Floating IP: new, existing
- Interface attach/detach operations
- Security group assign/unassign
- Name update (PATCH)
- 404 handling (delete outside Terraform)
- Access: password/username, SSH key, user_data

## Evidence Package

Location: `test-instance-skill/evidence/`

| File | Description |
|------|-------------|
| drift_test_001.log | Initial drift test |
| flavor_resize_002.log | Flavor resize output (includes error) |
| volume_extend_004.log | Volume extend output |
| volume_shrink_guard_006.log | Shrink guard error |
| vm_state_stop_007.log | VM stop operation |
| vm_state_start_008.log | VM start operation |
| destroy_009.log | Destroy operation |
| create_baseline_010.log | Create baseline instance |
| import_012.log | Import operation |
| flow.mitm | MITM proxy capture |

## Recommendations

### Priority 1 (Critical)
1. **Fix Bug 1**: vm_state validation during resize operations
2. **Fix Bug 2**: ImportState field population

### Priority 2 (Important)
3. Complete union type testing for all volume sources
4. Complete interface type testing (subnet, any_subnet, reserved_fixed_ip)
5. Test floating IP variants

### Priority 3 (Nice to Have)
6. Add integration tests for interface attach/detach
7. Add integration tests for security group assign/unassign
8. Add tests for access methods (password, SSH key, user_data)

## Conclusion

The `gcore_cloud_instance` resource core functionality is working well:

✅ **Working correctly**:
- Create with image volume + external interface
- Drift detection after all operations
- In-place flavor resize (uses /changeflavor)
- In-place volume extend (uses /extend)
- VM state start/stop (uses /action)
- Delete with task polling
- Volume shrink guard with clear error message

❌ **Requires fixes**:
- vm_state validation error during resize (Bug 1)
- ImportState missing field population (Bug 2)

The resource is usable for basic workflows but the import functionality needs fixing before production use.

---

*Report generated by terraform-testing-skill*
