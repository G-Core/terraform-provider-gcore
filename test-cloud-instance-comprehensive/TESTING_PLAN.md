# Comprehensive Testing Plan: gcore_cloud_instance

**Jira**: GCLOUD2-21138
**Date**: 2025-12-04
**Branch**: terraform-instances
**Fix Applied**: Removed `option.WithRequestBody` from polling methods - body now passed as params

---

## Phase 1: Analysis Summary

### OpenAPI → SDK Alignment

| API Operation | SDK Method | Async? | Notes |
|---------------|------------|--------|-------|
| POST /instances | NewAndPoll | Yes | Creates instance with task polling |
| GET /instances/{id} | Get | No | Should NOT have request body |
| PATCH /instances/{id} | Update | No | Updates name, tags |
| DELETE /instances/{id} | DeleteAndPoll | Yes | Deletes with task polling |
| POST /instances/{id}/changeflavor | ResizeAndPoll | Yes | Flavor resize |
| POST /instances/{id}/action | ActionAndPoll | Yes | Start/stop/reboot |
| POST /instances/{id}/attach_interface | Interfaces.AttachAndPoll | Yes | Attach new interface |
| POST /instances/{id}/detach_interface | Interfaces.DetachAndPoll | Yes | Detach interface |
| POST /instances/{id}/assignsecuritygroup | AssignSecurityGroup | No | Assign SG |
| POST /instances/{id}/unassignsecuritygroup | UnassignSecurityGroup | No | Unassign SG |

### Key Bug Fixed

**Issue**: GET requests during polling were sending request body
**Root Cause**: `option.WithRequestBody` was applied to all requests including polling GET
**Fix**: Refactored Create/Update methods to pass body as SDK params instead of `option.WithRequestBody`

---

## Phase 2: Resource Coverage Matrix

### Union Type Coverage

#### Interfaces (4 variants) - Discriminator: `type`

| Variant | Discriminator | Required Fields | Test # | Already Tested? |
|---------|---------------|-----------------|--------|-----------------|
| external | `"type": "external"` | - | 1 | ✅ Yes |
| subnet | `"type": "subnet"` | network_id, subnet_id | 2 | ✅ Yes |
| any_subnet | `"type": "any_subnet"` | network_id | 3 | ❌ No |
| reserved_fixed_ip | `"type": "reserved_fixed_ip"` | port_id | 4 | ❌ No |

#### Floating IP (2 variants) - Discriminator: `source`

| Variant | Discriminator | Required Fields | Test # | Already Tested? |
|---------|---------------|-----------------|--------|-----------------|
| new | `"source": "new"` | - | 5 | ✅ Yes |
| existing | `"source": "existing"` | existing_floating_id | 6 | ❌ No |

#### Volumes (1 variant - simplified) - Discriminator: `source`

| Variant | Discriminator | Required Fields | Test # | Already Tested? |
|---------|---------------|-----------------|--------|-----------------|
| existing-volume | `"source": "existing-volume"` | volume_id | 7 | ✅ Yes |

### Computed/Optional Field Coverage

| Field | Tag | Test # | Drift Check Needed? |
|-------|-----|--------|---------------------|
| vm_state | computed_optional | 8 | Yes |
| ip_address (interfaces) | computed_optional | 2,3,4 | Yes |
| boot_index (volumes) | computed_optional | 7 | Yes |
| status | computed | All | No |
| addresses | computed | All | No |

### Special API Endpoint Coverage (MITM Verification)

| Endpoint | Purpose | Test # | MITM Pattern to Verify |
|----------|---------|--------|------------------------|
| POST /instances | Create | 1-7 | Body has correct union structure |
| GET /instances/{id} | Polling | All | **NO request body** |
| PATCH /instances/{id} | Update | 9 | Has name/tags in body |
| POST /changeflavor | Resize | 10 | Has flavor in body |
| POST /action | Start/Stop | 8 | Has action in body |

### QA Checklist Mapping

| QA Category | Test # | Expected Behavior |
|-------------|--------|-------------------|
| Drift detection | All | No changes on 2nd plan |
| Update in-place | 9,10 | Same ID after update |
| Async operations | 1-8 | Task polling works without body in GET |
| Import | 11 | State matches API |

---

## Phase 3: Test Execution Plan

### Tests to Execute

| Test # | Description | MITM Verification | Priority |
|--------|-------------|-------------------|----------|
| 1 | External interface with ip_family | `"type": "external"` | ✅ Done |
| 2 | Subnet interface with network_id + subnet_id | `"type": "subnet"` | ✅ Done |
| 3 | **any_subnet interface** | `"type": "any_subnet"` | HIGH |
| 4 | **reserved_fixed_ip interface** | `"type": "reserved_fixed_ip"` | HIGH |
| 5 | Floating IP source=new | `"source": "new"` | ✅ Done |
| 6 | **Floating IP source=existing** | `"source": "existing"` | MEDIUM |
| 7 | Existing volume attachment | `"source": "existing-volume"` | ✅ Done |
| 8 | **vm_state change (stop/start)** | POST /action | HIGH |
| 9 | Name/tags update | PATCH body | ✅ Done |
| 10 | **Flavor change** | POST /changeflavor | HIGH |
| 11 | **MITM: Verify GET has no body** | GET during polling | CRITICAL |

### Dependencies Required

- Network: `cd2c62cd-9763-4766-8d36-6066ed92b3e3` (existing)
- Subnet: `4f144cc7-c377-445d-9c23-fa6576f1b945` (existing)
- Image: `f84ddba3-7a5a-4199-931a-250e981d16fb` (Ubuntu 25.04)
- Security Group: `9fa59dfb-df95-4860-965b-455556cbe7eb` (default)

### New Dependencies to Create

- Reserved Fixed IP (for test 4)
- Existing Floating IP (for test 6)

---

## Phase 4: Critical MITM Verification

### Primary Goal: Verify GET Requests Have No Body

After the fix, all GET requests during polling MUST NOT have a request body.

**Verification Steps**:
1. Start mitmproxy: `./scripts/run_mitm.sh`
2. Run terraform apply with proxy enabled
3. Check MITM logs for GET requests
4. Verify GET requests show `(no content)` or empty body

**Expected MITM Output**:
```
GET /cloud/v1/instances/379987/76/{instance_id} HTTP/1.1
(no content)
```

**FAILURE Indicator**:
```
GET /cloud/v1/instances/379987/76/{instance_id} HTTP/1.1
Content-Type: application/json
{"flavor": "...", "interfaces": [...], ...}
```

---

## Phase 5: Validation Checklist

### For Each Test

- [ ] Terraform apply succeeds
- [ ] Instance reaches ACTIVE status
- [ ] Terraform plan shows no drift
- [ ] MITM logs show correct payloads
- [ ] GET requests have no body during polling

### Overall

- [ ] All union variants tested
- [ ] All computed_optional fields verified for drift
- [ ] All special endpoints work correctly
- [ ] Resources cleaned up

---

## Approval Request

**Tests to run**:
1. Test 3: any_subnet interface (new)
2. Test 4: reserved_fixed_ip interface (new)
3. Test 6: Floating IP existing source (new)
4. Test 8: vm_state change (new)
5. Test 10: Flavor change (new)
6. Test 11: MITM verification of GET without body (CRITICAL)

**Estimated resources to create**:
- 3-4 instances (sequential, with cleanup between tests)
- 1 reserved fixed IP
- 1 floating IP
- Volumes as needed

**Do you approve this testing plan? (yes/no)**
