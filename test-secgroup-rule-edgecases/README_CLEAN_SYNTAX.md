# Clean Syntax Examples - Manual Testing Guide

## Files Overview

### 1. `manual-test-clean-syntax.tf` (RECOMMENDED FOR TESTING)
Comprehensive example with detailed comments explaining:
- ✅ How to use Terraform references (working solution)
- ❌ Why auto-fetch doesn't work (limitations explained)
- 📝 Alternative with environment variables
- 🔧 Manual test commands included

### 2. `main.tf` (UPDATED)
Your existing comprehensive edge case tests, now updated to use clean syntax:
- Security group: Uses hardcoded `project_id = 379987` and `region_id = 76`
- All 9 rules: Use references `project_id = gcore_cloud_security_group.test.project_id`

## Quick Test

```bash
# Navigate to test directory
cd /Users/user/repos/gcore-terraform/test-secgroup-rule-edgecases

# Set environment
export TF_CLI_CONFIG_FILE="/Users/user/repos/gcore-terraform/.terraformrc"
set -o allexport
source /Users/user/repos/gcore-terraform/.env
set +o allexport

# Test with manual-test-clean-syntax.tf
terraform plan -target=gcore_cloud_security_group.web \
               -target=gcore_cloud_security_group_rule.https \
               -target=gcore_cloud_security_group_rule.http \
               -target=gcore_cloud_security_group_rule.ssh

# Apply
terraform apply -auto-approve \
                -target=gcore_cloud_security_group.web \
                -target=gcore_cloud_security_group_rule.https \
                -target=gcore_cloud_security_group_rule.http \
                -target=gcore_cloud_security_group_rule.ssh

# Verify outputs
terraform output rule_details
# Should show:
# {
#   "https" = {
#     "id" = "..."
#     "project_id" = 379987  # ✅ Correctly inherited via reference
#     "region_id" = 76       # ✅ Correctly inherited via reference
#   }
#   ...
# }

# Cleanup
terraform destroy -auto-approve
```

## Key Differences

### ❌ Old Syntax (Repetitive)
```hcl
resource "gcore_cloud_security_group_rule" "https" {
  group_id   = gcore_cloud_security_group.test.id
  project_id = 379987  # Hardcoded repetition
  region_id  = 76      # Hardcoded repetition
  # ...
}
```

### ✅ New Syntax (DRY - Don't Repeat Yourself)
```hcl
resource "gcore_cloud_security_group_rule" "https" {
  group_id   = gcore_cloud_security_group.test.id
  project_id = gcore_cloud_security_group.test.project_id  # ✅ Reference parent
  region_id  = gcore_cloud_security_group.test.region_id   # ✅ Reference parent
  # ...
}
```

### 🎯 Benefits
1. **DRY Principle**: Values defined once in parent resource
2. **Maintainability**: Change project/region in one place
3. **No Runtime Overhead**: Terraform resolves references during plan phase
4. **Explicit**: Clear where values come from
5. **Type Safe**: Terraform validates references at plan time

## Current Limitations (With Solutions)

### LIMITATION 1: Nested `security_group` block

**Current (Required):**
```hcl
resource "gcore_cloud_security_group" "test" {
  security_group = {  # ← Nested block
    name        = "my-sg"
    description = "Description"
  }
}
```

**Desired (Doesn't work yet):**
```hcl
resource "gcore_cloud_security_group" "test" {
  name        = "my-sg"        # ← Flat attribute
  description = "Description"   # ← Flat attribute
}
```

**Why:** API request expects nested structure: `{"security_group": {"name": "..."}}`

**Fix:** Update OpenAPI spec to use `SingleCreateSecurityGroupSerializer` (flat) instead of `CreateSecurityGroupSerializer` (nested)

**File to change:** `/Users/user/repos/gcore-config/openapi.yml`

---

### LIMITATION 2: Must specify project_id/region_id (Can't auto-inherit)

**Why provider can't auto-fetch:**
1. Provider can't access other resources' state during CRUD operations
2. To fetch parent security group via API, we need project_id/region_id as parameters
3. Chicken-and-egg: Need values to fetch values

**Solutions:**
- ✅ **Use Terraform references** (shown in examples)
- ✅ **Use environment variables** `GCORE_CLOUD_PROJECT_ID`, `GCORE_CLOUD_REGION_ID`

---

## See Also

- **AUTO_INHERIT_IMPLEMENTATION_NOTES.md** - Detailed technical analysis
- **SECGROUP_RULES_TEST_REPORT.md** - Comprehensive test results
- **example-clean-syntax.tf** - Another example (similar to manual-test-clean-syntax.tf)
