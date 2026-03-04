# Comprehensive LB Resources Test Plan

## Phase 1 Analysis Summary

### Resources Under Test
| Resource | File | Async Operations | Import Support |
|----------|------|------------------|----------------|
| `gcore_cloud_load_balancer` | `internal/services/cloud_load_balancer/` | NewAndPoll, ResizeAndPoll, DeleteAndPoll | Yes (custom) |
| `gcore_cloud_load_balancer_listener` | `internal/services/cloud_load_balancer_listener/` | NewAndPoll, UpdateAndPoll, DeleteAndPoll | Yes |
| `gcore_cloud_load_balancer_pool` | `internal/services/cloud_load_balancer_pool/` | NewAndPoll, UpdateAndPoll, DeleteAndPoll | Yes (custom) |
| `gcore_cloud_load_balancer_pool_member` | `internal/services/cloud_load_balancer_pool_member/` | NewAndPoll, DeleteAndPoll | No |

### Key Implementation Details

**Load Balancer:**
- Flavor changes use dedicated `ResizeAndPoll` endpoint (not regular Update)
- PATCH returns empty `vrrp_ips`, requires GET refresh
- ImportState extracts `flavor.flavor_name` from nested API object
- Fields with `no_refresh`: flavor, name_template, vip_network_id, vip_subnet_id, floating_ip, tags

**Listener:**
- Standard CRUD with async polling
- `insert_x_forwarded` has `no_refresh` tag
- Computed optional: timeout fields, connection_limit, sni_secret_id, user_list

**Pool:**
- Health monitor deletion uses dedicated `HealthMonitors.Delete` endpoint
- ImportState extracts `listener_id` from `listeners` relationship
- `UseEmptyListWhenConfigNull` modifier for members list
- `listener_id` and `load_balancer_id` have `RequiresReplace()` + `UseStateForUnknown()`

**Pool Member (standalone resource):**
- Update rebuilds entire pool members list via pool PATCH
- Read fetches pool and finds member in members array
- No import support

---

## Phase 2: Test Coverage Matrix

### Test Group 1: Load Balancer Basic CRUD
| Test ID | Operation | Description | Expected API Calls |
|---------|-----------|-------------|-------------------|
| LB-001 | Create | Create LB with minimal config | POST /cloud/loadbalancers |
| LB-002 | Read | Refresh state | GET /cloud/loadbalancers/{id} |
| LB-003 | Update Name | Update LB name | PATCH /cloud/loadbalancers/{id} |
| LB-004 | Delete | Destroy LB | DELETE /cloud/loadbalancers/{id} |

### Test Group 2: Load Balancer Import & Drift
| Test ID | Operation | Description | Expected Outcome |
|---------|-----------|-------------|------------------|
| LB-010 | Import | Import existing LB | State populated with all fields |
| LB-011 | Plan after Import | Run plan after import | No changes (flavor extracted correctly) |
| LB-012 | Import with Flavor | Import LB, verify flavor | flavor = "lb1-1-2" in state |

### Test Group 3: Load Balancer Resize
| Test ID | Operation | Description | Expected API Calls |
|---------|-----------|-------------|-------------------|
| LB-020 | Resize | Change flavor lb1-1-2 -> lb1-2-4 | POST /cloud/loadbalancers/{id}/resize |
| LB-021 | Resize Back | Change flavor lb1-2-4 -> lb1-1-2 | POST /cloud/loadbalancers/{id}/resize |

### Test Group 4: Listener Basic CRUD
| Test ID | Operation | Description | Expected API Calls |
|---------|-----------|-------------|-------------------|
| LS-001 | Create HTTP | Create HTTP listener | POST /cloud/loadbalancers/listeners |
| LS-002 | Create TCP | Create TCP listener | POST /cloud/loadbalancers/listeners |
| LS-003 | Update Name | Update listener name | PATCH /cloud/loadbalancers/listeners/{id} |
| LS-004 | Update Timeouts | Change timeout values | PATCH /cloud/loadbalancers/listeners/{id} |
| LS-005 | Delete | Destroy listener | DELETE /cloud/loadbalancers/listeners/{id} |

### Test Group 5: Listener Import
| Test ID | Operation | Description | Expected Outcome |
|---------|-----------|-------------|------------------|
| LS-010 | Import | Import existing listener | State populated |
| LS-011 | Plan after Import | Run plan after import | No changes |

### Test Group 6: Pool Basic CRUD
| Test ID | Operation | Description | Expected API Calls |
|---------|-----------|-------------|-------------------|
| PL-001 | Create with listener_id | Create pool attached to listener | POST /cloud/loadbalancers/pools |
| PL-002 | Create with load_balancer_id | Create pool attached to LB | POST /cloud/loadbalancers/pools |
| PL-003 | Update Name | Update pool name | PATCH /cloud/loadbalancers/pools/{id} |
| PL-004 | Update Algorithm | Change lb_algorithm | PATCH /cloud/loadbalancers/pools/{id} |
| PL-005 | Delete | Destroy pool | DELETE /cloud/loadbalancers/pools/{id} |

### Test Group 7: Pool Import & Drift (CRITICAL - GCLOUD2-20778)
| Test ID | Operation | Description | Expected Outcome |
|---------|-----------|-------------|------------------|
| PL-010 | Import with listener_id | Import pool created with listener_id | listener_id populated from listeners[] |
| PL-011 | Plan after Import | Run plan after import | **NO CHANGES** (no replacement) |
| PL-012 | Import with load_balancer_id | Import pool created with lb_id | State populated |
| PL-013 | Plan after Import (lb_id) | Run plan after import | No changes |

### Test Group 8: Pool Health Monitor
| Test ID | Operation | Description | Expected API Calls |
|---------|-----------|-------------|-------------------|
| HM-001 | Add Health Monitor | Add healthmonitor block | PATCH /cloud/loadbalancers/pools/{id} |
| HM-002 | Update Health Monitor | Change delay/timeout | PATCH /cloud/loadbalancers/pools/{id} |
| HM-003 | Delete Health Monitor | Remove healthmonitor block | DELETE /cloud/loadbalancers/pools/{id}/healthmonitor |
| HM-004 | Re-add Health Monitor | Add healthmonitor again | PATCH /cloud/loadbalancers/pools/{id} |

### Test Group 9: Pool Members (inline)
| Test ID | Operation | Description | Expected API Calls |
|---------|-----------|-------------|-------------------|
| PM-001 | Add Members | Add members block with 2 members | PATCH /cloud/loadbalancers/pools/{id} |
| PM-002 | Remove One Member | Remove 1 of 2 members | PATCH /cloud/loadbalancers/pools/{id} |
| PM-003 | Remove All Members | Set members = [] | PATCH /cloud/loadbalancers/pools/{id} |
| PM-004 | Omit Members | Remove members attribute entirely | No API call (no drift) |
| PM-005 | Re-add Members | Add members back | PATCH /cloud/loadbalancers/pools/{id} |

### Test Group 10: Pool Session Persistence
| Test ID | Operation | Description | Expected API Calls |
|---------|-----------|-------------|-------------------|
| SP-001 | Add SOURCE_IP | Add session_persistence block | PATCH /cloud/loadbalancers/pools/{id} |
| SP-002 | Change to HTTP_COOKIE | Change persistence type | PATCH /cloud/loadbalancers/pools/{id} |
| SP-003 | Remove Persistence | Remove session_persistence block | PATCH /cloud/loadbalancers/pools/{id} |

### Test Group 11: Standalone Pool Member Resource
| Test ID | Operation | Description | Expected API Calls |
|---------|-----------|-------------|-------------------|
| SM-001 | Create Member | Create standalone member | POST /cloud/loadbalancers/pools/{id}/members |
| SM-002 | Update Weight | Change member weight | PATCH /cloud/loadbalancers/pools/{id} (rebuilds list) |
| SM-003 | Toggle Admin State | Set admin_state_up=false | PATCH /cloud/loadbalancers/pools/{id} |
| SM-004 | Delete Member | Destroy member | DELETE /cloud/loadbalancers/pools/{pool_id}/members/{id} |

---

## Phase 3: Test Environment Setup

### Prerequisites
- mitmproxy running on port 9092
- Environment variables configured
- Provider built with local changes

### Test Configuration
```hcl
# Provider config
terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

# Test constants
locals {
  project_id = 379987
  region_id  = 76  # Luxembourg-2
}
```

---

## Phase 4: Execution Order

1. **Setup Phase**
   - Create base LB (LB-001)
   - Create listener (LS-001)
   - Create pool with listener_id (PL-001)

2. **Import Tests (CRITICAL)**
   - Remove pool from state
   - Import pool (PL-010)
   - Verify no drift (PL-011) **<-- GCLOUD2-20778 fix validation**

3. **Health Monitor Tests**
   - Add healthmonitor (HM-001)
   - Update healthmonitor (HM-002)
   - Delete healthmonitor (HM-003)

4. **Members Tests**
   - Add members inline (PM-001)
   - Remove one member (PM-002)
   - Remove all members (PM-003)
   - Omit members attribute (PM-004)

5. **LB Import Test**
   - Remove LB from state
   - Import LB (LB-010)
   - Verify no drift on flavor (LB-011)

6. **Resize Test**
   - Resize LB (LB-020)

7. **Cleanup**
   - Destroy all resources

---

## Success Criteria

### Critical (Must Pass)
- [ ] PL-011: Pool import with listener_id shows NO CHANGES
- [ ] LB-011: LB import shows flavor correctly (no drift)
- [ ] HM-003: Health monitor deletion uses DELETE endpoint
- [ ] PM-003: Members removal via `members = []` works

### Important (Should Pass)
- [ ] All CRUD operations complete without errors
- [ ] All imports populate state correctly
- [ ] All async operations poll correctly

### Nice to Have
- [ ] API call verification via mitmproxy logs
- [ ] Performance metrics for async operations
