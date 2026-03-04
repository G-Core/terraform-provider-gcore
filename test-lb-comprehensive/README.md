# Load Balancer Comprehensive Testing

This directory contains comprehensive tests for GCore Load Balancer, Listener, and Pool resources.

## Purpose

Validate:
1. **Configuration drift detection** - No false drift when state matches infrastructure
2. **Update vs Replace behavior** - PATCH used when available, not unnecessary replacements
3. **Field-level updates** - Each PATCH-able field can be updated in-place
4. **Correct replacements** - RequiresReplace fields force recreation when changed

## Test Structure

```
test-lb-comprehensive/
├── TESTING_PLAN.md          # Comprehensive test plan with all 33 test cases
├── README.md                # This file
├── run_test.sh              # Test runner script
├── .env                     # GCore credentials
├── .terraformrc             # Provider override
├── common/                  # Shared configuration
│   ├── variables.tf
│   └── provider.tf
├── drift/                   # Phase 1: Drift detection tests (4 tests)
│   └── TC-DRIFT-001-lb-no-changes/
│       └── main.tf
├── update/                  # Phase 2: Update operation tests (15 tests)
├── replace/                 # Phase 3: Replacement tests (8 tests)
├── combined/                # Phase 4: Combined resource tests (2 tests)
└── edge/                    # Phase 4: Edge case tests (4 tests)
```

## Running Tests

### Run a single test:
```bash
./run_test.sh drift/TC-DRIFT-001-lb-no-changes
```

### Run all drift tests:
```bash
for test in drift/TC-DRIFT-*; do
    ./run_test.sh "$test"
done
```

### Cleanup after tests:
```bash
cd drift/TC-DRIFT-001-lb-no-changes
terraform destroy -auto-approve
```

## Test Phases

### Phase 1: Drift Detection (Highest Priority)
- TC-DRIFT-001: Load Balancer - No changes
- TC-DRIFT-002: Listener - No changes (all optional fields)
- TC-DRIFT-003: Pool - No changes (with health monitor)
- TC-DRIFT-004: Pool - No changes (with members)

**Goal**: Verify no false drift on second apply

### Phase 2: Update Operations (High Priority)
15 tests validating PATCH operations for:
- Load Balancer: name, tags
- Listener: name, allowed_cidrs, connection_limit, timeouts
- Pool: name, algorithm, protocol, healthmonitor, members, session_persistence, timeouts

**Goal**: Verify PATCH used instead of replacement

### Phase 3: Replacement Operations (Medium Priority)
8 tests validating RequiresReplace behavior for:
- Load Balancer: flavor, vip_network_id
- Listener: protocol, protocol_port, load_balancer_id
- Pool: listener_id, load_balancer_id

**Goal**: Verify correct replacement when needed

### Phase 4: Combined & Edge Cases (Lower Priority)
6 tests for:
- Full stack (LB + Listener + Pool)
- Update cascade prevention
- Computed fields behavior
- Large member lists

**Goal**: Verify resources work together correctly

## Success Criteria

- ✅ **0% false drift** - No changes on second plan
- ✅ **100% PATCH usage** - All updateable fields use PATCH
- ✅ **Correct replacement** - RequiresReplace fields force recreation
- ✅ **No cascading updates** - Parent updates don't affect children

## Current Status

See `TESTING_PLAN.md` for detailed status of each test case.

## Test Results Summary

Will be updated as tests complete.

### Phase 1: Drift Detection
- [ ] TC-DRIFT-001 - Running...
- [ ] TC-DRIFT-002 - Pending
- [ ] TC-DRIFT-003 - Pending
- [ ] TC-DRIFT-004 - Pending

### Phase 2: Update Operations
- [ ] 0/15 tests completed

### Phase 3: Replacement Operations
- [ ] 0/8 tests completed

### Phase 4: Combined & Edge Cases
- [ ] 0/6 tests completed

## Issues Found

*(Will be updated during test execution)*

## Recommendations

*(Will be updated after test completion)*
