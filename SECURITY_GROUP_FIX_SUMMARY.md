# Security Group Rules Fix - Summary

## Issue Found
After manual testing, you discovered that Terraform kept showing drift (wanting to delete default rules) even after multiple applies, but never actually sent the deletion request to the API.

**Root Cause:** You were using the **published provider** from the Terraform Registry (`registry.terraform.io/gcore/gcore`), which still has the old code with the `security_group_rules` field in the schema.

## Solution Implemented

### Changes Made

1. **Removed `SecurityGroupRules` field from model** (`internal/services/cloud_security_group/model.go:26`)
   - Removed the computed field that tracked backend rules in state
   - Added comment explaining this is intentional

2. **Removed `security_group_rules` from schema** (`internal/services/cloud_security_group/schema.go:187-290`)
   - Removed the entire 100+ line schema definition for the computed rules field
   - Removed unused `stringdefault` import

### Test Results

✅ **All Tests Passed:**
- Created security group without `security_group_rules` field in state
- NO drift detected across multiple `terraform apply` runs
- Backend-created default rules are completely ignored by Terraform
- State file is clean - only tracks security group metadata

### How to Use the Fixed Provider

#### Option 1: Development Override (For Testing)

1. Build the local provider:
```bash
cd /Users/user/repos/gcore-terraform
go build -o terraform-provider-gcore
```

2. Create a `.terraformrc` file in your test directory:
```hcl
provider_installation {
  dev_overrides {
    "gcore/gcore" = "/Users/user/repos/gcore-terraform"
  }

  direct {}
}
```

3. Set environment variable and run Terraform:
```bash
export TF_CLI_CONFIG_FILE="/path/to/your/.terraformrc"
terraform plan  # NO init needed with dev_overrides!
terraform apply
```

#### Option 2: Local Binary Installation (For Production Use)

1. Build and install the provider:
```bash
cd /Users/user/repos/gcore-terraform
make install
# or manually:
go build -o terraform-provider-gcore
mkdir -p ~/.terraform.d/plugins/registry.terraform.io/gcore/gcore/99.0.0/darwin_arm64/
cp terraform-provider-gcore ~/.terraform.d/plugins/registry.terraform.io/gcore/gcore/99.0.0/darwin_arm64/
```

2. Update your Terraform configuration:
```hcl
terraform {
  required_providers {
    gcore = {
      source  = "registry.terraform.io/gcore/gcore"
      version = "99.0.0"  # Use the version you installed
    }
  }
}
```

### Test Script

A comprehensive test script is available at:
`/Users/user/repos/gcore-terraform/test-secgroup-rules-only/test_drift.sh`

This script:
1. Builds the local provider
2. Creates a security group
3. Checks for drift multiple times
4. Verifies backend rules aren't being deleted
5. Confirms state doesn't track backend rules

### Example Terraform Configuration

```hcl
# Security group WITHOUT any inline rules
resource "gcore_cloud_security_group" "test" {
  project_id = 379987
  region_id  = 76

  security_group = {
    name        = "test-rules-only"
    description = "Test: rules only via separate resources"
    # NO security_group_rules field!
  }
}

# Rules managed as separate resources
resource "gcore_cloud_security_group_rule" "https" {
  group_id   = gcore_cloud_security_group.test.id
  project_id = gcore_cloud_security_group.test.project_id
  region_id  = gcore_cloud_security_group.test.region_id

  direction        = "ingress"
  ethertype        = "IPv4"
  protocol         = "tcp"
  port_range_min   = 443
  port_range_max   = 443
  remote_ip_prefix = "0.0.0.0/0"
  description      = "HTTPS"
}
```

### Before vs After

#### Before (With Bug)
```
~ security_group_rules = [
    - { id = "..." direction = "egress" ... },
    # ... 40 backend-created rules shown as drift
  ] -> (known after apply)

Plan: 0 to add, 1 to change, 0 to destroy
```
**Problem:** Shows drift but never actually deletes rules

#### After (Fixed)
```
No changes. Your infrastructure matches the configuration.
```
**Solution:** Backend rules not tracked in state, no drift shown

### Files Modified

1. `/Users/user/repos/gcore-terraform/internal/services/cloud_security_group/model.go`
   - Line 26: Removed `SecurityGroupRules` field

2. `/Users/user/repos/gcore-terraform/internal/services/cloud_security_group/schema.go`
   - Lines 187-290: Removed `security_group_rules` schema definition
   - Line 18: Removed unused `stringdefault` import

### Resources Cleaned Up

All test security groups and networks created during testing have been deleted:
- ✅ Security groups: test-auto-inherit, test-rules-only, and all test instances
- ✅ Networks: qa-terr-nw (2 instances)

### Next Steps

1. **For Development:** Use the dev_overrides method above
2. **For Production:** Work with the Stainless team to incorporate these changes into the code generation configuration
3. **Documentation:** Update provider documentation to emphasize separate rule resource pattern
4. **Migration Guide:** Create guide for users to migrate from inline rules (if any exist) to separate resources

### Important Notes

- ⚠️ After running `stainless generate`, these manual changes will be overwritten
- 📝 Keep this file as documentation of intentional changes
- 🔄 Consider creating a patch file or working with Stainless to configure generation to exclude these fields
- ✅ This pattern aligns with AWS, Azure, and OpenStack provider best practices
