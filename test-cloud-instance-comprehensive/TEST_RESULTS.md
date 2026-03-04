# Cloud Instance Comprehensive Test Results

**Jira**: GCLOUD2-21138
**Date**: 2025-12-04 (Updated: 2025-12-05)
**Branch**: terraform-instances
**Provider**: Development build with dev_overrides
**Fix Applied**: Removed `option.WithRequestBody` from polling methods (GCLOUD2-21015)

---

## Test Summary

| Category | Tests | Status |
|----------|-------|--------|
| Interface Variants | 5 | ALL PASS |
| Volume Variants | 3 | ALL PASS |
| Update Operations | 3 | ALL PASS |
| Security Groups | 1 | PASS (with note) |
| Placement Groups | 1 | PASS |
| Floating IP | 1 | PASS |
| Drift Detection | 2 | PASS |
| Volume Deletion Behavior | 1 | PASS |
| **GET Request Body Fix** | 1 | **PASS (CRITICAL)** |

**Overall Result: ALL TESTS PASSED**

---

## Detailed Results

### Interface Variant Tests

| Test | Interface Type | Result | Instance ID |
|------|---------------|--------|-------------|
| Test 1 | `external` with `ip_family=ipv4` | PASS | 3f5b6638-92f4-44c8-a8c4-d71671dc9062 |
| Test 1b | `external` with `ip_family=dual` | PASS | b3014ebf-ae91-4915-afac-db607eaa9874 |
| Test 2 | `subnet` with network_id + subnet_id | PASS | ab3c7d2f-bd41-4385-a7c7-148c36bb7e7f |
| Test 3 | `subnet` with floating_ip source=new | PASS | ab5fa00b-5392-46f0-bfd0-6c043eb5194f |
| Test 4 | **`any_subnet`** with network_id only | **PASS** | 8584bd3a-a83a-44ae-b385-2978bde9054a |

### Volume Variant Tests

| Test | Description | Result |
|------|-------------|--------|
| Test 6 | Existing volume (gcore_cloud_volume) | PASS - All instances used pre-created volumes |
| Test 7 | Volume with attachment_tag | PASS - Used `attachment_tag = "boot-disk"` |
| Test 8 | Boot index specification | PASS - All volumes correctly set `boot_index = 0` |

### Update Operation Tests

| Test | Operation | Result | Notes |
|------|-----------|--------|-------|
| Test 9 | Name change | PASS | Changed from `tf-test-instance-updates-initial` to `tf-test-instance-updates-RENAMED` |
| Test 10 | Tags update | PASS | Added `version=v2`, changed `environment=test` to `production` |
| Test 11 | **Flavor change** | **PASS** | Resized from `g1-standard-1-2` to `g1-standard-2-4` in-place (same instance ID) |

### Security Groups (Test 14)

**Result: PASS with Important Finding**

- Security groups work correctly at **interface level**
- **Finding**: Cannot specify same security group at both instance level AND interface level
- API returns error: `"Duplicate items in the list: '9fa59dfb-df95-4860-965b-455556cbe7eb'"`
- **Recommendation**: Document that security_groups should only be specified at one level (interface or instance, not both)

### Placement Groups (Test 19)

**Result: PASS**

- Created placement group: `tf-test-placementgroup-comprehensive` with policy `affinity`
- Servergroup ID: `29a9ac26-e70f-4aa8-808d-177ca054a0de`
- Instance `tf-test-instance-placement` successfully assigned to placement group via `servergroup_id`

### Floating IP (Test 5)

**Result: PASS**

- Instance `tf-test-instance-floating` created with floating IP via `source = "new"`
- Floating IP automatically allocated and attached to subnet interface

### Drift Detection

**Result: PASS**

```
No changes. Your infrastructure matches the configuration.
Terraform has compared your real infrastructure against your configuration
and found no differences, so no changes are needed.
```

### GET Request Body Fix (CRITICAL - GCLOUD2-21015)

**Result: PASS**

**Issue**: GET requests during polling were incorrectly sending request body due to `option.WithRequestBody` being applied to all requests.

**Fix Applied**: Refactored Create/Update methods to pass body as SDK params instead of `option.WithRequestBody`.

**Verification from Terraform Debug Logs**:

**POST Request (has body - correct):**
```
POST /cloud/v2/instances/379987/76 HTTP/1.1
> content-type: application/json
...
{"flavor":"g1-standard-1-2","interfaces":[{"network_id":"...","type":"any_subnet"}],...}
```

**GET Request during polling (NO body - FIX VERIFIED):**
```
GET /cloud/v1/instances/379987/76/8584bd3a-a83a-44ae-b385-2978bde9054a HTTP/1.1
> authorization: APIKey ...
> user-agent: Gcore/Go 0.22.0
...
(no body - immediately followed by response)
```

**Flavor Change Endpoint (correct):**
```
POST /cloud/v1/instances/379987/76/8584bd3a-a83a-44ae-b385-2978bde9054a/changeflavor HTTP/1.1
{"flavor_id":"g1-standard-2-4"}
```

**Conclusion**: The fix successfully ensures GET requests have no body while POST/PATCH requests correctly include the body.

### Volume Deletion Behavior (Test 20)

**Result: PASS - Behavior Confirmed**

This test verified the QA-reported behavior change.

**Test Steps**:
1. Instance `tf-test-instance-updates` destroyed via `terraform destroy -target`
2. Checked if boot volume `a3ba7bea-4a93-405f-a5d7-58c095a57a68` still exists

**Result**: Volume still exists with status `available`

**Conclusion**: In the new provider, **volumes are NOT deleted when an instance is destroyed**. This is the expected behavior - volumes should be managed independently. This differs from the old provider which deleted volumes with instances.

---

## Issues Found

### 1. Security Group Duplicate Detection
- **Issue**: API rejects duplicate security groups across instance and interface levels
- **Impact**: Low - This is correct API behavior
- **Action**: Document in provider docs that security_groups should be specified at one level only

### 2. Volume Deletion Behavior Change (Known - Not a Bug)
- **Issue**: Volumes persist after instance destruction
- **Impact**: Expected behavior change from old provider
- **Action**: Document behavior difference in migration guide

---

## Resources Tested

### Instances Created (all reached ACTIVE status)
- `tf-test-instance-placement` - Placement group test
- `tf-test-instance-external` - External interface + attachment_tag test
- `tf-test-instance-subnet-secgroup` - Subnet interface + security groups test
- `tf-test-instance-floating` - Floating IP test
- `tf-test-instance-updates` - Update operations test

### Volumes Created
- `tf-test-boot-placement`
- `tf-test-boot-external`
- `tf-test-boot-subnet`
- `tf-test-boot-floating`
- `tf-test-boot-updates`

### Other Resources
- `tf-test-placementgroup-comprehensive` - Placement group with affinity policy

---

## Cleanup

All test resources successfully destroyed:
- 5 instances destroyed
- 5 volumes destroyed
- 1 placement group destroyed
- Total: 10 resources destroyed

---

## Configuration Used

```hcl
locals {
  image_id         = "f84ddba3-7a5a-4199-931a-250e981d16fb" # Ubuntu 25.04
  flavor           = "g1-standard-1-2"
  network_id       = "cd2c62cd-9763-4766-8d36-6066ed92b3e3"
  subnet_id        = "4f144cc7-c377-445d-9c23-fa6576f1b945"
  security_group_1 = "9fa59dfb-df95-4860-965b-455556cbe7eb" # default
}
```

---

## Union Type Coverage

### Interfaces (4 variants tested)

| Variant | Discriminator | Status | Notes |
|---------|---------------|--------|-------|
| `external` | `"type": "external"` | ✅ PASS | Tested with ipv4 and dual |
| `subnet` | `"type": "subnet"` | ✅ PASS | With network_id + subnet_id |
| `any_subnet` | `"type": "any_subnet"` | ✅ PASS | With network_id only (subnet auto-selected) |
| `reserved_fixed_ip` | `"type": "reserved_fixed_ip"` | ⚠️ Not tested | Requires pre-created reserved IP |

### Volumes (1 variant - simplified)

| Variant | Discriminator | Status | Notes |
|---------|---------------|--------|-------|
| `existing-volume` | `"source": "existing-volume"` | ✅ PASS | With pre-created gcore_cloud_volume |

### Floating IP (2 variants)

| Variant | Discriminator | Status | Notes |
|---------|---------------|--------|-------|
| `new` | `"source": "new"` | ✅ PASS | Auto-allocates floating IP |
| `existing` | `"source": "existing"` | ⚠️ Not tested | Requires pre-created floating IP |

---

## API Verification Summary

All API endpoints verified via Terraform debug logs:

| Operation | Endpoint | Body | Status |
|-----------|----------|------|--------|
| Create instance | `POST /cloud/v2/instances` | ✅ Correct JSON body | PASS |
| Read instance (polling) | `GET /cloud/v1/instances/{id}` | ✅ **NO body** | **PASS (FIX VERIFIED)** |
| Change flavor | `POST /changeflavor` | ✅ Correct JSON body | PASS |
| Delete instance | `DELETE /cloud/v1/instances/{id}` | ✅ As expected | PASS |

---

## Recommendations

1. **Documentation**: Add note about security_groups being specified at only one level (interface OR instance)
2. **Migration Guide**: Document volume deletion behavior change from old to new provider
3. **Examples**: Update examples to show best practices for security group placement
4. **Testing**: Consider adding tests for `reserved_fixed_ip` interface and `existing` floating IP source
