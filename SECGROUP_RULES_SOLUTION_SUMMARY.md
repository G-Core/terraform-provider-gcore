# Security Group Rules - Solution Summary

## Issues Addressed

### ✅ Issue #1: Default Rules Auto-Deletion
**Status**: FIXED

The Gcore API creates 39 default egress rules when a security group is created. We implemented automatic deletion of these rules during the Create() operation so only user-defined rules remain.

**Changes**: `/Users/user/repos/gcore-terraform/internal/services/cloud_security_group/resource.go:103-161`

### ✅ Issue #2: Security Group Rule Update Error
**Status**: FIXED

When updating a rule (e.g., changing description), the API was returning error:
```
ValidationError: {'security_group_id': ['Field required']}
```

**Root Cause**: The API's PUT endpoint for updating rules actually performs a **replace operation** (delete old + create new with different ID) instead of updating in-place as documented.

**Solution**: Mark all mutable fields with `RequiresReplace` so Terraform explicitly recreates rules on any change, matching the API's actual behavior.

### ✅ Issue #3: Computed Fields Drift
**Status**: FIXED

Computed fields were showing as "(known after apply)" on every terraform plan, causing drift.

**Solution**: Added `UseStateForUnknown()` plan modifiers to all computed fields.

### ✅ Issue #4: Top-Level name Field Drift
**Status**: FIXED

The top-level `name` field was drifting and causing provider errors during Create().

**Solution**:
- Changed schema from Optional to Computed
- Changed model JSON tag from `optional` to `computed`

## Files Modified

### 1. Security Group Resource
**File**: `/Users/user/repos/gcore-terraform/internal/services/cloud_security_group/schema.go`
- Added `UseStateForUnknown()` plan modifiers to computed fields
- Changed top-level `name` to Computed

**File**: `/Users/user/repos/gcore-terraform/internal/services/cloud_security_group/model.go`
- Changed `name` JSON tag to `computed`

**File**: `/Users/user/repos/gcore-terraform/internal/services/cloud_security_group/resource.go`
- Added default rule deletion in Create() method
- Added empty list preservation in Read() method

### 2. Security Group Rule Resource
**File**: `/Users/user/repos/gcore-terraform/internal/services/cloud_security_group_rule/schema.go`
- Added `RequiresReplace` plan modifier to all mutable fields:
  - description
  - direction
  - ethertype
  - port_range_max
  - port_range_min
  - protocol
  - remote_group_id
  - remote_ip_prefix
- Changed `security_group_id` to Computed with `UseStateForUnknown()`

**File**: `/Users/user/repos/gcore-terraform/internal/services/cloud_security_group_rule/model.go`
- Changed `security_group_id` JSON tag to `computed_optional`

**File**: `/Users/user/repos/gcore-terraform/internal/services/cloud_security_group_rule/resource.go`
- Simplified Update() method to return error (should never be called)

## Current Behavior

### Security Group Creation
1. User creates security group with nested `security_group.name`
2. API creates the group + 39 default egress rules
3. Provider automatically deletes all default rules
4. Provider populates top-level computed fields (`name`, `description`, etc.)
5. State shows `security_group_rules = []` (empty, as expected)

### Security Group Rules
1. User creates rules via separate `gcore_cloud_security_group_rule` resources
2. Rules support for-each loops for bulk creation
3. **Any change to a rule attribute triggers a recreate** (destroy old, create new)
4. This matches the API's actual behavior (replace operation)

### Drift Detection
- ✅ No drift on security group computed fields
- ✅ No drift on `security_group_rules` field (stays empty)
- ✅ No drift on rule resources after creation
- ✅ State remains consistent after rule updates (recreates)

## API Issue Documentation

**File**: `/Users/user/repos/gcore-terraform/SECGROUP_RULE_PUT_ISSUE.md`

Comprehensive documentation of the API's PUT endpoint behavior, including:
- Detailed test results with curl
- Before/after ID comparison
- Impact analysis
- Test script for reproduction

**Key Finding**: The PUT endpoint returns a different rule ID, indicating the API performs a replace (delete + create) instead of an in-place update.

## Test Resources

### Test Directory
`/Users/user/repos/gcore-terraform/test-secgroup-manual/`

### Test Scripts
- `test.sh` - Basic functionality test
- `comprehensive_test.sh` - Full test suite (5 scenarios)
- `test_rule_update_api.sh` - Direct API testing with curl
- `cleanup.sh` - Quick cleanup

### Test Configuration
`main.tf` - Demonstrates:
- Security group creation with nested config
- 2 individual rule resources
- 2 rules via for-each loop
- Total: 4 user-defined rules (no defaults)

## Verification Results

### Comprehensive Test Suite
All scenarios PASSED:
- ✅ Initial create & verify (4 rules, no defaults)
- ✅ Immediate drift detection (no drift)
- ✅ Second apply is no-op (0 changes)
- ✅ No drift after second apply
- ✅ Rule add/remove operations work

### Rule Update Test
- ✅ Description change triggers `must be replaced`
- ✅ Plan shows `-/+` (destroy and create)
- ✅ Plan output: "1 to add, 0 to change, 1 to destroy"
- ✅ Apply succeeds with new rule ID
- ✅ No drift after replace operation

## User Impact

### What Users See

**Before (broken):**
```
Error: Provider produced inconsistent result after apply
When applying changes to gcore_cloud_security_group_rule.ports["http"],
provider produced an unexpected new value: .id: was
cty.StringVal("old-id"), but now cty.StringVal("new-id").
```

**After (working):**
```
Terraform will perform the following actions:

  # gcore_cloud_security_group_rule.ports["http"] must be replaced
  -/+ resource "gcore_cloud_security_group_rule" "ports" {
      ~ description = "old" -> "new" # forces replacement
      ...
    }

Plan: 1 to add, 0 to change, 1 to destroy.
```

Users now see clear messaging that rule changes require replacement, which matches the API's actual behavior.

### User Experience
- ✅ Security groups work as expected
- ✅ Default rules automatically deleted
- ✅ Rules managed via separate resources
- ✅ For-each loops supported
- ✅ No unexpected drift
- ⚠️ Rule updates are disruptive (recreate) - but this matches API reality

## Next Steps / Future Improvements

### 1. Schema Unification (Deferred)
**Request**: Flatten the nested `security_group` block so users can write:
```hcl
resource "gcore_cloud_security_group" "test" {
  name = "my-secgroup"  # Top-level, not nested
}
```

**Status**: Requires changes to `/Users/user/repos/gcore-config/openapi.yml` and Stainless Studio regeneration.

### 2. Report API Issue to Gcore
The PUT endpoint for security group rules performs a replace operation (returns new ID) instead of updating in-place as documented. This should be reported as either:
- A bug in the API implementation
- Missing documentation about the replace behavior

### 3. Add Rule Import Support
With the current replace-on-update approach, rule imports should work reliably since we're not trying to update existing rules.

## Conclusion

All critical issues are **RESOLVED**:
- ✅ Default rules automatically deleted
- ✅ No drift on any fields
- ✅ Rule updates work (via recreate)
- ✅ State stays consistent
- ✅ Comprehensive test coverage

The provider now accurately reflects the Gcore API's actual behavior, ensuring predictable and reliable resource management.
