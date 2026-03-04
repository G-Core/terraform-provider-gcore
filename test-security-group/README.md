# Security Group Rules Test Examples

This directory contains examples demonstrating the security group rules functionality in the Gcore Terraform provider.

## Overview

Security group rules can be managed as **separate resources** (recommended) or as inline attributes within the security group resource.

**Important:** You cannot mix both approaches for the same security group. Choose one pattern and stick with it.

## Files

- **`main.tf`** - Original example using inline rules (legacy pattern)
- **`main-minimal.tf`** - Minimal example with one security group and one rule (recommended for quick testing)
- **`main-separate-rules.tf`** - Comprehensive examples using separate rule resources (recommended for production)

## Recommended Pattern: Separate Resources

### Why Use Separate Resources?

1. ✅ **Modularity** - Different modules can add rules to the same security group
2. ✅ **Granular Control** - Update/delete individual rules without affecting others
3. ✅ **Team Collaboration** - Multiple team members can manage different rules
4. ✅ **Dynamic Rules** - Easy to use with `for_each` and `count`
5. ✅ **Industry Standard** - Matches AWS, Azure, and OpenStack patterns

### Basic Usage

```hcl
resource "gcore_cloud_security_group" "example" {
  project_id = var.project_id
  region_id  = var.region_id

  security_group = {
    name        = "example-sg"
    description = "Example security group"
  }
}

resource "gcore_cloud_security_group_rule" "https" {
  group_id   = gcore_cloud_security_group.example.id
  project_id = var.project_id
  region_id  = var.region_id

  direction        = "ingress"
  ethertype        = "IPv4"
  protocol         = "tcp"
  port_range_min   = 443
  port_range_max   = 443
  remote_ip_prefix = "0.0.0.0/0"
  description      = "Allow HTTPS"
}
```

## AWS-Like Default Rules Behavior

The Gcore provider follows the AWS pattern for handling backend-created default rules:

### First Apply
```bash
$ terraform apply
```
Creates the security group. The backend automatically adds default rules (typically egress allow-all).

### Second Apply (Drift Detection)
```bash
$ terraform plan
```
Terraform detects the backend-created rules as drift and shows:
```
# gcore_cloud_security_group.example has changed
~ resource "gcore_cloud_security_group" "example" {
    + security_group_rules = [
        + {
            + id               = "rule-xyz"
            + direction        = "egress"
            + protocol         = "any"
            # ...
          },
      ]
  }

# Terraform will suggest removing these rules
```

### Your Options

**Option 1: Let Terraform Remove Default Rules (Secure by Default)**
```bash
$ terraform apply
```
Terraform removes the backend-created rules. Your security group now only has the rules you explicitly defined.

**Option 2: Adopt Default Rules into Terraform**

Import the default rule:
```bash
$ terraform import gcore_cloud_security_group_rule.default_egress <rule-id>
```

Add to your configuration:
```hcl
resource "gcore_cloud_security_group_rule" "default_egress" {
  group_id   = gcore_cloud_security_group.example.id
  project_id = var.project_id
  region_id  = var.region_id

  direction        = "egress"
  ethertype        = "IPv4"
  protocol         = "any"
  remote_ip_prefix = "0.0.0.0/0"
  description      = "Allow all outbound traffic"
}
```

## Testing

### Quick Test (Minimal Example)
```bash
# Use the minimal example
terraform init
terraform plan -target=module.minimal
terraform apply -target=module.minimal

# On second apply, observe drift detection
terraform plan
```

### Comprehensive Test
```bash
# Use the full example
cp main-separate-rules.tf test.tf
terraform init
terraform plan
terraform apply

# Verify rules are created correctly
terraform show
```

### Test Drift Detection
```bash
# After first apply, check for default rules
terraform plan

# You should see backend-created rules detected as drift
```

## Environment Variables

Set these before running terraform:
```bash
export GCORE_API_KEY="your-api-key"
export GCORE_CLOUD_PROJECT_ID=379987
export GCORE_CLOUD_REGION_ID=76
```

Or use the `.env` file pattern from the parent directory.

## Common Issues

### Issue: "Cannot mix inline and separate rules"
**Solution:** Choose one pattern. If using separate resources, don't define `security_group_rules` in the `security_group` block.

### Issue: "Rule not found in security group"
**Solution:** The rule was deleted outside Terraform. Run `terraform apply` to recreate it or remove it from your configuration.

### Issue: Default rules keep appearing
**Solution:** This is expected behavior. Either:
1. Let Terraform remove them (secure by default)
2. Import them and add to your configuration

## Related Documentation

- [SECURITY_GROUP_RULES_UX_DECISION.md](../SECURITY_GROUP_RULES_UX_DECISION.md) - Full analysis and decision rationale
- [STAINLESS_CONFIG_CHANGES_SECURITY_GROUP_RULES.md](../STAINLESS_CONFIG_CHANGES_SECURITY_GROUP_RULES.md) - Technical implementation details

## Ticket

GCLOUD2-20783 - Implement security group rules as separate resource
