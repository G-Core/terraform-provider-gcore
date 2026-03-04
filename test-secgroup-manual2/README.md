# Multi-Security Group Test with Loops

This test demonstrates using loop syntax to add rules to multiple security groups.

## Test Architecture

```
┌─────────────────────────────────┐
│  Web Security Group             │
│  - HTTP (80)                    │
│  - HTTPS (443)                  │
│  - SSH (22)                     │
│  Source: 0.0.0.0/0              │
└─────────────────────────────────┘

┌─────────────────────────────────┐
│  Database Security Group        │
│  - PostgreSQL (5432)            │
│  - MySQL (3306)                 │
│  - Redis (6379)                 │
│  Source: 10.0.0.0/8             │
└─────────────────────────────────┘
```

## Key Features Tested

1. **Multiple Security Groups**: Creates 2 independent security groups
2. **Loop Syntax**: Uses `for_each` to add rules from variables
3. **Different Rule Sets**: Each security group has its own rule definitions
4. **Flexible Configuration**: Rules defined in variables for easy modification
5. **In-Place Updates**: Tests that updating one security group doesn't affect the other's rules

## Resources Created

- 2 security groups:
  - `gcore_cloud_security_group.web`
  - `gcore_cloud_security_group.database`
- 6 security group rules:
  - 3 rules for web tier (http, https, ssh)
  - 3 rules for database tier (postgresql, mysql, redis)

## Usage

### Quick Test

```bash
./test.sh
```

### Manual Testing

```bash
# 1. Initialize
terraform init

# 2. Create resources
terraform apply

# 3. Check outputs
terraform output

# 4. Verify no drift
terraform plan -detailed-exitcode

# 5. Test update (change web tier name)
# Edit main.tf: change "test-web-tier" to "test-web-tier-updated"
terraform apply

# 6. Verify database rules were NOT replaced
terraform show | grep database_rules
```

## Expected Behavior

### Initial Apply
- Creates 2 security groups
- Creates 6 rules total (3 per group)
- Default backend rules automatically deleted
- `security_group_rules` field remains empty for both groups

### After Web Tier Name Update
- Web security group updates in-place (PATCH)
- Web security group ID remains unchanged
- **Database security group NOT modified**
- **All database rules remain unchanged** (no replacement)
- No drift detected

## Verification Checklist

- [ ] Both security groups created successfully
- [ ] All 6 rules created and attached to correct groups
- [ ] No drift after initial apply
- [ ] Web tier name update works in-place
- [ ] Database tier completely unaffected by web tier changes
- [ ] No rule replacements when updating unrelated security group
- [ ] State shows `security_group_rules = []` for both groups

## Cleanup

```bash
terraform destroy -auto-approve
```

## Variables

You can customize the rules by modifying variables:

```hcl
# Add a new web rule
variable "web_rules" {
  default = {
    # ... existing rules ...
    ftp = {
      port        = 21
      protocol    = "tcp"
      description = "FTP access"
    }
  }
}

# Add a new database rule
variable "database_rules" {
  default = {
    # ... existing rules ...
    mongodb = {
      port        = 27017
      protocol    = "tcp"
      description = "MongoDB access"
      cidr        = "10.0.0.0/8"
    }
  }
}
```

## Key Testing Points

1. **Isolation**: Rules for one security group don't affect the other
2. **Scalability**: Easy to add/remove rules via variable changes
3. **Stability**: Security group updates don't trigger rule replacements
4. **Clarity**: Loop syntax makes it clear which rules belong to which group
