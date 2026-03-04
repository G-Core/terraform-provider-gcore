# Manual Testing Guide: Security Group Rules

This directory contains a manual test configuration for security group rules with automatic default rule deletion.

## Test Configuration Overview

The `main.tf` includes:

1. **Security Group** - Automatically deletes 39 backend default egress rules during creation
2. **2 Individual Rules** - SSH (port 22) and HTTPS (port 443)
3. **2 Loop Rules** - HTTP (port 80) and MySQL (port 3306) created via `for_each`
4. **Total: 4 rules** - All user-managed, zero default rules

## Quick Start

```bash
# 1. Load environment variables
set -o allexport && source ../.env && set +o allexport

# 2. Initialize
terraform init

# 3. Apply
terraform apply -auto-approve

# 4. Verify outputs
terraform output
```

## Verification Steps

### Step 1: Check State (should show empty security_group_rules)

```bash
terraform state show gcore_cloud_security_group.test | grep security_group_rules
# Expected: security_group_rules = []
```

### Step 2: Check API Rule Count

```bash
SG_ID=$(terraform output -raw security_group_id)

curl -s -H "Authorization: APIKey $GCORE_API_KEY" \
  "https://api.gcore.com/cloud/v1/securitygroups/$GCORE_CLOUD_PROJECT_ID/$GCORE_CLOUD_REGION_ID/$SG_ID" | \
  jq '.security_group_rules | length'

# Expected: 4 (NOT 43 which would be 4 + 39 defaults)
```

### Step 3: View All Rules Details

```bash
curl -s -H "Authorization: APIKey $GCORE_API_KEY" \
  "https://api.gcore.com/cloud/v1/securitygroups/$GCORE_CLOUD_PROJECT_ID/$GCORE_CLOUD_REGION_ID/$SG_ID" | \
  jq '.security_group_rules[] | {direction, protocol, port: .port_range_min, desc: .description}'
```

Expected output:
```json
{"direction":"ingress","protocol":"tcp","port":22,"desc":"SSH access"}
{"direction":"ingress","protocol":"tcp","port":443,"desc":"HTTPS traffic"}
{"direction":"ingress","protocol":"tcp","port":80,"desc":"HTTP web traffic"}
{"direction":"ingress","protocol":"tcp","port":3306,"desc":"MySQL database access"}
```

### Step 4: Test No Drift (CRITICAL TEST)

```bash
terraform plan -detailed-exitcode
echo "Exit code: $?"

# Expected: Exit code 0 (no changes)
# Should NOT show drift on security_group_rules field
```

## Advanced Testing

### Test A: Add More Rules via Loop

Modify the `ports` variable in `main.tf`:

```hcl
default = {
  http = {
    port        = 80
    protocol    = "tcp"
    description = "HTTP web traffic"
  }
  mysql = {
    port        = 3306
    protocol    = "tcp"
    description = "MySQL database access"
  }
  # Add these:
  redis = {
    port        = 6379
    protocol    = "tcp"
    description = "Redis cache"
  }
  postgres = {
    port        = 5432
    protocol    = "tcp"
    description = "PostgreSQL database"
  }
}
```

Then run:
```bash
terraform apply -auto-approve
# Should create 2 additional rules
```

### Test B: Override Loop Values

```bash
# Override the ports variable
terraform apply -auto-approve -var='ports={custom={port=8080,protocol="tcp",description="Custom port"}}'

# This will destroy existing loop rules and create only the custom one
```

### Test C: Add Rule via Gcore UI

1. Go to Gcore Cloud Console
2. Navigate to your security group
3. Add a new ingress rule (e.g., port 9000)
4. Run `terraform plan`
5. Should show the rule detected outside Terraform

### Test D: Remove a Rule

Comment out the `https` rule in `main.tf`, then:
```bash
terraform apply -auto-approve
# Should delete the HTTPS rule from API
```

## Success Criteria

✅ Security group created with 0 default rules (only 4 user rules)
✅ State shows `security_group_rules = []`
✅ API shows exactly 4 rules
✅ `terraform plan` exit code 0 (no drift)
✅ Can add/remove rules without affecting security_group resource
✅ Loop pattern works correctly

## Cleanup

```bash
terraform destroy -auto-approve
```

## Troubleshooting

**Issue**: API shows 43 rules instead of 4
- **Cause**: Default rule deletion failed
- **Check**: Look for errors in `terraform apply` output during security group creation

**Issue**: `terraform plan` shows drift on security_group_rules
- **Cause**: Drift prevention not working
- **Check**: Ensure Read() method is preserving empty list

**Issue**: Rules created via loop not appearing
- **Cause**: for_each syntax error
- **Check**: `terraform validate` and check variable definition
