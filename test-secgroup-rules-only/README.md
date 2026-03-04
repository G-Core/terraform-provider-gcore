# Test: Security Group Rules as Separate Resources Only

## Purpose

Test that:
1. ✅ Inline `security_group_rules` field is removed from schema (cannot be used)
2. ✅ Backend-created default rules do NOT show as nested attribute drift
3. ✅ Separate `gcore_cloud_security_group_rule` resources work correctly

## Changes Made

**Schema (`internal/services/cloud_security_group/schema.go`)**:
- Removed optional `security_group.security_group_rules` field (lines 55-134)
- Removed computed `security_group_rules` field (lines 195-298)

**Model (`internal/services/cloud_security_group/model.go`)**:
- Removed `SecurityGroupRules` field from `CloudSecurityGroupModel`
- Removed `SecurityGroupRules` field from nested `CloudSecurityGroupSecurityGroupModel`
- Removed `CloudSecurityGroupSecurityGroupSecurityGroupRulesModel` type
- Removed `CloudSecurityGroupSecurityGroupRulesModel` type

## Test Steps

### Test 1: Verify inline rules are disabled

```bash
cd /Users/user/repos/gcore-terraform/test-secgroup-rules-only

# This config would cause an error:
# cat > test-inline-error.tf <<'EOF'
# resource "gcore_cloud_security_group" "test" {
#   project_id = 379987
#   region_id  = 76
#   security_group = {
#     name = "test"
#     security_group_rules = [  # ← Should fail: attribute doesn't exist
#       {
#         direction = "ingress"
#         protocol = "tcp"
#         port_range_min = 80
#         port_range_max = 80
#       }
#     ]
#   }
# }
# EOF
#
# terraform validate
# # Expected: Error: Unsupported argument: security_group_rules
```

### Test 2: Verify separate rules work

```bash
export TF_CLI_CONFIG_FILE="/Users/user/repos/gcore-terraform/.terraformrc"
set -o allexport
source /Users/user/repos/gcore-terraform/.env
set +o allexport

# First apply - creates security group and rule
terraform init
terraform apply -auto-approve

# Should succeed:
# + resource "gcore_cloud_security_group" "test" {
#     + id = (known after apply)
#     ...
#   }
# + resource "gcore_cloud_security_group_rule" "https" {
#     + id = (known after apply)
#     ...
#   }
```

### Test 3: Verify no nested attribute drift

```bash
# After first apply, backend creates default rules
# Run plan again
terraform plan

# Expected behavior:
# ✅ NO changes shown for gcore_cloud_security_group.test
# ✅ NO nested attribute drift like:
#    ~ resource "gcore_cloud_security_group" "test" {
#      ~ security_group_rules = [ ... ]  # ← This should NOT appear!
#    }
#
# ✅ Only our explicitly managed rule exists:
#    No changes. Your infrastructure matches the configuration.
#
# Why no drift?
# - Backend creates default egress rules automatically
# - Old behavior: Showed as nested attribute change in security group
# - New behavior: Rules NOT tracked in security group state
# - Default rules just exist on backend, not tracked by Terraform
# - Only explicitly managed rule resources show in state
```

### Test 4: Verify separate rule lifecycle

```bash
# Update rule
# Change port from 443 to 8443 in main.tf
terraform plan
# Should show:
# ~ resource "gcore_cloud_security_group_rule" "https" {
#     ~ port_range_min = 443 -> 8443
#     ~ port_range_max = 443 -> 8443
#   }

terraform apply -auto-approve

# Delete rule
terraform destroy -target=gcore_cloud_security_group_rule.https -auto-approve
# Should only delete the rule, not the security group

# Cleanup
terraform destroy -auto-approve
```

## Expected Results

### ✅ Task 1: Inline rules disabled
- Users CANNOT define `security_group.security_group_rules`
- Terraform validation fails if attempted
- Error message: "Unsupported argument"

### ✅ Task 2: No nested attribute drift
- Backend-created default rules do NOT show in plan
- Security group resource shows NO changes on second plan
- Drift visualization no longer shows nested attributes

### Comparison: Before vs After

**Before (showing nested attribute drift):**
```
Terraform will perform the following actions:

  # gcore_cloud_security_group.test will be updated in-place
  ~ resource "gcore_cloud_security_group" "test" {
      ~ security_group_rules = [
          - {
              - id = "42c15bc9-7ab0-452f-96ef-b0a45c3dd27a"
              - direction = "egress"
              - protocol = "vrrp"
              ...
            },
        ]
    }

Plan: 0 to add, 1 to change, 0 to destroy.
```

**After (no drift shown):**
```
No changes. Your infrastructure matches the configuration.

Terraform has compared your real infrastructure against your configuration
and found no differences, so no changes are needed.
```

## Why This Design?

1. **Industry Best Practice**: AWS, Azure, OpenStack all recommend separate rule resources
2. **API Alignment**: Gcore API treats rules as first-class resources
3. **Modularity**: Different modules can add rules to same security group
4. **Granular Control**: Update/delete individual rules without affecting others
5. **Clean State**: Security group state doesn't track rules, avoiding drift noise
6. **Explicit Management**: Users explicitly manage only the rules they define

## References

- [SECURITY_GROUP_RULES_UX_DECISION.md](/Users/user/repos/gcore-terraform/SECURITY_GROUP_RULES_UX_DECISION.md)
- AWS documentation: Security Group Rules best practices
- Azure documentation: Network Security Rules
