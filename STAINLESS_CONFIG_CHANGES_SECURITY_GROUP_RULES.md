# Stainless Configuration Changes: Security Group Rules as Separate Resource

## Date
2025-11-10

## Objective
Configure Stainless to generate security group rules as a separate Terraform resource (`gcore_securitygroup_rule`) instead of only nested attributes within `gcore_securitygroup`.

## Problem Statement

### Original Configuration
Security group rules were only available as a subresource in the SDK/API client:
```go
client.Cloud.SecurityGroups.Rules.New(...)
client.Cloud.SecurityGroups.Rules.Replace(...)
client.Cloud.SecurityGroups.Rules.Delete(...)
```

For Terraform, this pattern alone generates only inline rules within the security group resource, NOT a separate `gcore_securitygroup_rule` resource.

### API Limitations
The Gcore API does not provide a GET endpoint for individual security group rules:
- ✅ `POST /cloud/v1/securitygroups/{project_id}/{region_id}/{group_id}/rules` - Create rule
- ✅ `PUT /cloud/v1/securitygrouprules/{project_id}/{region_id}/{rule_id}` - Update rule
- ✅ `DELETE /cloud/v1/securitygrouprules/{project_id}/{region_id}/{rule_id}` - Delete rule
- ❌ `GET /cloud/v1/securitygrouprules/{project_id}/{region_id}/{rule_id}` - **DOES NOT EXIST**

Rules can only be read via the parent security group:
- `GET /cloud/v1/securitygroups/{project_id}/{region_id}/{group_id}` returns `SecurityGroupSerializer` which includes `security_group_rules` array

## Solution

### Configuration Changes

**File:** `/Users/user/repos/gcore-config/openapi.stainless.yml`

**Location:** `resources.cloud.subresources.security_groups.subresources.rules`

**Before:**
```yaml
subresources:
  rules:
    methods:
      create: post /cloud/v1/securitygroups/{project_id}/{region_id}/{group_id}/rules
      replace: put /cloud/v1/securitygrouprules/{project_id}/{region_id}/{rule_id}
      delete: delete /cloud/v1/securitygrouprules/{project_id}/{region_id}/{rule_id}
```

**After:**
```yaml
subresources:
  rules:
    terraform:
      resource: true
      data_source: false # No data source - no GET method for individual rules
    methods:
      create: post /cloud/v1/securitygroups/{project_id}/{region_id}/{group_id}/rules
      replace:
        endpoint: put /cloud/v1/securitygrouprules/{project_id}/{region_id}/{rule_id}
        terraform:
          method: update
      delete: delete /cloud/v1/securitygrouprules/{project_id}/{region_id}/{rule_id}
```

### Key Changes Explained

1. **Added Terraform Configuration Block**
   ```yaml
   terraform:
     resource: true        # Generate Terraform resource
     data_source: false    # Skip data source (no GET endpoint)
   ```
   - `resource: true` - Instructs Stainless to generate a separate Terraform resource
   - `data_source: false` - Explicitly disables data source generation since there's no GET method

2. **Mapped `replace` to Terraform `update`**
   ```yaml
   replace:
     endpoint: put /cloud/v1/securitygrouprules/{project_id}/{region_id}/{rule_id}
     terraform:
       method: update
   ```
   - The API uses `PUT` (replace semantics) but Terraform expects `update`
   - Explicitly mapping this ensures correct CRUD behavior

3. **Pattern Match: Load Balancer Pool Members**
   This configuration follows the same pattern as `load_balancers.pools.members`:
   ```yaml
   members:
     terraform:
       resource: true
       data_source: false
     methods:
       add:
         endpoint: post /cloud/v1/lbpools/{project_id}/{region_id}/{pool_id}/member
         terraform:
           method: create
       remove: delete /cloud/v1/lbpools/{project_id}/{region_id}/{pool_id}/member/{member_id}
   ```

## How Stainless Handles Missing GET

### Read Operation Strategy

When a Terraform resource lacks a dedicated GET endpoint, Stainless generates code that:

1. **On Resource Read (Refresh)**
   - Calls the parent resource GET: `GET /cloud/v1/securitygroups/{project_id}/{region_id}/{group_id}`
   - Receives `SecurityGroupSerializer` with `security_group_rules` array
   - Finds the matching rule by `id` in the array
   - Updates Terraform state with the found rule

2. **On Resource Not Found**
   - If rule ID not in parent's `security_group_rules` array → resource deleted
   - Terraform removes from state

This is a standard pattern when child resources don't have independent GET endpoints but are returned as part of parent resource reads.

## Expected Generated Resources

### Separate Resource (Primary Pattern - Recommended)

```hcl
resource "gcore_securitygroup" "web" {
  name        = "web-sg"
  description = "Web server security group"

  project_id = 123456
  region_id  = 1
}

resource "gcore_securitygroup_rule" "http" {
  security_group_id = gcore_securitygroup.web.id

  direction        = "ingress"
  ethertype        = "IPv4"
  protocol         = "tcp"
  port_range_min   = 80
  port_range_max   = 80
  remote_ip_prefix = "0.0.0.0/0"
  description      = "Allow HTTP"

  project_id = 123456
  region_id  = 1
}
```

### Inline Rules (Secondary Pattern - Still Supported)

```hcl
resource "gcore_securitygroup" "simple" {
  name        = "simple-sg"
  project_id  = 123456
  region_id   = 1

  security_group_rules = [
    {
      direction        = "ingress"
      protocol         = "tcp"
      port_range_min   = 22
      port_range_max   = 22
      remote_ip_prefix = "10.0.0.0/8"
    }
  ]
}
```

**Note:** The inline pattern is still supported because `SecurityGroupSerializer.security_group_rules` has:
```yaml
security_group_rules:
  type: array
  x-stainless-terraform-configurability: computed_optional
```

- `computed_optional` = can be set by user OR computed by backend
- This allows inline rules to work, but they're optional

## Important: Cannot Mix Both Patterns

⚠️ **Critical Constraint:** Users CANNOT use both inline `security_group_rules` AND separate `gcore_securitygroup_rule` resources for the same security group.

This must be enforced in the provider implementation with:

1. **Conflict Detection in CustomizeDiff**
   ```go
   CustomizeDiff: func(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) error {
       if hasInlineRules(diff) && hasSeparateRules(ctx, diff, meta) {
           return fmt.Errorf(
               "cannot use inline security_group_rules with separate " +
               "gcore_securitygroup_rule resources for the same security group"
           )
       }
       return nil
   }
   ```

2. **Documentation Warning**
   Both resources must include prominent warnings about this constraint.

## Validation & Testing

### Stainless Generation Test

After regenerating the SDK and Terraform provider:

1. **Verify Resource Generation**
   ```bash
   cd /Users/user/repos/gcore-terraform
   # Check that gcore_securitygroup_rule resource exists
   ls -la internal/*securitygroup_rule*.go
   ```

2. **Check Resource Structure**
   - Resource has Create, Read, Update, Delete methods
   - Read method queries parent security group and finds matching rule
   - Schema includes all required fields

3. **Terraform Acceptance Tests**
   ```go
   func TestAccSecurityGroupRule_basic(t *testing.T) {
       resource.Test(t, resource.TestCase{
           Steps: []resource.TestStep{
               {
                   Config: testAccSecurityGroupRuleConfig_basic(),
                   Check: resource.ComposeTestCheckFunc(
                       resource.TestCheckResourceAttr("gcore_securitygroup_rule.test", "direction", "ingress"),
                       resource.TestCheckResourceAttr("gcore_securitygroup_rule.test", "protocol", "tcp"),
                   ),
               },
           },
       })
   }
   ```

### API Endpoint Verification

Confirm endpoints match configuration:

```bash
# Check OpenAPI spec has these endpoints
grep -A5 "securitygrouprules" /Users/user/repos/gcore-config/openapi.yml

# Expected:
# POST   /cloud/v1/securitygroups/{project_id}/{region_id}/{group_id}/rules
# PUT    /cloud/v1/securitygrouprules/{project_id}/{region_id}/{rule_id}
# DELETE /cloud/v1/securitygrouprules/{project_id}/{region_id}/{rule_id}
```

## Next Steps

### 1. Regenerate SDK and Provider

```bash
# In gcore-config repository
cd /Users/user/repos/gcore-config

# Trigger Stainless generation
# (Follow your team's Stainless generation process)
```

### 2. Review Generated Code

Check the generated Terraform provider code:

**Expected files:**
- `internal/provider/resource_securitygroup.go`
- `internal/provider/resource_securitygroup_rule.go` ← NEW
- `internal/provider/data_source_securitygroup.go`

**Key points to verify:**
- `resource_securitygroup_rule.go` exists
- Read method queries parent security group
- Schema has all rule fields
- Proper ID handling (rule_id from path parameter)

### 3. Implement Conflict Detection

Add `CustomizeDiff` to security group resource to prevent mixing patterns:

```go
// In resource_securitygroup.go
CustomizeDiff: customdiff.All(
    func(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) error {
        // Check if inline rules are defined
        if _, ok := diff.GetOk("security_group_rules"); !ok {
            return nil
        }

        // Check if separate rule resources exist
        // This requires checking Terraform state for gcore_securitygroup_rule
        // resources with matching security_group_id

        // If both exist, return error
        return errors.New(
            "cannot use inline security_group_rules with separate " +
            "gcore_securitygroup_rule resources for the same security group"
        )
    },
),
```

### 4. Update Documentation

**gcore_securitygroup docs:**
```markdown
## Important: Rule Management

Security group rules can be managed in two ways:

### Recommended: Separate Rule Resources

Use the `gcore_securitygroup_rule` resource for better modularity:

\`\`\`hcl
resource "gcore_securitygroup" "web" {
  name = "web-sg"
}

resource "gcore_securitygroup_rule" "http" {
  security_group_id = gcore_securitygroup.web.id
  direction         = "ingress"
  protocol          = "tcp"
  # ...
}
\`\`\`

### Alternative: Inline Rules

For simple, static security groups:

\`\`\`hcl
resource "gcore_securitygroup" "simple" {
  name = "simple-sg"

  security_group_rules = [
    { direction = "ingress", protocol = "tcp", # ... }
  ]
}
\`\`\`

⚠️ **WARNING:** You cannot mix both approaches for the same security group.
```

**gcore_securitygroup_rule docs:**
```markdown
# gcore_securitygroup_rule

Manages individual security group rules.

⚠️ **WARNING:** Do not use this resource if you have defined inline
`security_group_rules` in the `gcore_securitygroup` resource. The two
approaches cannot be mixed for the same security group.

## Example Usage

\`\`\`hcl
resource "gcore_securitygroup_rule" "http" {
  security_group_id = gcore_securitygroup.web.id

  direction        = "ingress"
  protocol         = "tcp"
  port_range_min   = 80
  port_range_max   = 80
  remote_ip_prefix = "0.0.0.0/0"
}
\`\`\`
```

### 5. Migration Guide

Create migration guide for users moving from inline to separate:

```markdown
# Migrating from Inline to Separate Security Group Rules

## Step 1: Import Existing Rules

\`\`\`bash
# Get rule IDs from state
terraform state show gcore_securitygroup.example

# Import each rule
terraform import gcore_securitygroup_rule.rule1 <rule-id-1>
terraform import gcore_securitygroup_rule.rule2 <rule-id-2>
\`\`\`

## Step 2: Update Configuration

Remove inline rules from security group:

\`\`\`diff
 resource "gcore_securitygroup" "example" {
   name = "example-sg"
-  security_group_rules = [...]
 }

+resource "gcore_securitygroup_rule" "rule1" {
+  security_group_id = gcore_securitygroup.example.id
+  # ... rule config
+}
\`\`\`

## Step 3: Plan and Apply

\`\`\`bash
terraform plan  # Should show no changes
terraform apply
\`\`\`
```

## References

- **Configuration File:** `/Users/user/repos/gcore-config/openapi.stainless.yml` (lines 461-485)
- **OpenAPI Spec:** `/Users/user/repos/gcore-config/openapi.yml` (security group rules section)
- **Stainless Docs:** https://www.stainless.com/docs/guides/getting-started-with-terraform/
- **UX Decision Doc:** `/Users/user/repos/gcore-terraform/SECURITY_GROUP_RULES_UX_DECISION.md`
- **Similar Pattern:** Load Balancer Pool Members (lines 357-366 in config)

## Benefits of This Approach

1. ✅ **Modularity** - Different modules can add rules to the same security group
2. ✅ **Granular Control** - Update/delete individual rules without affecting others
3. ✅ **Team Collaboration** - Multiple team members can manage different rules
4. ✅ **Industry Standard** - Matches AWS, Azure, and OpenStack patterns
5. ✅ **Backward Compatible** - Inline rules still work for simple use cases
6. ✅ **API Alignment** - Direct mapping to Gcore API endpoints
7. ✅ **Default Rule Handling** - Clean handling of backend-created default rules

## Risks & Mitigations

### Risk 1: Users Mix Both Patterns

**Mitigation:**
- Implement `CustomizeDiff` validation
- Clear documentation warnings
- Error messages guide users to choose one approach

### Risk 2: Read Performance

**Issue:** Reading rules requires fetching parent security group

**Mitigation:**
- This is standard practice (AWS, Azure do the same)
- Single API call per security group (not per rule)
- Rules are cached in memory after first read

### Risk 3: State Drift

**Issue:** Backend-created default rules appear as drift

**Mitigation:**
- This is intentional (AWS pattern)
- Users can import or let Terraform remove them
- Documented behavior with examples

## Success Criteria

- [x] Configuration updated in `openapi.stainless.yml`
- [ ] SDK regenerated with separate `SecurityGroupRuleService`
- [ ] Terraform provider generates `gcore_securitygroup_rule` resource
- [ ] Read operation works (queries parent, finds matching rule)
- [ ] CRUD operations map correctly (create, update, delete)
- [ ] Acceptance tests pass
- [ ] Documentation updated with warnings and examples
- [ ] Migration guide created

## Conclusion

This configuration change enables generating security group rules as separate Terraform resources while maintaining backward compatibility with inline rules. The approach follows industry best practices and aligns with the Gcore API architecture.

**Key Takeaway:** Stainless can generate Terraform resources without dedicated GET endpoints by reading through the parent resource, as demonstrated with load balancer pool members and now with security group rules.
