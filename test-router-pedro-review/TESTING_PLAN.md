# Testing Plan: gcore_cloud_network_router (GCLOUD2-21144)

## Context: Pedro's PR #38 Review Fixes

This test validates fixes applied from Pedro's code review:
1. ✅ Removed excessive comments and IMPORTANT/CRITICAL labels
2. ✅ Fixed partial route deletion condition: `len(dataRoutes) < len(stateRoutes)`
3. ✅ Added planned routes preservation after early PATCH
4. ⚠️  Need to verify: Does `routes=[]` serialize without WithJSONSet?

---

## Phase 1: Analysis Summary

### OpenAPI → SDK Alignment

| API Operation | SDK Method | Async? | Notes |
|---------------|------------|--------|-------|
| POST /routers | NewAndPoll | Yes | ✅ Implemented |
| PATCH /routers/{id} | Update | No | ✅ Implemented |
| DELETE /routers/{id} | DeleteAndPoll | Yes | ✅ Implemented |
| POST /routers/{id}/attach | AttachSubnet | No | ✅ Implemented |
| POST /routers/{id}/detach | DetachSubnet | No | ✅ Implemented |

### Legacy Provider Comparison

| Feature | Old Provider | New Provider | Parity |
|---------|-------------|--------------|--------|
| Create | WaitTaskAndReturn | NewAndPoll | ✅ |
| Update | HasChange logic | MarshalForPatch | ✅ |
| Interface attach | attach_interface | AttachSubnet | ✅ |
| Interface detach | detach_interface | DetachSubnet | ✅ |
| Routes clearing | Empty array | Empty array + ModifyPlan | ✅ |

### Migration Nuances

- ✅ **Interface management**: Uses attach/detach endpoints (not PATCH)
- ✅ **Routes clearing**: ModifyPlan detects removal and sets to empty array
- ✅ **Partial deletions**: NEW - now correctly detected with `len(dataRoutes) < len(stateRoutes)`
- ⚠️  **WithJSONSet usage**: May be unnecessary if routes=[] serializes automatically

### JIRA/Known Issues

- **GCLOUD2-21144**: Router routes - apply Pedro's PR #38 review comments
  - Partial deletion bug fixed
  - Code comments cleaned up
  - Need to verify empty slice serialization

---

## Phase 2: Resource Coverage Matrix

### Computed/Optional Field Coverage

| Field | Tag | Default Value | Test # | Drift Check? |
|-------|-----|---------------|--------|--------------|
| interfaces | computed_optional | [] | 1,2,3 | Yes |
| routes | computed_optional | [] | 1,4,5,6 | Yes |
| external_gateway_info | computed_optional | null | 1 | Yes |
| type (interface) | computed_optional | subnet | 2 | Yes |
| type (gateway) | computed_optional | manual | 1 | Yes |
| enable_snat | computed_optional | false | 1 | Yes |

### Special API Endpoint Coverage

| Endpoint | Purpose | Test # | MITM Pattern |
|----------|---------|--------|--------------|
| POST /routers | Create | 1 | Check payload has name |
| PATCH /routers/{id} | Update name | 7 | Check PATCH body |
| POST /routers/{id}/attach | Attach interface | 2 | Verify subnet_id in body |
| POST /routers/{id}/detach | Detach interface | 3 | Verify subnet_id in body |
| PATCH /routers/{id} (routes) | Update routes | 4,5,6 | Verify routes in body |
| PATCH /routers/{id} (routes=[]) | Delete all routes | 6 | **CRITICAL**: Verify `"routes":[]` |
| GET /routers/{id} | Final state check | All | Verify after each operation |

### QA Checklist Mapping

| QA Category | Test # | Expected Behavior |
|-------------|--------|-------------------|
| Drift detection | 1-10 | No changes on 2nd plan |
| Update in-place | 7 | Same ID after name change |
| Async operations | 1 | Task polling works |
| Empty array handling (routes) | 6a,6b | Sends [] not null or omitted |
| Empty array handling (interfaces) | 8,9 | Sends detach request |
| Partial deletion (NEW) | 5 | 3 routes → 1 route works |
| Full deletion - explicit [] | 6a,8 | Attribute set to [] in config |
| Full deletion - attr removed | 6b,9 | Attribute removed from config |
| Interface operations | 2,3,8,9 | Attach/detach work correctly |
| Import | 10 | State matches API |
| ModifyPlan behavior | 6b,9 | Detects removed attrs → [] |

---

## MITM Request Body Logging

**NEW**: The `run_mitm.sh` script now captures request/response bodies to `mitm_requests.log`

**How to view:**
```bash
# View all requests
cat mitm_requests.log

# Find specific operation
grep "PATCH.*routers" mitm_requests.log -A 30

# Find route deletions
grep -B 5 '"routes"' mitm_requests.log | grep -A 30 "PATCH"
```

**Critical verification points:**
- Test #6a: `"routes": []` must appear in PATCH body
- Test #6b: `"routes": []` must appear in PATCH body (ModifyPlan effect)
- Tests #8,9: Detach requests for interfaces

---

## Phase 3: Test Execution Plan

### Test #1: Minimal Router (Drift Test)
**Purpose**: Verify no drift with minimal config, test computed fields

**Config**:
```hcl
resource "gcore_cloud_network_router" "test" {
  name = "pedro-test-minimal"
}
```

**Expected**:
- Create succeeds
- `terraform plan` shows no changes
- Computed fields populated: status, distributed, region

**MITM Verification**:
- POST /routers with `{"name": "pedro-test-minimal"}`
- Response has all computed fields

---

### Test #2: Add Interface (Attach Operation)
**Purpose**: Verify AttachSubnet endpoint is used

**Config**:
```hcl
resource "gcore_cloud_network_router" "test" {
  name = "pedro-test-minimal"

  interfaces = [{
    subnet_id = gcore_cloud_subnet.test.id
  }]
}
```

**Expected**:
- Update in-place (same router ID)
- POST to /routers/{id}/attach endpoint
- No drift after apply

**MITM Verification**:
- POST /routers/{id}/attach with `{"subnet_id": "..."}`
- Response includes interface in list

---

### Test #3: Remove Interface (Detach Operation)
**Purpose**: Verify DetachSubnet endpoint is used

**Config**:
```hcl
resource "gcore_cloud_network_router" "test" {
  name = "pedro-test-minimal"

  interfaces = []
}
```

**Expected**:
- Update in-place
- POST to /routers/{id}/detach endpoint
- No drift after apply

**MITM Verification**:
- POST /routers/{id}/detach with `{"subnet_id": "..."}`
- Response shows empty interfaces list

---

### Test #4: Add Routes
**Purpose**: Verify routes can be added

**Config**:
```hcl
resource "gcore_cloud_network_router" "test" {
  name = "pedro-test-minimal"

  routes = [
    {
      destination = "10.0.1.0/24"
      nexthop     = "192.168.1.1"
    },
    {
      destination = "10.0.2.0/24"
      nexthop     = "192.168.1.2"
    },
    {
      destination = "10.0.3.0/24"
      nexthop     = "192.168.1.3"
    }
  ]
}
```

**Expected**:
- Update in-place
- PATCH to /routers/{id} with routes array
- No drift after apply

**MITM Verification**:
- PATCH /routers/{id} with `"routes": [{...}, {...}, {...}]`

---

### Test #5: Partial Route Deletion (PEDRO'S BUG FIX)
**Purpose**: Verify fix for partial deletion condition

**Config**:
```hcl
resource "gcore_cloud_network_router" "test" {
  name = "pedro-test-minimal"

  routes = [
    {
      destination = "10.0.1.0/24"
      nexthop     = "192.168.1.1"
    }
  ]
}
```

**Expected**:
- Update in-place
- PATCH to /routers/{id} with single route
- **OLD BUG**: Would NOT detect this as deletion (condition was `len(data) == 0`)
- **NEW FIX**: DOES detect as deletion (condition is `len(data) < len(state)`)
- No drift after apply

**MITM Verification**:
- PATCH /routers/{id} with `"routes": [{...}]` (only 1 route, not 3)

---

### Test #6a: Full Route Deletion - Method 1 (Empty Array)
**Purpose**: Verify routes=[] is sent correctly (Pedro's claim about omitzero)

**Config**:
```hcl
resource "gcore_cloud_network_router" "test" {
  name = "pedro-test-minimal"

  routes = []  # Explicitly set to empty array
}
```

**Expected**:
- Update in-place
- PATCH to /routers/{id} with empty routes array
- No drift after apply

**MITM Verification** (CRITICAL):
- Check `mitm_requests.log` for PATCH body
- PATCH body MUST contain `"routes": []`
- NOT `"routes": null`
- NOT omitted entirely
- This verifies if WithJSONSet is necessary

---

### Test #6b: Full Route Deletion - Method 2 (Attribute Removed)
**Purpose**: Verify attribute removal from config works (ModifyPlan should set to [])

**Config**:
```hcl
resource "gcore_cloud_network_router" "test" {
  name = "pedro-test-minimal"

  # routes attribute completely removed from config
}
```

**Expected**:
- ModifyPlan detects removal and sets routes to []
- Update in-place
- PATCH to /routers/{id} with empty routes array
- No drift after apply

**MITM Verification** (CRITICAL):
- Check `mitm_requests.log` for PATCH body
- PATCH body MUST contain `"routes": []` (ModifyPlan should force this)
- Verify ModifyPlan is working correctly for computed_optional fields

---

### Test #7: Name Update (with empty routes from 6b)
**Purpose**: Verify simple updates work and preserve ID

**Config**:
```hcl
resource "gcore_cloud_network_router" "test" {
  name = "pedro-test-renamed"

  # routes still removed from config
}
```

**Expected**:
- Update in-place (same router ID)
- PATCH to /routers/{id} with new name
- No drift after apply

**MITM Verification**:
- PATCH /routers/{id} with `{"name": "pedro-test-renamed"}`

---

### Test #8: Interface Deletion - Method 1 (Empty Array)
**Purpose**: Test interface removal with explicit empty array

**Config**:
```hcl
resource "gcore_cloud_network_router" "test" {
  name = "pedro-test-renamed"

  interfaces = []  # Explicitly empty
}
```

**Expected**:
- POST to /routers/{id}/detach endpoint
- No drift after apply

**MITM Verification**:
- Check `mitm_requests.log` for detach call

---

### Test #9: Interface Deletion - Method 2 (Attribute Removed)
**Purpose**: Test interface removal by removing attribute from config

**Config**:
```hcl
resource "gcore_cloud_network_router" "test" {
  name = "pedro-test-renamed"

  # interfaces attribute completely removed
}
```

**Expected**:
- POST to /routers/{id}/detach endpoint
- No drift after apply

**MITM Verification**:
- Check `mitm_requests.log` for detach call
- Same behavior as Test #8

---

### Test #10: Import
**Purpose**: Verify import works correctly

**Command**:
```bash
terraform import gcore_cloud_network_router.test PROJECT_ID/REGION_ID/ROUTER_ID
```

**Expected**:
- State populated correctly
- No drift after import

---

## Dependencies Required

- ✅ Network: Will create new network for testing
- ✅ Subnet: Will create new subnet in network
- ❌ Pre-created resources: None needed
- ✅ Valid CIDR blocks: Using 192.168.1.0/24 for network, 10.0.x.0/24 for routes

---

## Cleanup Plan

Resources to destroy (in order):
1. Router (will auto-detach interfaces)
2. Subnet
3. Network

Estimated cost: Minimal (router is free, network/subnet are free)

---

## Reviewer Sign-off

- [x] All computed_optional fields have drift tests
- [x] All special endpoints have MITM verification
- [x] All QA categories covered
- [x] Dependencies identified and available
- [x] Cleanup plan reviewed
- [x] Pedro's bug fix (partial deletion) has dedicated test (Test #5)
- [x] Empty array serialization will be verified via MITM (Tests #6a, #6b)
- [x] Both attribute deletion methods tested ([] and removal)
- [x] MITM now captures request bodies to mitm_requests.log

**Critical Tests**:
- **Test #5**: Validates Pedro's partial deletion fix (`len(data) < len(state)`)
- **Test #6a**: Validates explicit `routes = []` in config
- **Test #6b**: Validates removed routes attribute (ModifyPlan should set to [])
- **Tests #8/9**: Validates both interface deletion methods

**Total Tests**: 10 comprehensive tests

**Approve this plan? (Waiting for user approval...)**
