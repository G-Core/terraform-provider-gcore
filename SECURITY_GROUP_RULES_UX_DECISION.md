# Security Group Rules: Separate Resource vs Nested Attribute

## Executive Summary

**Recommendation: Use BOTH approaches with separate resources as primary pattern**

After comprehensive analysis of the Gcore API, SDK architecture, and industry best practices from AWS, Azure, OpenStack, and GCP providers, the recommended approach is to provide **both** a separate `gcore_securitygroup_rule` resource AND support for inline rules within `gcore_securitygroup`, with clear guidance that **they cannot be mixed for the same security group**.

The separate resource should be the **primary documented pattern**, while inline rules should be supported for simple use cases with prominent warnings about the mixing constraint.

---

## Context

### Background
Your team has decided to follow the AWS-like approach for security group design where:
- A default security group is created by the backend
- On the next Terraform apply, users see a plan to remove it
- Users must either remove the default rules or explicitly add them to their `.tf` files

### Key Question
Should security group rules be:
1. **Separate resources** (`gcore_securitygroup_rule`)
2. **Nested attributes** (inline within `gcore_securitygroup`)
3. **Both** (with constraints)

---

## Analysis Framework

### 1. Gcore API & SDK Architecture

**Key Findings:**

✅ **API supports individual rule operations**
- `POST /cloud/v1/securitygroups/{project_id}/{region_id}/{group_id}/rules` - Create single rule
- `PUT /cloud/v1/securitygrouprules/{project_id}/{region_id}/{rule_id}` - Update single rule
- `DELETE /cloud/v1/securitygrouprules/{project_id}/{region_id}/{rule_id}` - Delete single rule

✅ **SDK has dedicated service**
```go
type SecurityGroupService struct {
    Options []option.RequestOption
    Rules   SecurityGroupRuleService  // ← Separate service!
}
```

✅ **Rules are first-class API resources**
- Each rule has its own UUID
- Rules have independent lifecycle (created_at, updated_at, revision_number)
- Rules can be managed individually without affecting the security group

**Implication:** The API architecture naturally supports separate resource management.

### 2. Security Group Update Mechanism

The security group PATCH endpoint uses a `changed_rules` array:

```json
{
  "changed_rules": [
    {"action": "create", ...rule_fields},
    {"action": "delete", "security_group_rule_id": "uuid"}
  ]
}
```

This "delta-based" update mechanism is critical:
- ❌ **NOT a full replacement** (unlike simple nested attributes)
- ✅ **Requires tracking which rules to add/remove**
- ⚠️ **Complex state management** for inline rules

**Implication:** Inline rules require complex reconciliation logic to generate the correct delta.

### 3. Industry Best Practices Analysis

#### AWS (Primary Reference Model)

**Approach:** Provides BOTH with strong recommendations

```hcl
# Inline rules (simpler, less flexible)
resource "aws_security_group" "example" {
  name = "example"

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# Separate rules (recommended)
resource "aws_security_group" "example" {
  name = "example"
}

resource "aws_security_group_rule" "http" {
  type              = "ingress"
  from_port         = 80
  to_port           = 80
  protocol          = "tcp"
  cidr_blocks       = ["0.0.0.0/0"]
  security_group_id = aws_security_group.example.id
}
```

**Critical Warning in AWS Documentation:**
> ⚠️ **NOTE on Security Groups and Security Group Rules:** At this time you cannot use a Security Group with in-line rules in conjunction with any Security Group Rule resources. Doing so will cause a conflict of rule settings and will overwrite rules.

**AWS Community Consensus:**
- Separate resources are the **modern best practice**
- Inline rules are considered "legacy" but still supported for simple cases
- The mixing constraint is universally accepted as necessary

**Trade-offs Discussed:**

| Aspect | Inline Rules | Separate Resources |
|--------|--------------|-------------------|
| Simplicity | ✅ Better for simple cases | ❌ More verbose |
| Modularity | ❌ Cannot extend from other modules | ✅ Can be added by any module |
| Visibility | ⚠️ Won't detect manual changes | ⚠️ Won't detect manual changes |
| State Management | ✅ Simpler | ❌ More state entries |
| Team Collaboration | ❌ All rules in one place | ✅ Rules can be in different files/modules |

#### Azure (Strongly Recommends Separate)

**Approach:** BOTH supported, but documentation strongly favors separate

```hcl
# Separate (recommended)
resource "azurerm_network_security_group" "example" {
  name                = "example-nsg"
  location            = "westus"
  resource_group_name = azurerm_resource_group.example.name
}

resource "azurerm_network_security_rule" "example" {
  name                        = "allow-http"
  priority                    = 100
  direction                   = "Inbound"
  access                      = "Allow"
  protocol                    = "Tcp"
  source_port_range          = "*"
  destination_port_range     = "80"
  network_security_group_name = azurerm_network_security_group.example.name
  resource_group_name        = azurerm_resource_group.example.name
}
```

**Azure Documentation:**
> ⚠️ **NOTE:** Using separate resources is recommended because inline rules make it impossible for other modules to add security group rules if needed.

**Same mixing constraint:** Cannot use both inline and separate rules for the same NSG.

#### OpenStack (Separate Resources Standard)

**Approach:** BOTH supported, separate is standard pattern

```hcl
resource "openstack_networking_secgroup_v2" "secgroup" {
  name                 = "my-secgroup"
  delete_default_rules = true  # ← Important for Terraform control
}

resource "openstack_networking_secgroup_rule_v2" "rule_ssh" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 22
  port_range_max    = 22
  security_group_id = openstack_networking_secgroup_v2.secgroup.id
}
```

**Key Insight:** The `delete_default_rules = true` flag shows OpenStack's solution to the default rules problem - same issue Gcore faces!

**Historical Issue:** GitHub issue #3601 titled "Decouple Security Group Rules from Security Groups" shows community demand for separate resources due to:
- Circular dependency problems with inline rules
- Inability to manage rules independently
- Conflicts with OpenStack's auto-created default rules

#### GCP (Different Model - Per-Rule Resources)

**Approach:** Each firewall rule is a top-level resource (no nesting)

```hcl
resource "google_compute_firewall" "allow_ssh" {
  name    = "allow-ssh"
  network = google_compute_network.vpc.name

  allow {
    protocol = "tcp"
    ports    = ["22"]
  }

  source_ranges = ["0.0.0.0/0"]
}
```

**Note:** GCP doesn't have "security groups" - each firewall rule is independent and applies to the VPC. Not directly comparable but shows that Google strongly favors discrete, separate rule management.

---

## Detailed Comparison: Separate vs Nested

### Approach 1: Separate Resources (Primary Recommendation)

#### Structure

```hcl
resource "gcore_securitygroup" "web" {
  name        = "web-sg"
  description = "Web server security group"

  project_id = 123456
  region_id  = 1
}

resource "gcore_securitygroup_rule" "http" {
  security_group_id = gcore_securitygroup.web.id

  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 80
  port_range_max    = 80
  remote_ip_prefix  = "0.0.0.0/0"
  description       = "Allow HTTP"

  project_id = 123456
  region_id  = 1
}

resource "gcore_securitygroup_rule" "https" {
  security_group_id = gcore_securitygroup.web.id

  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 443
  port_range_max    = 443
  remote_ip_prefix  = "0.0.0.0/0"
  description       = "Allow HTTPS"

  project_id = 123456
  region_id  = 1
}
```

#### Advantages ✅

1. **API Alignment**
   - Maps directly to Gcore API endpoints
   - Each rule operation = one API call
   - No complex reconciliation logic needed

2. **Modularity & Reusability**
   - Different modules can add rules to the same security group
   - Rules can be defined in separate files
   - Team collaboration: different team members can manage different rules
   ```hcl
   # Module A adds HTTP
   resource "gcore_securitygroup_rule" "from_module_a" {
     security_group_id = var.security_group_id
     # ... HTTP rule
   }

   # Module B adds database access
   resource "gcore_securitygroup_rule" "from_module_b" {
     security_group_id = var.security_group_id
     # ... DB rule
   }
   ```

3. **Granular Control**
   - Update/delete individual rules without affecting others
   - Clear resource dependencies in Terraform graph
   - Better error isolation (one rule fails, others unaffected)

4. **Count/For-Each Friendly**
   ```hcl
   locals {
     http_ports = [80, 443, 8080, 8443]
   }

   resource "gcore_securitygroup_rule" "web_ports" {
     for_each = toset(local.http_ports)

     security_group_id = gcore_securitygroup.web.id
     port_range_min    = each.value
     port_range_max    = each.value
     # ...
   }
   ```

5. **Handles Default Rules Cleanly**
   - Users see each backend-created rule as drift
   - Can selectively import specific rules: `terraform import gcore_securitygroup_rule.keep <rule-id>`
   - Or let Terraform remove them (AWS approach)

6. **State Clarity**
   - Each rule = one state entry
   - Easy to `terraform state rm` a single rule
   - Clear import paths

7. **Industry Standard**
   - AWS, Azure, OpenStack all prefer this approach
   - Users familiar with other clouds expect this pattern

#### Disadvantages ❌

1. **Verbosity**
   - More lines of code for simple security groups
   - Multiple resource blocks instead of one

2. **State Size**
   - More entries in Terraform state file
   - More API calls during refresh

3. **Cannot Mix with Inline**
   - If you start with separate resources, you're committed
   - Must be clearly documented

4. **Initial Learning Curve**
   - New users might find it less intuitive than inline

### Approach 2: Nested Attributes (Inline Rules)

#### Structure

```hcl
resource "gcore_securitygroup" "web" {
  name        = "web-sg"
  description = "Web server security group"

  project_id = 123456
  region_id  = 1

  rules = [
    {
      direction        = "ingress"
      ethertype        = "IPv4"
      protocol         = "tcp"
      port_range_min   = 80
      port_range_max   = 80
      remote_ip_prefix = "0.0.0.0/0"
      description      = "Allow HTTP"
    },
    {
      direction        = "ingress"
      ethertype        = "IPv4"
      protocol         = "tcp"
      port_range_min   = 443
      port_range_max   = 443
      remote_ip_prefix = "0.0.0.0/0"
      description      = "Allow HTTPS"
    }
  ]
}
```

#### Advantages ✅

1. **Simplicity for Basic Cases**
   - One resource block for entire security group
   - Good for simple, static security groups
   - Less verbose for small number of rules

2. **Atomic Operations**
   - All rules updated in one operation
   - Single state entry
   - Easier to see all rules at once

3. **Familiar Pattern**
   - Matches common Terraform patterns (like aws_instance security_groups)
   - Intuitive for beginners

#### Disadvantages ❌

1. **Complex State Reconciliation**
   - Must calculate diff between desired state and current state
   - Generate `changed_rules` array with create/delete actions
   - Risk of bugs in diff logic

2. **API Mismatch**
   - API updates rules individually, but Terraform needs to batch them
   - The `changed_rules` PATCH approach requires:
     ```go
     // Pseudo-code for complexity
     old_rule_ids := getCurrentRuleIDs()
     new_rule_specs := getDesiredRules()

     changed_rules := []
     for old_id := range old_rule_ids {
         if !matches(new_rule_specs, old_id) {
             changed_rules.append({action: "delete", id: old_id})
         }
     }
     for new_spec := range new_rule_specs {
         if !exists(old_rule_ids, new_spec) {
             changed_rules.append({action: "create", ...new_spec})
         }
     }
     ```

3. **No Modularity**
   - Cannot add rules from different modules
   - All rules must be in the same resource block
   - Breaks module composition patterns

4. **All-or-Nothing Updates**
   - Changing one rule triggers update of entire rule list
   - Harder to isolate errors
   - More complex plan outputs

5. **Default Rules Problem**
   - Backend-created default rules appear in state
   - Must track which rules are "ours" vs "theirs"
   - Complicated ignore logic or filtering needed

6. **ForEach Limitations**
   - Cannot easily iterate over rules from external sources
   - Complex dynamic block syntax required

7. **No Partial Management**
   - Cannot let Terraform manage some rules but not others
   - Cannot import individual rules

### Approach 3: Both (Recommended for Gcore)

**Recommendation:** Provide both patterns with clear constraints and guidance.

#### Primary Pattern: Separate Resources

**Default/Primary Documentation:**

```hcl
# Recommended approach
resource "gcore_securitygroup" "web" {
  name        = "web-sg"
  project_id  = 123456
  region_id   = 1
}

resource "gcore_securitygroup_rule" "http" {
  security_group_id = gcore_securitygroup.web.id
  direction         = "ingress"
  protocol          = "tcp"
  port_range_min    = 80
  port_range_max    = 80
  remote_ip_prefix  = "0.0.0.0/0"

  project_id = 123456
  region_id  = 1
}
```

#### Secondary Pattern: Inline Rules

**For simple use cases only:**

```hcl
# Alternative for simple, static security groups
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

#### Critical Implementation Requirements

1. **Mutually Exclusive Detection**
   ```go
   // In the security group resource Update/Read logic
   if len(plan.SecurityGroupRules) > 0 {
       // Check if any separate gcore_securitygroup_rule resources reference this SG
       // This requires provider-level coordination
       resp.Diagnostics.AddError(
           "Conflicting Rule Management",
           "This security group has inline rules defined. You cannot use "+
           "gcore_securitygroup_rule resources for this security group. "+
           "Choose one approach or the other.",
       )
   }
   ```

2. **Clear Documentation**

   In `gcore_securitygroup` docs:
   > ⚠️ **IMPORTANT:** You cannot use `security_group_rules` (inline rules) in conjunction with separate `gcore_securitygroup_rule` resources for the same security group. Using both will cause rule conflicts and unpredictable behavior. Choose one approach:
   >
   > - **Recommended:** Use separate `gcore_securitygroup_rule` resources for better modularity and granular control
   > - **Alternative:** Use inline `security_group_rules` for simple, static security groups

   In `gcore_securitygroup_rule` docs:
   > ⚠️ **IMPORTANT:** Do not use this resource if you have defined inline `security_group_rules` in the `gcore_securitygroup` resource. The two approaches cannot be mixed for the same security group.

3. **State Management Strategy**
   - When using inline rules: Store rules as computed attribute in state
   - When using separate resources: Each rule is independent state entry
   - Read operation must detect which mode is active

---

## Handling Default Rules (AWS Pattern)

Regardless of inline vs separate, the default rule handling follows AWS pattern:

### Backend Behavior
When creating a security group, Gcore backend automatically adds default rules (likely egress allow-all).

### Terraform Behavior (AWS Pattern)

#### On First Apply
```
Terraform will perform the following actions:

  # gcore_securitygroup.test will be created
  + resource "gcore_securitygroup" "test" {
      + id   = (known after apply)
      + name = "test-sg"
    }

Plan: 1 to add, 0 to change, 0 to destroy.
```

#### On Second Apply (Terraform Detects Drift)
```
Terraform detected the following changes made outside of Terraform since the last "terraform apply":

  # gcore_securitygroup.test has changed
  ~ resource "gcore_securitygroup" "test" {
        id   = "abc-123"
        name = "test-sg"
      + security_group_rules = [
          + {
              + id               = "rule-xyz"
              + direction        = "egress"
              + protocol         = "any"
              + ethertype        = "IPv4"
              + remote_ip_prefix = "0.0.0.0/0"
              # ...
            },
        ]
    }

Unless you have made equivalent changes to your configuration, or ignored the relevant attributes,
the following plan may include actions to undo or respond to these changes.

Terraform will perform the following actions:

  # gcore_securitygroup.test will be updated in-place
  ~ resource "gcore_securitygroup" "test" {
        id   = "abc-123"
        name = "test-sg"
      - security_group_rules = [
          - {
              - id               = "rule-xyz"
              - direction        = "egress"
              - protocol         = "any"
              # ...
            },
        ]
    }

Plan: 0 to add, 1 to change, 0 to destroy.
```

#### User Options

**Option 1: Let Terraform Remove Them (AWS Pattern)**
- User applies, Terraform removes default rules
- Security group now has no rules (locked down)
- User must explicitly add rules they need

**Option 2: Adopt Them into Terraform**

With separate resources:
```hcl
# Import the default rule
terraform import gcore_securitygroup_rule.default_egress <rule-id>

# Add to config
resource "gcore_securitygroup_rule" "default_egress" {
  security_group_id = gcore_securitygroup.test.id
  direction         = "egress"
  protocol          = "any"
  ethertype         = "IPv4"
  remote_ip_prefix  = "0.0.0.0/0"

  lifecycle {
    # Prevent accidental deletion
    prevent_destroy = true
  }
}
```

With inline rules:
```hcl
resource "gcore_securitygroup" "test" {
  name = "test-sg"

  security_group_rules = [
    {
      direction        = "egress"
      protocol         = "any"
      ethertype        = "IPv4"
      remote_ip_prefix = "0.0.0.0/0"
    }
  ]
}
```

---

## Implementation Recommendations

### Phase 1: MVP - Separate Resources Only

**Rationale:** Start with the cleaner, more maintainable approach

```go
// resource_gcore_securitygroup.go
func resourceSecurityGroup() *schema.Resource {
    return &schema.Resource{
        Schema: map[string]*schema.Schema{
            "name": {Type: schema.TypeString, Required: true},
            "description": {Type: schema.TypeString, Optional: true},
            // NO "security_group_rules" attribute

            // Computed attribute shows current rules (read-only)
            "current_rules": {
                Type:     schema.TypeList,
                Computed: true,
                Elem: &schema.Resource{
                    Schema: ruleSchema(),
                },
            },
        },
        // CRUD operations
    }
}

// resource_gcore_securitygroup_rule.go
func resourceSecurityGroupRule() *schema.Resource {
    return &schema.Resource{
        Schema: map[string]*schema.Schema{
            "security_group_id": {Type: schema.TypeString, Required: true, ForceNew: true},
            "direction":         {Type: schema.TypeString, Required: true},
            "protocol":          {Type: schema.TypeString, Optional: true},
            // ... other rule fields
        },
        // CRUD operations - direct API calls
    }
}
```

**Benefits:**
- Simpler implementation
- Better testing
- Clear user expectations
- Matches AWS/Azure best practices

### Phase 2: Add Inline Support (Optional)

If user demand exists, add inline rules with safeguards:

```go
func resourceSecurityGroup() *schema.Resource {
    return &schema.Resource{
        Schema: map[string]*schema.Schema{
            "security_group_rules": {
                Type:     schema.TypeList,
                Optional: true,
                Elem: &schema.Resource{
                    Schema: ruleSchema(),
                },
                ConflictsWith: []string{"use_separate_rules"},
            },
            "use_separate_rules": {
                Type:     schema.TypeBool,
                Optional: true,
                Default:  true,
                ConflictsWith: []string{"security_group_rules"},
            },
        },

        CustomizeDiff: func(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) error {
            // Detect if user is mixing approaches
            if hasInlineRules(diff) && hasSeparateRules(ctx, diff, meta) {
                return fmt.Errorf("cannot use inline security_group_rules with separate gcore_securitygroup_rule resources")
            }
            return nil
        },
    }
}
```

### Phase 3: Migration Tools

Provide migration path from inline to separate:

```bash
# Terraform command to split inline rules
terraform state show gcore_securitygroup.example

# User splits into separate resources, then:
terraform state rm gcore_securitygroup.example.security_group_rules[0]
terraform import gcore_securitygroup_rule.http <rule-id>
```

Or provide a helper command:
```bash
gcore-terraform migrate-rules --from inline --to separate --security-group sg-123
```

---

## Comparative Decision Matrix

| Criterion | Separate Resources | Inline Rules | Both |
|-----------|-------------------|--------------|------|
| **API Alignment** | ⭐⭐⭐⭐⭐ Perfect | ⭐⭐⭐ Needs reconciliation | ⭐⭐⭐⭐ |
| **Modularity** | ⭐⭐⭐⭐⭐ Excellent | ⭐ Limited | ⭐⭐⭐⭐ |
| **Simplicity (small scale)** | ⭐⭐⭐ More verbose | ⭐⭐⭐⭐⭐ Very simple | ⭐⭐⭐⭐ |
| **Simplicity (large scale)** | ⭐⭐⭐⭐ Good patterns | ⭐⭐ Complex blocks | ⭐⭐⭐⭐ |
| **State Management** | ⭐⭐⭐⭐ Clear | ⭐⭐⭐ Complex diff | ⭐⭐⭐ |
| **Default Rule Handling** | ⭐⭐⭐⭐⭐ Clean | ⭐⭐⭐ Needs filtering | ⭐⭐⭐⭐ |
| **Team Collaboration** | ⭐⭐⭐⭐⭐ Best | ⭐⭐ Conflicts | ⭐⭐⭐⭐⭐ |
| **Testing** | ⭐⭐⭐⭐⭐ Granular | ⭐⭐⭐ Monolithic | ⭐⭐⭐ |
| **Implementation Cost** | ⭐⭐⭐⭐ Straightforward | ⭐⭐⭐ Complex | ⭐⭐ More work |
| **Maintenance Cost** | ⭐⭐⭐⭐ Low | ⭐⭐⭐ Medium | ⭐⭐ Higher |
| **Industry Practice** | ⭐⭐⭐⭐⭐ Standard | ⭐⭐⭐ Legacy | ⭐⭐⭐⭐ AWS/Azure |

**Legend:** ⭐ = Poor, ⭐⭐⭐ = Acceptable, ⭐⭐⭐⭐⭐ = Excellent

---

## Final Recommendation

### Primary Approach: Separate Resources

**Implement `gcore_securitygroup_rule` as the primary pattern**

**Reasons:**
1. ✅ **API Architecture Alignment** - Gcore API treats rules as independent resources with dedicated endpoints
2. ✅ **Industry Standard** - AWS, Azure, OpenStack all recommend this pattern
3. ✅ **Default Rules Handling** - Cleanest way to handle backend-created default rules (users can selectively import or let Terraform remove)
4. ✅ **Modularity** - Enables module composition and team collaboration
5. ✅ **SDK Structure** - SDK already has separate `SecurityGroupRuleService`
6. ✅ **Simpler Implementation** - Direct API calls, no complex reconciliation
7. ✅ **Better Testing** - Granular resource lifecycle testing

### Optional Secondary: Inline Rules

**Consider adding inline `security_group_rules` attribute if:**
- User feedback indicates strong demand
- You have capacity for additional complexity
- You can implement robust conflict detection

**If implemented:**
- Mark as "alternative approach" in docs
- Prominently warn about mixing constraint
- Recommend separate resources for production use
- Implement the `ConflictsWith` or `CustomizeDiff` safeguards

### Do NOT Implement: Inline-Only

**Do not make inline rules the only option** because:
- ❌ Goes against Gcore API architecture
- ❌ Conflicts with industry best practices
- ❌ Creates modularity problems
- ❌ Makes default rule handling complex
- ❌ Limits user flexibility

---

## Documentation Examples

### Example 1: Simple Web Server (Recommended Pattern)

```hcl
terraform {
  required_providers {
    gcore = {
      source  = "G-Core/gcore"
      version = "~> 1.0"
    }
  }
}

resource "gcore_securitygroup" "web" {
  name        = "web-server-sg"
  description = "Security group for web servers"

  project_id = var.project_id
  region_id  = var.region_id
}

# Allow HTTP
resource "gcore_securitygroup_rule" "http" {
  security_group_id = gcore_securitygroup.web.id

  direction        = "ingress"
  ethertype        = "IPv4"
  protocol         = "tcp"
  port_range_min   = 80
  port_range_max   = 80
  remote_ip_prefix = "0.0.0.0/0"
  description      = "Allow HTTP from anywhere"

  project_id = var.project_id
  region_id  = var.region_id
}

# Allow HTTPS
resource "gcore_securitygroup_rule" "https" {
  security_group_id = gcore_securitygroup.web.id

  direction        = "ingress"
  ethertype        = "IPv4"
  protocol         = "tcp"
  port_range_min   = 443
  port_range_max   = 443
  remote_ip_prefix = "0.0.0.0/0"
  description      = "Allow HTTPS from anywhere"

  project_id = var.project_id
  region_id  = var.region_id
}

# Allow SSH from management network only
resource "gcore_securitygroup_rule" "ssh" {
  security_group_id = gcore_securitygroup.web.id

  direction        = "ingress"
  ethertype        = "IPv4"
  protocol         = "tcp"
  port_range_min   = 22
  port_range_max   = 22
  remote_ip_prefix = "10.0.0.0/24"
  description      = "Allow SSH from management network"

  project_id = var.project_id
  region_id  = var.region_id
}

# Allow egress (explicitly manage default rule)
resource "gcore_securitygroup_rule" "egress_all" {
  security_group_id = gcore_securitygroup.web.id

  direction        = "egress"
  ethertype        = "IPv4"
  protocol         = "any"
  remote_ip_prefix = "0.0.0.0/0"
  description      = "Allow all outbound traffic"

  project_id = var.project_id
  region_id  = var.region_id
}
```

### Example 2: Dynamic Rules with for_each

```hcl
locals {
  web_ports = {
    http    = { port = 80,   description = "HTTP" }
    https   = { port = 443,  description = "HTTPS" }
    http_alt = { port = 8080, description = "Alternative HTTP" }
  }
}

resource "gcore_securitygroup" "web" {
  name       = "dynamic-web-sg"
  project_id = var.project_id
  region_id  = var.region_id
}

resource "gcore_securitygroup_rule" "web_ports" {
  for_each = local.web_ports

  security_group_id = gcore_securitygroup.web.id

  direction        = "ingress"
  ethertype        = "IPv4"
  protocol         = "tcp"
  port_range_min   = each.value.port
  port_range_max   = each.value.port
  remote_ip_prefix = "0.0.0.0/0"
  description      = "Allow ${each.value.description}"

  project_id = var.project_id
  region_id  = var.region_id
}
```

### Example 3: Module Composition

```hcl
# modules/security-group/main.tf
resource "gcore_securitygroup" "this" {
  name        = var.name
  description = var.description
  project_id  = var.project_id
  region_id   = var.region_id
}

output "security_group_id" {
  value = gcore_securitygroup.this.id
}

# modules/ssh-rule/main.tf
resource "gcore_securitygroup_rule" "ssh" {
  security_group_id = var.security_group_id

  direction        = "ingress"
  protocol         = "tcp"
  port_range_min   = 22
  port_range_max   = 22
  remote_ip_prefix = var.allowed_cidr

  project_id = var.project_id
  region_id  = var.region_id
}

# root/main.tf
module "web_sg" {
  source = "./modules/security-group"
  name   = "web-sg"
  # ...
}

# Team A adds SSH access
module "ssh_access" {
  source            = "./modules/ssh-rule"
  security_group_id = module.web_sg.security_group_id
  allowed_cidr      = "10.0.1.0/24"
}

# Team B adds web rules (in different file/module)
module "web_access" {
  source            = "./modules/web-rules"
  security_group_id = module.web_sg.security_group_id
}
```

### Example 4: Handling Default Rules

```hcl
# After first apply, Terraform detects backend-created default rule
# Option 1: Let Terraform remove it (recommended for secure-by-default)

# Option 2: Adopt the default rule into Terraform
resource "gcore_securitygroup" "test" {
  name       = "test-sg"
  project_id = var.project_id
  region_id  = var.region_id
}

# Import command:
# terraform import gcore_securitygroup_rule.default_egress <rule-id-from-drift-message>

resource "gcore_securitygroup_rule" "default_egress" {
  security_group_id = gcore_securitygroup.test.id

  direction        = "egress"
  ethertype        = "IPv4"
  protocol         = "any"
  remote_ip_prefix = "0.0.0.0/0"
  description      = "Allow all outbound (backend default)"

  project_id = var.project_id
  region_id  = var.region_id
}
```

---

## Testing Strategy

### Unit Tests (Provider Code)

```go
func TestAccSecurityGroupRule_basic(t *testing.T) {
    resource.Test(t, resource.TestCase{
        PreCheck:                 func() { testAccPreCheck(t) },
        ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
        Steps: []resource.TestStep{
            {
                Config: testAccSecurityGroupRuleConfig_basic(),
                Check: resource.ComposeTestCheckFunc(
                    resource.TestCheckResourceAttr("gcore_securitygroup_rule.test", "direction", "ingress"),
                    resource.TestCheckResourceAttr("gcore_securitygroup_rule.test", "protocol", "tcp"),
                    resource.TestCheckResourceAttr("gcore_securitygroup_rule.test", "port_range_min", "80"),
                ),
            },
        },
    })
}

func TestAccSecurityGroupRule_update(t *testing.T) {
    // Test updating a single rule doesn't affect others
}

func TestAccSecurityGroupRule_delete(t *testing.T) {
    // Test deleting a rule doesn't delete the security group
}
```

### Integration Tests

```hcl
# Test 1: Multiple independent rules
# Test 2: Rules with dependencies
# Test 3: Default rule handling
# Test 4: Module composition
# Test 5: for_each patterns
```

---

## Migration Path (Future Consideration)

If you currently have inline rules and want to migrate to separate resources:

### Automated Migration Tool

```bash
#!/bin/bash
# migrate-sg-rules.sh

SG_ID=$1
STATE_FILE=${2:-terraform.tfstate}

echo "Migrating security group $SG_ID to separate rule resources..."

# Extract rule count
RULE_COUNT=$(terraform state show "gcore_securitygroup.$SG_ID" | grep -c "security_group_rules")

echo "Found $RULE_COUNT inline rules"

# For each rule, generate import commands and resource blocks
for i in $(seq 0 $((RULE_COUNT - 1))); do
    RULE_ID=$(terraform state show "gcore_securitygroup.$SG_ID" | grep "id =" | sed -n "$((i+1))p" | awk '{print $3}')

    echo "terraform import gcore_securitygroup_rule.rule_$i $RULE_ID"

    # Generate resource block (user must fill in attributes)
    cat <<EOF >> generated_rules.tf

resource "gcore_securitygroup_rule" "rule_$i" {
  security_group_id = gcore_securitygroup.$SG_ID.id
  # TODO: Fill in rule attributes from state
}
EOF
done

echo "Migration commands written. Review generated_rules.tf and run the import commands."
```

---

## Conclusion

**Implement separate `gcore_securitygroup_rule` resources as the primary pattern.**

This decision is based on:

1. ✅ **Alignment with Gcore API architecture** - Rules are first-class resources
2. ✅ **Industry best practices** - AWS, Azure, OpenStack recommend separate resources
3. ✅ **Modularity and flexibility** - Enables advanced Terraform patterns
4. ✅ **Simpler implementation** - Direct API mapping, less complex state logic
5. ✅ **Better default rule handling** - Users can selectively import or remove backend rules
6. ✅ **Team collaboration** - Multiple modules can contribute rules
7. ✅ **Testing and maintenance** - Granular resource testing, clearer state management

Inline rules can be added later as a secondary option if strong user demand emerges, but should always include:
- Clear warnings about mixing constraints
- Recommendation to use separate resources for production
- Robust conflict detection

The AWS model of supporting both approaches (with clear guidance toward separate resources) has proven successful and should be followed.

---

## References

- Gcore API Documentation: https://gcore.com/docs/api-reference/cloud/security-groups/
- Gcore SDK: `./sdk-gcore-go/cloud/securitygroup*.go`
- AWS Security Groups: https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/security_group
- Azure NSG: https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/network_security_rule
- OpenStack Security Groups: https://registry.terraform.io/providers/terraform-provider-openstack/openstack/latest/docs/resources/networking_secgroup_rule_v2
- Default Rules Research: `./secgroup_research.txt`
- Terraform Framework Documentation: Resource vs Nested Attributes
