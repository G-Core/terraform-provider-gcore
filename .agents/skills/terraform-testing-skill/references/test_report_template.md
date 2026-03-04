# Test Report Template

## Resource: `gcore_cloud_[RESOURCE_NAME]`

**Date**: [DATE]  
**Tester**: [NAME]  
**Provider Version**: [VERSION]  
**Test Environment**: Project [PROJECT_ID], Region [REGION_ID]

## Executive Summary

- **Overall Status**: ⚠️ PASS WITH ISSUES | ✅ PASS | ❌ FAIL
- **Critical Issues Found**: [NUMBER]
- **Recommendations**: [BRIEF_SUMMARY]

## Test Coverage

### 1. Drift Testing

| Test Case | Description | Result | Notes |
|-----------|-------------|---------|-------|
| TC-DRIFT-001 | Minimal config no drift | ✅ PASS | No drift on second plan |
| TC-DRIFT-002 | With computed fields | ❌ FAIL | http_method shows drift |
| TC-DRIFT-003 | Full configuration | ✅ PASS | All fields stable |

### 2. Update Testing

| Test Case | Field | Method | Result | Notes |
|-----------|-------|---------|---------|-------|
| TC-UPDATE-001 | name | PATCH | ✅ PASS | Resource ID unchanged |
| TC-UPDATE-002 | routes | PATCH | ✅ PASS | Correctly sends empty array |
| TC-UPDATE-003 | interfaces | Special | ✅ PASS | Uses attach/detach endpoints |

### 3. Corner Cases

| Test Case | Scenario | Result | Notes |
|-----------|----------|---------|-------|
| TC-CORNER-001 | Empty arrays | ✅ PASS | Correctly clears all items |
| TC-CORNER-002 | Null fields | ✅ PASS | Maintains existing values |
| TC-CORNER-003 | Max values | ✅ PASS | Handles limits correctly |

### 4. API Call Verification

**Total API Calls Captured**: [NUMBER]

| Operation | Expected Method | Actual Method | Body Validation | Result |
|-----------|-----------------|---------------|-----------------|---------|
| Create | POST | POST | ✅ Valid | ✅ PASS |
| Update name | PATCH | PATCH | ✅ Valid | ✅ PASS |
| Update routes | PATCH | PATCH | ✅ routes=[] | ✅ PASS |
| Delete | DELETE | DELETE | N/A | ✅ PASS |

## Issues Found

### Critical Issues

#### Issue 1: Drift in computed fields
- **Severity**: HIGH
- **Description**: Fields `http_method` and `max_retries_down` show perpetual drift
- **Root Cause**: Using `apijson.Unmarshal` instead of `UnmarshalComputed`
- **Fix Applied**: Updated to use `UnmarshalComputed` in Read and ImportState
- **Status**: ✅ FIXED

### Minor Issues

#### Issue 2: [TITLE]
- **Severity**: LOW
- **Description**: [DESCRIPTION]
- **Impact**: [IMPACT]
- **Recommendation**: [FIX]

## Test Execution Details

### Environment Setup
```bash
export GCORE_API_KEY=***
export GCORE_CLOUD_PROJECT_ID=379987
export GCORE_CLOUD_REGION_ID=76
export TF_CLI_CONFIG_FILE=.terraformrc
```

### Test Commands Used
```bash
# Drift test
terraform apply -auto-approve && terraform plan -detailed-exitcode

# Update test
terraform apply -var="name_updated=true"

# API capture
mitmdump -s capture_api.py -p 8888
```

## Performance Metrics

- **Average Create Time**: [X] seconds
- **Average Update Time**: [X] seconds
- **Average Delete Time**: [X] seconds
- **Task Polling Intervals**: [X] polls over [Y] seconds

## Code Changes Required

### File: `internal/services/cloud_[resource]/resource.go`
```diff
- err = apijson.Unmarshal(bytes, &data)
+ err = apijson.UnmarshalComputed(bytes, &data)
```

### File: `internal/services/cloud_[resource]/model.go`
```diff
- HTTPMethod types.String `json:"http_method,optional"`
+ HTTPMethod types.String `json:"http_method,computed_optional"`
```

### File: `internal/services/cloud_[resource]/schema.go`
```diff
  "http_method": schema.StringAttribute{
+     Computed: true,
      Optional: true,
  }
```

## Verification After Fixes

- [ ] Rebuild provider
- [ ] Re-run all drift tests
- [ ] Re-run all update tests
- [ ] Verify API calls still correct
- [ ] No new issues introduced

## Sign-off

- **QA Engineer**: [NAME] - [DATE]
- **Developer**: [NAME] - [DATE]
- **Lead**: [NAME] - [DATE]

## Appendix

### A. Full Test Logs
- Location: `tests/[RESOURCE]/logs/`

### B. API Call Captures
- Location: `mitm_api_calls.json`

### C. Terraform Debug Logs
- Location: `terraform-debug.log`

### D. Test Configurations Used
- Minimal: `configs/[RESOURCE]/minimal.tf`
- Update: `configs/[RESOURCE]/update_*.tf`
- Corner: `configs/[RESOURCE]/corner_*.tf`