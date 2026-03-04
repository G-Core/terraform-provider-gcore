# Comprehensive Testing Plan for Load Balancer Resources

**Status**: In Progress
**Date Created**: 2025-11-04
**Purpose**: Verify configuration drift detection and proper update vs replacement behavior for Load Balancer, Listener, and Pool resources

---

## Executive Summary

This testing plan validates:
1. **Configuration drift detection** - No false positives when state matches infrastructure
2. **Update vs Replace behavior** - PATCH operations used when available instead of resource replacement
3. **Field-level updates** - Each PATCH-able field can be updated in-place
4. **Correct replacements** - Fields marked RequiresReplace force recreation when changed

---

## API Capabilities Summary

### Load Balancer (`gcore_cloud_load_balancer`)

**PATCH-able fields** (via `/cloud/v1/loadbalancers/.../PATCH`):
- ✅ `name` - Load balancer name
- ✅ `logging` - Logging configuration object
- ✅ `preferred_connectivity` - L2 or L3 connectivity
- ✅ `tags` - Key-value tags

**RequiresReplace fields** (not in PATCH, requires recreation):
- ⛔ `flavor` - LB flavor/size
- ⛔ `vip_network_id` - Network for VIP
- ⛔ `vip_subnet_id` - Subnet for VIP
- ⛔ `floating_ip` - Floating IP configuration
- ⛔ `vip_ip_family` - IP family (dual/ipv4/ipv6)
- ⛔ `vip_port_id` - Reserved fixed IP port
- ⛔ `name_template` - Name template

### Listener (`gcore_cloud_load_balancer_listener`)

**PATCH-able fields** (via `/cloud/v2/lblisteners/.../PATCH`):
- ✅ `name` - Listener name
- ✅ `allowed_cidrs` - CIDR allowlist
- ✅ `connection_limit` - Connection limit
- ✅ `secret_id` - TLS certificate secret
- ✅ `sni_secret_id` - SNI certificate secrets list
- ✅ `timeout_client_data` - Client timeout (ms)
- ✅ `timeout_member_connect` - Backend connect timeout (ms)
- ✅ `timeout_member_data` - Backend data timeout (ms)
- ✅ `user_list` - Basic auth users list

**RequiresReplace fields** (not in PATCH, requires recreation):
- ⛔ `load_balancer_id` - Parent LB
- ⛔ `protocol` - Listener protocol (HTTP/HTTPS/TCP/etc)
- ⛔ `protocol_port` - Port number
- ⛔ `insert_x_forwarded` - X-Forwarded headers flag

### Pool (`gcore_cloud_lbpool`)

**PATCH-able fields** (via `/cloud/v2/lbpools/.../PATCH`):
- ✅ `name` - Pool name
- ✅ `lb_algorithm` - Load balancing algorithm
- ✅ `protocol` - Pool protocol
- ✅ `ca_secret_id` - CA certificate secret
- ✅ `crl_secret_id` - CRL secret
- ✅ `secret_id` - TLS client cert secret
- ✅ `healthmonitor` - Health monitor config (full object)
- ✅ `session_persistence` - Session persistence config (full object)
- ✅ `members` - Pool members list (full array)
- ✅ `timeout_client_data` - Client timeout (ms)
- ✅ `timeout_member_connect` - Backend connect timeout (ms)
- ✅ `timeout_member_data` - Backend data timeout (ms)

**RequiresReplace fields** (not in PATCH, requires recreation):
- ⛔ `listener_id` - Parent listener
- ⛔ `load_balancer_id` - Parent LB

---

## Test Categories

- **Phase 1**: Drift Detection (4 tests) - Highest Priority
- **Phase 2**: Update Operations (15 tests) - High Priority
- **Phase 3**: Replace Operations (8 tests) - Medium Priority
- **Phase 4**: Combined & Edge Cases (6 tests) - Lower Priority

**Total**: 33 test cases

---

## Phase 1: Configuration Drift Detection Tests

**Goal**: Verify no false drift when applying same configuration twice

### TC-DRIFT-001: Load Balancer - No Changes
**Status**: ⏳ Pending

**Configuration**:
```hcl
resource "gcore_cloud_load_balancer" "test" {
  name       = "test-lb-drift-01"
  flavor     = "lb1-1-2"
  project_id = var.project_id
  region_id  = var.region_id

  tags = {
    environment = "test"
    purpose     = "drift-detection"
  }
}
```

**Steps**:
1. `terraform apply -auto-approve`
2. Wait for LB to become ACTIVE
3. `terraform plan`

**Expected**: No changes detected - output should show "No changes. Your infrastructure matches the configuration."

**Verification**:
```bash
terraform plan | grep "No changes" && echo "PASS" || echo "FAIL"
```

**Result**:

---

### TC-DRIFT-002: Listener - No Changes (All Optional Fields)
**Status**: ⏳ Pending

**Configuration**:
```hcl
resource "gcore_cloud_load_balancer_listener" "test" {
  name              = "test-listener-drift-02"
  load_balancer_id  = gcore_cloud_load_balancer.test.id
  protocol          = "HTTP"
  protocol_port     = 80
  project_id        = var.project_id
  region_id         = var.region_id

  allowed_cidrs          = ["0.0.0.0/0"]
  connection_limit       = 100000
  timeout_client_data    = 50000
  timeout_member_connect = 5000
  timeout_member_data    = 50000
}
```

**Steps**: Same as TC-DRIFT-001

**Expected**: No changes detected

**Result**:

---

### TC-DRIFT-003: Pool - No Changes (With Health Monitor)
**Status**: ⏳ Pending

**Configuration**:
```hcl
resource "gcore_cloud_lbpool" "test" {
  name            = "test-pool-drift-03"
  lb_algorithm    = "ROUND_ROBIN"
  protocol        = "HTTP"
  loadbalancer_id = gcore_cloud_load_balancer.test.id
  project_id      = var.project_id
  region_id       = var.region_id

  healthmonitor = {
    delay          = 10
    max_retries    = 3
    timeout        = 5
    type           = "HTTP"
    url_path       = "/health"
    expected_codes = "200"
  }

  timeout_client_data    = 50000
  timeout_member_connect = 5000
  timeout_member_data    = 50000
}
```

**Steps**: Same as TC-DRIFT-001

**Expected**: No changes detected

**Result**:

---

### TC-DRIFT-004: Pool - No Changes (With Members)
**Status**: ⏳ Pending

**Configuration**:
```hcl
resource "gcore_cloud_lbpool" "test" {
  name            = "test-pool-drift-04"
  lb_algorithm    = "ROUND_ROBIN"
  protocol        = "HTTP"
  loadbalancer_id = gcore_cloud_load_balancer.test.id
  project_id      = var.project_id
  region_id       = var.region_id

  members = [
    {
      address        = "192.168.1.10"
      protocol_port  = 8080
      subnet_id      = var.subnet_id
      weight         = 1
      admin_state_up = true
      backup         = false
    },
    {
      address        = "192.168.1.11"
      protocol_port  = 8080
      subnet_id      = var.subnet_id
      weight         = 1
      admin_state_up = true
      backup         = false
    }
  ]
}
```

**Steps**: Same as TC-DRIFT-001

**Expected**: No changes detected

**Result**:

---

## Phase 2: Update Operations (PATCH) Tests

**Goal**: Verify fields are updated via PATCH, not replaced

### TC-UPDATE-LB-001: Update Load Balancer Name
**Status**: ⏳ Pending

**Initial Configuration**:
```hcl
resource "gcore_cloud_load_balancer" "test" {
  name       = "test-lb-update-01-v1"
  flavor     = "lb1-1-2"
  project_id = var.project_id
  region_id  = var.region_id
}
```

**Updated Configuration**:
```hcl
resource "gcore_cloud_load_balancer" "test" {
  name       = "test-lb-update-01-v2"  # Changed
  flavor     = "lb1-1-2"
  project_id = var.project_id
  region_id  = var.region_id
}
```

**Expected Behavior**:
- Plan shows: `~ name: "test-lb-update-01-v1" => "test-lb-update-01-v2"`
- Update in-place (no `-/+ must be replaced`)
- API call: `PATCH /cloud/v1/loadbalancers/.../`

**Verification**:
```bash
terraform plan | grep "must be replaced" && echo "FAIL" || echo "PASS"
```

**Result**:

---

### TC-UPDATE-LB-002: Update Load Balancer Tags
**Status**: ⏳ Pending

**Initial**:
```hcl
tags = {
  environment = "dev"
  team        = "platform"
}
```

**Updated**:
```hcl
tags = {
  environment = "staging"  # Modified
  team        = "platform" # Unchanged
  owner       = "john"     # Added
}
```

**Expected**: Update in-place via PATCH

**Result**:

---

### TC-UPDATE-LISTENER-001: Update Listener Name
**Status**: ⏳ Pending

**Initial**: `name = "listener-v1"`
**Updated**: `name = "listener-v2"`

**Expected**: Update in-place via PATCH

**Result**:

---

### TC-UPDATE-LISTENER-002: Update Listener Allowed CIDRs
**Status**: ⏳ Pending

**Initial**: `allowed_cidrs = ["0.0.0.0/0"]`
**Updated**: `allowed_cidrs = ["10.0.0.0/8", "192.168.0.0/16"]`

**Expected**: Update in-place via PATCH

**Result**:

---

### TC-UPDATE-LISTENER-003: Update Listener Connection Limit
**Status**: ⏳ Pending

**Initial**: `connection_limit = 100000`
**Updated**: `connection_limit = 50000`

**Expected**: Update in-place via PATCH

**Result**:

---

### TC-UPDATE-LISTENER-004: Update Listener Timeouts
**Status**: ⏳ Pending

**Initial**:
```hcl
timeout_client_data    = 50000
timeout_member_connect = 5000
timeout_member_data    = 50000
```

**Updated**:
```hcl
timeout_client_data    = 60000
timeout_member_connect = 10000
timeout_member_data    = 60000
```

**Expected**: Update in-place via PATCH

**Result**:

---

### TC-UPDATE-POOL-001: Update Pool Name
**Status**: ⏳ Pending

**Initial**: `name = "pool-v1"`
**Updated**: `name = "pool-v2"`

**Expected**: Update in-place via PATCH

**Result**:

---

### TC-UPDATE-POOL-002: Update Pool Algorithm
**Status**: ⏳ Pending

**Initial**: `lb_algorithm = "ROUND_ROBIN"`
**Updated**: `lb_algorithm = "LEAST_CONNECTIONS"`

**Expected**: Update in-place via PATCH

**Result**:

---

### TC-UPDATE-POOL-003: Update Pool Protocol
**Status**: ⏳ Pending

**Initial**: `protocol = "HTTP"`
**Updated**: `protocol = "HTTPS"`

**Expected**: Update in-place via PATCH

**Result**:

---

### TC-UPDATE-POOL-004: Update Pool Health Monitor
**Status**: ⏳ Pending

**Initial**:
```hcl
healthmonitor = {
  delay       = 10
  max_retries = 3
  timeout     = 5
  type        = "HTTP"
  url_path    = "/health"
}
```

**Updated**:
```hcl
healthmonitor = {
  delay          = 15
  max_retries    = 5
  timeout        = 10
  type           = "HTTP"
  url_path       = "/healthz"
  expected_codes = "200,204"
}
```

**Expected**: Update in-place via PATCH

**Result**:

---

### TC-UPDATE-POOL-005: Add Pool Member
**Status**: ⏳ Pending

**Initial**: 1 member
**Updated**: 2 members (add one)

**Expected**: Update in-place via PATCH

**Result**:

---

### TC-UPDATE-POOL-006: Remove Pool Member
**Status**: ⏳ Pending

**Initial**: 2 members
**Updated**: 1 member (remove one)

**Expected**: Update in-place via PATCH

**Result**:

---

### TC-UPDATE-POOL-007: Update Pool Member Weight
**Status**: ⏳ Pending

**Initial**: `weight = 1`
**Updated**: `weight = 10`

**Expected**: Update in-place via PATCH

**Result**:

---

### TC-UPDATE-POOL-008: Update Pool Session Persistence
**Status**: ⏳ Pending

**Initial**: No session persistence

**Updated**:
```hcl
session_persistence = {
  type        = "HTTP_COOKIE"
  cookie_name = "SERVERID"
}
```

**Expected**: Update in-place via PATCH

**Result**:

---

### TC-UPDATE-POOL-009: Update Pool Timeouts
**Status**: ⏳ Pending

**Initial**:
```hcl
timeout_client_data    = 50000
timeout_member_connect = 5000
timeout_member_data    = 50000
```

**Updated**:
```hcl
timeout_client_data    = 60000
timeout_member_connect = 10000
timeout_member_data    = 60000
```

**Expected**: Update in-place via PATCH

**Result**:

---

## Phase 3: Replace Operations (RequiresReplace) Tests

**Goal**: Verify fields that SHOULD force replacement do so correctly

### TC-REPLACE-LB-001: Change Load Balancer Flavor
**Status**: ⏳ Pending

**Initial**: `flavor = "lb1-1-2"`
**Updated**: `flavor = "lb1-2-4"`

**Expected Behavior**:
- Plan shows: `-/+ gcore_cloud_load_balancer.test must be replaced`
- Reason: `flavor` changes require replacement

**Verification**:
```bash
terraform plan | grep "must be replaced" || echo "FAIL: Should force replacement"
```

**Result**:

---

### TC-REPLACE-LB-002: Change VIP Network
**Status**: ⏳ Pending

**Initial**: `vip_network_id = "network-1-uuid"`
**Updated**: `vip_network_id = "network-2-uuid"`

**Expected**: Force replacement

**Result**:

---

### TC-REPLACE-LISTENER-001: Change Listener Protocol
**Status**: ⏳ Pending

**Initial**: `protocol = "HTTP"`
**Updated**: `protocol = "HTTPS"`

**Expected**: Force replacement (protocol cannot be changed)

**Result**:

---

### TC-REPLACE-LISTENER-002: Change Listener Port
**Status**: ⏳ Pending

**Initial**: `protocol_port = 80`
**Updated**: `protocol_port = 8080`

**Expected**: Force replacement (port cannot be changed)

**Result**:

---

### TC-REPLACE-LISTENER-003: Change Load Balancer ID
**Status**: ⏳ Pending

**Initial**: `load_balancer_id = gcore_cloud_load_balancer.test1.id`
**Updated**: `load_balancer_id = gcore_cloud_load_balancer.test2.id`

**Expected**: Force replacement (cannot move listener to different LB)

**Result**:

---

### TC-REPLACE-POOL-001: Change Listener ID
**Status**: ⏳ Pending

**Initial**: `listener_id = gcore_cloud_load_balancer_listener.test1.id`
**Updated**: `listener_id = gcore_cloud_load_balancer_listener.test2.id`

**Expected**: Force replacement (cannot move pool to different listener)

**Result**:

---

### TC-REPLACE-POOL-002: Change Load Balancer ID
**Status**: ⏳ Pending

**Initial**: `loadbalancer_id = gcore_cloud_load_balancer.test1.id`
**Updated**: `loadbalancer_id = gcore_cloud_load_balancer.test2.id`

**Expected**: Force replacement (cannot move pool to different LB)

**Result**:

---

## Phase 4: Combined Resource & Edge Case Tests

### TC-COMBINED-001: Full Stack Creation
**Status**: ⏳ Pending

**Configuration**: LB + Listener + Pool with health monitor + members

**Steps**:
1. Apply full configuration
2. Verify all resources created
3. Apply again - verify no drift
4. Update LB name - verify only LB updated
5. Update listener timeouts - verify only listener updated
6. Add pool member - verify only pool updated

**Expected**: Each resource updates independently via PATCH

**Result**:

---

### TC-COMBINED-002: Update Cascade Prevention
**Status**: ⏳ Pending

**Scenario**: Update parent LB name, verify children not affected

**Expected**:
- Only LB updated (PATCH)
- Listener and Pool not touched
- No cascading updates

**Result**:

---

### TC-EDGE-001: Empty to Populated (Computed Fields)
**Status**: ⏳ Pending

**Scenario**: Initially don't specify computed_optional timeout fields, then explicitly set them to API defaults

**Expected**:
- If value matches API default, no update
- If value differs, update via PATCH

**Result**:

---

### TC-EDGE-002: Large Member List (50+ members)
**Status**: ⏳ Pending

**Scenario**: Pool with 50+ members

**Expected**:
- No drift on re-apply
- Updates handled via PATCH

**Result**:

---

## Test Execution Summary

### Phase 1: Drift Detection
| Test ID | Status | Result | Notes |
|---------|--------|--------|-------|
| TC-DRIFT-001 | ⏳ Pending | | LB drift check |
| TC-DRIFT-002 | ⏳ Pending | | Listener drift check |
| TC-DRIFT-003 | ⏳ Pending | | Pool + healthmonitor drift |
| TC-DRIFT-004 | ⏳ Pending | | Pool + members drift |

### Phase 2: Update Operations
| Test ID | Status | Result | Notes |
|---------|--------|--------|-------|
| TC-UPDATE-LB-001 | ⏳ Pending | | LB name |
| TC-UPDATE-LB-002 | ⏳ Pending | | LB tags |
| TC-UPDATE-LISTENER-001 | ⏳ Pending | | Listener name |
| TC-UPDATE-LISTENER-002 | ⏳ Pending | | Listener CIDRs |
| TC-UPDATE-LISTENER-003 | ⏳ Pending | | Listener conn limit |
| TC-UPDATE-LISTENER-004 | ⏳ Pending | | Listener timeouts |
| TC-UPDATE-POOL-001 | ⏳ Pending | | Pool name |
| TC-UPDATE-POOL-002 | ⏳ Pending | | Pool algorithm |
| TC-UPDATE-POOL-003 | ⏳ Pending | | Pool protocol |
| TC-UPDATE-POOL-004 | ⏳ Pending | | Pool healthmonitor |
| TC-UPDATE-POOL-005 | ⏳ Pending | | Pool add member |
| TC-UPDATE-POOL-006 | ⏳ Pending | | Pool remove member |
| TC-UPDATE-POOL-007 | ⏳ Pending | | Pool member weight |
| TC-UPDATE-POOL-008 | ⏳ Pending | | Pool session persist |
| TC-UPDATE-POOL-009 | ⏳ Pending | | Pool timeouts |

### Phase 3: Replace Operations
| Test ID | Status | Result | Notes |
|---------|--------|--------|-------|
| TC-REPLACE-LB-001 | ⏳ Pending | | LB flavor |
| TC-REPLACE-LB-002 | ⏳ Pending | | LB vip_network |
| TC-REPLACE-LISTENER-001 | ⏳ Pending | | Listener protocol |
| TC-REPLACE-LISTENER-002 | ⏳ Pending | | Listener port |
| TC-REPLACE-LISTENER-003 | ⏳ Pending | | Listener LB ID |
| TC-REPLACE-POOL-001 | ⏳ Pending | | Pool listener ID |
| TC-REPLACE-POOL-002 | ⏳ Pending | | Pool LB ID |

### Phase 4: Combined & Edge Cases
| Test ID | Status | Result | Notes |
|---------|--------|--------|-------|
| TC-COMBINED-001 | ⏳ Pending | | Full stack |
| TC-COMBINED-002 | ⏳ Pending | | Update cascade |
| TC-EDGE-001 | ⏳ Pending | | Computed fields |
| TC-EDGE-002 | ⏳ Pending | | Large member list |

---

## Success Metrics

- ✅ **0% false drift** - No changes reported when running `terraform plan` twice on same config
- ✅ **100% PATCH usage** - All updateable field changes use PATCH, not replace
- ✅ **Correct replacement** - Fields marked RequiresReplace do force replacement
- ✅ **No cascading updates** - Updating parent doesn't unnecessarily update children

---

## Issues Found

*(This section will be populated during test execution)*

---

## Recommendations

*(This section will be populated after test completion)*
