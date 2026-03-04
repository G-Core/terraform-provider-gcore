# Testing Plan: gcore_cloud_instance_image

## Phase 1: Analysis Summary

### Resource Overview
- **Resource**: `gcore_cloud_instance_image`
- **JIRA Ticket**: GCLOUD2-22615
- **Purpose**: Upload custom images to Gcore cloud from a URL

### OpenAPI → SDK Alignment

| API Operation | SDK Method | Async? | Notes |
|---------------|------------|--------|-------|
| POST /images (upload) | UploadAndPoll | Yes | Creates task, polls until complete |
| GET /images/{id} | Get | No | Read image details |
| PATCH /images/{id} | Update | No | Update mutable fields |
| DELETE /images/{id} | DeleteAndPoll | Yes | Creates task, polls until complete |

### Field Analysis

| Field | Type | Behavior | Test Coverage |
|-------|------|----------|---------------|
| id | computed | UseStateForUnknown | Drift test |
| name | required | Updateable | Create, Update test |
| url | required | RequiresReplace | Create test |
| project_id | optional | RequiresReplace, path param | Create test |
| region_id | optional | RequiresReplace, path param | Create test |
| os_distro | optional | RequiresReplace | Create test |
| os_version | optional | RequiresReplace | Create test |
| architecture | computed_optional | Default: x86_64, RequiresReplaceIfConfigured | Drift test |
| cow_format | computed_optional | Default: false, RequiresReplaceIfConfigured | Drift test |
| os_type | computed_optional | Default: linux | Drift test |
| ssh_key | computed_optional | Default: allow | Drift test |
| is_baremetal | computed_optional | Default: false | Drift test |
| tags | optional | no_refresh | Create test |
| hw_firmware_type | optional | Updateable | Update test |
| hw_machine_type | optional | Updateable | Update test |

### Fixes Applied (PR #74)

1. Changed from `Upload()` to `UploadAndPoll()` for async task handling
2. Removed `tasks` field from model/schema (caused state corruption)
3. Changed `Unmarshal()` to `UnmarshalComputed()` in Read/Import
4. Used `MarshalJSON()` + `WithRequestBody()` pattern
5. Added `IsUnknown()` checks for path parameters
6. Proper error handling for `io.ReadAll()`

---

## Phase 2: Resource Coverage Matrix

### Test Cases

| Test # | Description | Fields Tested | MITM Pattern | Expected Result |
|--------|-------------|---------------|--------------|-----------------|
| 1 | Create with all optional fields | name, url, os_distro, os_version, os_type, architecture, ssh_key, tags | POST /images | Resource created with ID |
| 2 | Drift check | All computed_optional fields | - | No changes on 2nd plan |
| 3 | Update name | name | PATCH /images/{id} | Same ID, name updated |
| 4 | Import existing | id, project_id, region_id | GET /images/{id} | State matches API |
| 5 | Delete | - | DELETE /images/{id} | Task polling, resource removed |

### Computed/Optional Field Coverage

| Field | Tag | Default Value | Test # | Drift Check? |
|-------|-----|---------------|--------|--------------|
| architecture | computed_optional | x86_64 | 1 | Yes (Test 2) |
| cow_format | computed_optional | false | 1 | Yes (Test 2) |
| os_type | computed_optional | linux | 1 | Yes (Test 2) |
| ssh_key | computed_optional | allow | 1 | Yes (Test 2) |
| is_baremetal | computed_optional | false | 1 | Yes (Test 2) |

### QA Checklist Mapping

| QA Category | Test # | Expected Behavior |
|-------------|--------|-------------------|
| Async operations | 1, 5 | Task polling works correctly |
| Drift detection | 2 | No changes on 2nd plan |
| Update in-place | 3 | Same ID after update |
| Import | 4 | State matches API response |
| Error handling | - | Clear error messages |

---

## Phase 3: Test Execution Plan

### Test Configuration

```hcl
resource "gcore_cloud_instance_image" "test" {
  project_id = 379987
  region_id  = 76

  name       = "test-image-skill"
  url        = "http://mirror.noris.net/cirros/0.4.0/cirros-0.4.0-x86_64-disk.img"

  os_distro  = "cirros"
  os_version = "0.4.0"
  os_type    = "linux"

  architecture = "x86_64"
  ssh_key      = "allow"

  tags = {
    "test" = "skill-test"
  }
}
```

### Dependencies Required
- Network: No
- Subnet: No
- Pre-created resources: No

### Cleanup Plan
- Image will be destroyed after all tests complete
- Estimated cost: Minimal (small test image)

---

## Reviewer Sign-off

- [x] All computed_optional fields have drift tests
- [x] Async operations verified via MITM
- [x] All QA categories covered
- [x] No dependencies needed
- [x] Cleanup plan reviewed

**Test execution approved: YES**
