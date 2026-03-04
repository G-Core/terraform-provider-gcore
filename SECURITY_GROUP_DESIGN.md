# Security Group Resource Design Documentation

## Overview

The `gcore_cloud_security_group` resource implements an AWS-style pattern where security group rules can be managed either:
1. **Inline** - As a list within the `security_group_rules` attribute (not recommended)
2. **Separate Resources** - Using `gcore_cloud_security_group_rule` resources (recommended)

## Architecture Decisions

### 1. Optional+Computed `security_group_rules` Field

**Location:** `internal/services/cloud_security_group/model.go:29`

```go
SecurityGroupRules customfield.NestedObjectList[CloudSecurityGroupSecurityGroupRulesModel]
    `tfsdk:"security_group_rules" json:"security_group_rules,optional,computed"`
```

**Why Optional+Computed:**
- **Optional:** Users can choose to not specify this field (null), indicating they want ZERO inline rules
- **Computed:** Terraform reads the actual rules from the API to detect drift
- **No `no_refresh`:** Field must be readable during Read() operations to detect backend-created default rules

### 2. Backend Default Rules Handling

**Problem:** When a security group is created, the Gcore backend automatically adds ~39 default egress rules for various protocols (TCP, UDP, ICMP, etc. for IPv4 and IPv6).

**Solution:** AWS-style drift detection:
1. **First Apply:** Creates security group, backend adds default rules
2. **Second Apply:** Terraform detects default rules as drift and removes them via `changed_rules` API field
3. **Third Apply:** No drift, infrastructure matches configuration

### 3. The `changed_rules` API Field

**Location:** `internal/services/cloud_security_group/model.go:20`

```go
ChangedRules *[]*CloudSecurityGroupChangedRulesModel
    `tfsdk:"changed_rules" json:"changed_rules,optional,no_refresh"`
```

**API Behavior (from Gcore documentation):**
- Accepts array of rule operations: `{"action": "create"}` or `{"action": "delete"}`
- For delete: requires `security_group_rule_id` UUID
- **Selective deletion:** Rules not mentioned in `changed_rules` remain untouched
- **Critical:** If `security_group_rules: []` is sent in Update request, API deletes ALL rules

**Why `no_refresh`:**
- Write-only field for API operations
- Not tracked in Terraform state
- Only used to send delete/create operations to API

### 4. Revision Number & Updated At

**Location:** `internal/services/cloud_security_group/resource.go:188-192`

```go
// Fix for "inconsistent result" error
data.RevisionNumber = state.RevisionNumber
data.UpdatedAt = state.UpdatedAt
```

**Problem:** During Update operations (e.g., deleting rules), the API increments `revision_number` and updates `updated_at` as side-effects. This causes Terraform validation error:

```
Error: Provider produced inconsistent result after apply
.revision_number: was cty.NumberIntVal(1), but now cty.NumberIntVal(2)
```

**Solution:** Preserve prior state values for these computed fields. They will sync on next `terraform refresh`.

## Update Logic (Drift Removal)

### Location
`internal/services/cloud_security_group/resource.go:136-175`

### Algorithm

```go
// Step 1: Identify backend-created rules (those without descriptions)
if !state.SecurityGroupRules.IsNullOrUnknown() {
    stateRules := state.SecurityGroupRules.AsStructSliceT(ctx)

    for _, rule := range stateRules {
        // Only delete rules without descriptions - these are backend defaults
        // User rules from separate resources have descriptions
        if rule.Description.IsNull() || rule.Description.ValueString() == "" {
            changedRules.append(delete operation for this rule)
        }
    }
}

// Step 2: Prevent sending security_group_rules in Update request
// CRITICAL: If we send security_group_rules: [], API deletes ALL rules
origSecurityGroupRules := data.SecurityGroupRules
data.SecurityGroupRules = state.SecurityGroupRules  // Set to same value so MarshalForUpdate skips it

dataBytes := data.MarshalJSONForUpdate(*state)

data.SecurityGroupRules = origSecurityGroupRules  // Restore after marshaling
```

### Why Filter by Description

**Assumption:** Backend-created default rules never have descriptions, while user-created rules (from separate `gcore_cloud_security_group_rule` resources) always have descriptions.

**Verification from mitmproxy logs:**
- Backend rules: No `description` field in JSON
- User rules: Always include `"description": "Allow HTTPS"` or similar

## Read Logic (Prevent False Drift)

### Location
`internal/services/cloud_security_group/resource.go:254-275`

### Algorithm

```go
// Step 1: Read all rules from API
err = apijson.Unmarshal(bytes, &data)  // Unmarshals all 41 rules

// Step 2: Filter to only include backend-created rules (no description)
if !data.SecurityGroupRules.IsNullOrUnknown() {
    allRules := data.SecurityGroupRules.AsStructSliceT(ctx)

    backendRules := make([]CloudSecurityGroupSecurityGroupRulesModel, 0)
    for _, rule := range allRules {
        // Only keep rules WITHOUT descriptions (backend defaults)
        if rule.Description.IsNull() || rule.Description.ValueString() == "" {
            backendRules = append(backendRules, rule)
        }
    }

    // Update data with filtered rules
    filteredRules := customfield.NewObjectList(ctx, backendRules)
    data.SecurityGroupRules = filteredRules
}

// Step 3: Set filtered state
resp.State.Set(ctx, &data)
```

**Result:** User-created rules (with descriptions) are excluded from `security_group_rules` attribute, preventing false drift detection on subsequent applies.

## Testing Methodology

### Test Scenario: 3-Apply Test

```bash
# Apply 1: Create security group + user rules via separate resources
terraform apply -auto-approve
# Expected:
# - Security group created
# - 2 user rules created (with descriptions)
# - Backend adds 39 default rules (no descriptions)
# - Total: 41 rules in API

# Apply 2: Remove backend default rules (drift)
terraform apply -auto-approve
# Expected:
# - changed_rules sent with 39 delete operations (only rules without descriptions)
# - 2 user rules preserved
# - 39 backend rules removed
# - Total: 2 rules in API
# - No "inconsistent result" errors

# Apply 3: Verify no drift
terraform apply -auto-approve
# Expected:
# - "No changes. Your infrastructure matches the configuration."
# - No rules recreated
# - No drift detected
```

### Verification Commands

```bash
# Check API state
curl -H "Authorization: APIKey $GCORE_API_KEY" \
  "https://api.gcore.com/cloud/v1/securitygroups/$PROJECT_ID/$REGION_ID/$SG_ID" | \
  jq '.security_group_rules | length'

# Check Terraform state
terraform show | grep security_group_rules
```

### Test Results ✅ ALL PASSING

**Apply 1:**
- ✅ Security group created
- ✅ 2 user rules created (SSH, HTTPS with descriptions)
- ✅ 39 backend default rules added by API
- ✅ Total: 41 rules in API
- ✅ State: `security_group_rules = []` (correct - Create doesn't trigger Read filtering)

**Apply 2:**
- ✅ Drift detected: 39 backend rules shown as needing removal
- ✅ Update generated 39 delete operations (only rules without descriptions)
- ✅ All 39 backend rules removed from API
- ✅ 2 user rules preserved
- ✅ Final API count: 2 rules
- ✅ State: `security_group_rules = []`
- ✅ No "inconsistent result" errors

**Apply 3:**
- ✅ "No changes. Your infrastructure matches the configuration."
- ✅ No drift detected
- ✅ User rules not appearing in security_group_rules attribute
- ✅ All resources stable

### mitmproxy Monitoring

Use mitmproxy to capture API calls:

```bash
# Expected calls on Apply 2:
# 1. GET  - Read current state (shows 41 rules)
# 2. PATCH/PUT - Update with changed_rules containing 39 delete operations
# 3. GET  - Verify update (shows 2 rules remaining)
```

## Known Issues & Fixes

### Issue 1: Revision Number Inconsistency ✅ FIXED

**Error:**
```
Error: Provider produced inconsistent result after apply
.revision_number: was cty.NumberIntVal(1), but now cty.NumberIntVal(2)
```

**Fix:** Preserve `revision_number` and `updated_at` from prior state during Update operations.

### Issue 2: User Rules Deleted on Apply 2 ✅ FIXED

**Problem:** All rules (including user-created ones) were being deleted on second apply.

**Root Cause:**
1. Generated delete operations for ALL rules in state
2. Sent `security_group_rules: []` which caused API to delete everything

**Fix:**
1. Filter rules: Only delete rules WITHOUT descriptions
2. Prevent sending `security_group_rules` field by setting it to same value as state before marshaling

### Issue 3: User Rules Show as Drift on Apply 3 ✅ FIXED

**Problem:** User-created rules (from separate `gcore_cloud_security_group_rule` resources) appear in `security_group_rules` state attribute, causing drift detection on subsequent applies.

**Example:**
```hcl
# On Apply 3, Terraform wants to remove user rules:
~ security_group_rules = [
  - {
      - description = "Allow HTTPS"  # User-created rule showing as drift
      - direction   = "ingress"
      ...
    },
]
```

**Fix:** Added filtering in Read() function (lines 254-275 in resource.go) to exclude user-created rules (those with descriptions) from the `security_group_rules` attribute. This prevents false drift detection on subsequent applies.

## Implementation Files

- **Schema:** `internal/services/cloud_security_group/schema.go`
- **Model:** `internal/services/cloud_security_group/model.go`
- **Resource Logic:** `internal/services/cloud_security_group/resource.go`
- **Rule Resource:** `internal/services/cloud_security_group_rule/` (separate resource)

## Related Documentation

- [SECURITY_GROUP_RULES_UX_DECISION.md](./SECURITY_GROUP_RULES_UX_DECISION.md) - UX design decisions
- [Gcore API Documentation](https://gcore.com/docs/api-reference/cloud/security-groups/update-security-group)
