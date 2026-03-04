# Workarounds for Nested `security_group` Block

## Problem Statement

**Current (Required):**
```hcl
resource "gcore_cloud_security_group" "test" {
  project_id = 379987
  region_id  = 76

  security_group = {  # ← Nested block (poor UX)
    name        = "my-sg"
    description = "Description"
  }
}
```

**Desired:**
```hcl
resource "gcore_cloud_security_group" "test" {
  project_id  = 379987
  region_id   = 76
  name        = "my-sg"        # ← Flat attribute
  description = "Description"   # ← Flat attribute
}
```

**Root Cause:** OpenAPI spec uses `CreateSecurityGroupSerializer` which has nested structure:
```yaml
CreateSecurityGroupSerializer:
  properties:
    security_group:  # ← Nesting layer
      $ref: '#/components/schemas/SingleCreateSecurityGroupSerializer'
```

**Long-term Solution:** Change OpenAPI spec to use `SingleCreateSecurityGroupSerializer` directly (requires coordination with API team, regeneration, testing).

---

## Workaround Options (Provider-Level)

### Option 1: Custom Schema Override with Alias Fields ✅ RECOMMENDED

**Approach:** Add flat `name` and `description` fields alongside the nested block, mark them as mutually exclusive.

**Implementation:**

#### Step 1: Modify `internal/services/cloud_security_group/model.go`

Add flat fields with conflict detection:

```go
type CloudSecurityGroupModel struct {
	ID                 types.String                                                            `tfsdk:"id" json:"id,computed"`
	ProjectID          types.Int64                                                             `tfsdk:"project_id" path:"project_id,optional"`
	RegionID           types.Int64                                                             `tfsdk:"region_id" path:"region_id,optional"`

	// WORKAROUND: Support both nested and flat syntax
	// These are mutually exclusive - user must choose one approach
	SecurityGroup      *CloudSecurityGroupSecurityGroupModel                                   `tfsdk:"security_group" json:"-"` // Don't serialize directly

	// NEW: Flat fields (preferred)
	Name               types.String                                                            `tfsdk:"name" json:"-"` // Don't serialize directly
	Description        types.String                                                            `tfsdk:"description" json:"-"` // Don't serialize directly

	// Existing fields...
	Instances          *[]types.String                                                         `tfsdk:"instances" json:"instances,optional,no_refresh"`
	Tags               *map[string]types.String                                                `tfsdk:"tags" json:"tags,optional,no_refresh"`
	// ... rest of fields
}
```

#### Step 2: Custom MarshalJSON to handle both formats

```go
func (m CloudSecurityGroupModel) MarshalJSON() (data []byte, err error) {
	// Create the API request structure
	type APIRequest struct {
		SecurityGroup *CloudSecurityGroupSecurityGroupModel `json:"security_group,omitempty"`
		Instances     *[]types.String                       `json:"instances,omitempty"`
		Tags          *map[string]types.String              `json:"tags,omitempty"`
		// ... other fields
	}

	req := APIRequest{}

	// WORKAROUND LOGIC: Support both nested and flat syntax
	if m.SecurityGroup != nil {
		// User used nested syntax (old way)
		req.SecurityGroup = m.SecurityGroup
	} else if !m.Name.IsNull() {
		// User used flat syntax (new way)
		req.SecurityGroup = &CloudSecurityGroupSecurityGroupModel{
			Name:        m.Name,
			Description: m.Description,
		}
	}

	req.Instances = m.Instances
	req.Tags = m.Tags
	// ... copy other fields

	return apijson.MarshalRoot(req)
}
```

#### Step 3: Add schema validation to prevent both being set

Modify `internal/services/cloud_security_group/resource.go`:

```go
func (r *CloudSecurityGroupResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var plan *CloudSecurityGroupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validation: Ensure user doesn't set both nested and flat fields
	hasNested := plan.SecurityGroup != nil
	hasFlat := !plan.Name.IsNull() || !plan.Description.IsNull()

	if hasNested && hasFlat {
		resp.Diagnostics.AddError(
			"Conflicting configuration",
			"Cannot use both 'security_group' block and flat 'name'/'description' attributes. "+
			"Please use only one approach:\n"+
			"Option 1 (Nested): security_group = { name = \"...\", description = \"...\" }\n"+
			"Option 2 (Flat): name = \"...\", description = \"...\"",
		)
		return
	}

	if !hasNested && !hasFlat {
		resp.Diagnostics.AddError(
			"Missing required fields",
			"Must specify either 'security_group' block or 'name' attribute",
		)
		return
	}
}
```

#### Step 4: Schema definition

The schema would need to mark fields appropriately. You might need to add custom schema overrides.

**Pros:**
- ✅ Backward compatible (existing configs still work)
- ✅ Allows gradual migration to flat syntax
- ✅ No OpenAPI changes needed
- ✅ Clear error messages for users
- ✅ Can be implemented entirely in provider code

**Cons:**
- ⚠️ Changes will be lost on next Stainless regeneration (need manual re-application)
- ⚠️ More complex code to maintain
- ⚠️ Need to document both syntaxes
- ⚠️ Eventually want to deprecate nested syntax

**Maintenance Strategy:**
1. Create a patch file with the changes
2. After each `stainless generate`, reapply the patch
3. Document in `CONTRIBUTING.md` or similar
4. Eventually deprecate nested syntax with migration guide

---

### Option 2: Custom JSON Marshaling Only (Minimal Changes)

**Approach:** Keep the nested schema but add custom JSON marshaling to also accept flat fields.

**Implementation:**

Only modify `MarshalJSON` in model.go to handle flat fields if they exist:

```go
func (m CloudSecurityGroupModel) MarshalJSON() (data []byte, err error) {
	// If name is set directly (not in nested block), construct the nested structure for API
	if m.SecurityGroup == nil && !m.Name.IsNull() {
		m.SecurityGroup = &CloudSecurityGroupSecurityGroupModel{
			Name:        m.Name,
			Description: m.Description,
		}
	}

	return apijson.MarshalRoot(m)
}
```

**Pros:**
- ✅ Minimal code changes
- ✅ Smaller patch to maintain

**Cons:**
- ❌ Still requires adding `name` and `description` fields to the model
- ❌ Won't validate in Terraform schema (fields not defined)
- ❌ Doesn't actually solve the UX problem

**Verdict:** Not sufficient on its own

---

### Option 3: Post-Generation Script (Automated Patching)

**Approach:** Create a script that automatically modifies generated code after `stainless generate`.

**Implementation:**

Create `scripts/patch-security-group-schema.sh`:

```bash
#!/bin/bash
# Automatically flatten security_group schema after Stainless generation

set -e

MODEL_FILE="internal/services/cloud_security_group/model.go"
RESOURCE_FILE="internal/services/cloud_security_group/resource.go"

echo "Patching security group schema to support flat attributes..."

# 1. Add flat fields to model
# ... sed commands to modify model.go

# 2. Add custom MarshalJSON
# ... sed commands or Go AST manipulation

# 3. Add validation to ModifyPlan
# ... sed commands

echo "✅ Security group schema patched successfully"
echo "⚠️  Remember to test: cd test-secgroup-rule-edgecases && terraform plan"
```

**Pros:**
- ✅ Can be run automatically after regeneration
- ✅ Repeatable and documented
- ✅ Can be part of CI/CD

**Cons:**
- ⚠️ Script complexity (sed/awk fragile)
- ⚠️ May break if generated code structure changes significantly
- ⚠️ Still requires maintaining patch logic

---

### Option 4: Stainless Custom Template (If Supported)

**Approach:** Check if Stainless supports custom templates or hooks for code generation.

**Investigation Needed:**
- Check Stainless documentation for template customization
- Look for `x-stainless-*` extensions in OpenAPI that might help
- Contact Stainless support for guidance

**Current Finding:**
In the OpenAPI spec, there's already a hint:
```yaml
security_group:
  $ref: '#/components/schemas/SingleCreateSecurityGroupSerializer'
  x-stainless-terraform-configurability: required
```

**Possible Research:**
- Are there other `x-stainless-terraform-*` extensions?
- Can we use `x-stainless-terraform-flatten: true`?
- Can we override the nesting behavior?

**Action:** Research Stainless documentation or contact support

---

## Recommended Immediate Action Plan

### Phase 1: Quick Win (1-2 hours)

**Option 1A: Minimal Manual Patch**

1. **Backup current files:**
   ```bash
   cp internal/services/cloud_security_group/model.go internal/services/cloud_security_group/model.go.backup
   cp internal/services/cloud_security_group/resource.go internal/services/cloud_security_group/resource.go.backup
   ```

2. **Add flat fields to model** (lines 18-20 in model.go):
   ```go
   // Add after SecurityGroup field:
   Name        types.String `tfsdk:"name" json:"-"`        // Flat alternative
   Description types.String `tfsdk:"description" json:"-"` // Flat alternative
   ```

3. **Modify MarshalJSON** to support both:
   ```go
   func (m CloudSecurityGroupModel) MarshalJSON() (data []byte, err error) {
       // Support both nested and flat syntax
       if m.SecurityGroup == nil && !m.Name.IsNull() {
           m.SecurityGroup = &CloudSecurityGroupSecurityGroupModel{
               Name:        m.Name,
               Description: m.Description,
           }
       }
       return apijson.MarshalRoot(m)
   }
   ```

4. **Add validation in ModifyPlan** (prevent both being set)

5. **Test:**
   ```bash
   cd test-secgroup-rule-edgecases

   # Create test-flat-syntax.tf with flat syntax
   terraform plan
   ```

6. **Document the patch** in `WORKAROUNDS_APPLIED.md`

**Time:** 1-2 hours
**Risk:** Low (backward compatible)
**Maintenance:** Need to reapply after Stainless regeneration

### Phase 2: Long-term Solution (coordinate with API team)

1. Update OpenAPI spec in `/Users/user/repos/gcore-config/openapi.yml`
2. Change POST endpoint to use `SingleCreateSecurityGroupSerializer`
3. Regenerate provider
4. Deprecate nested syntax
5. Provide migration guide

---

## Decision Matrix

| Option | Time | Complexity | Backward Compatible | Maintenance | Recommended |
|--------|------|------------|---------------------|-------------|-------------|
| Option 1: Alias Fields | 2h | Medium | ✅ Yes | Manual patch | ✅ **YES** |
| Option 2: Marshal Only | 1h | Low | ❌ No | Manual patch | ❌ No |
| Option 3: Auto-patch Script | 4h | High | ✅ Yes | Script updates | ⚠️ Maybe |
| Option 4: Stainless Template | ? | ? | ✅ Yes | None | ⚠️ Research |
| Long-term: Fix OpenAPI | 1w+ | Medium | ❌ Breaking | None | 🎯 **Ultimate goal** |

---

## Testing Plan

After implementing workaround:

```bash
# Test 1: Flat syntax (new)
cat > test-flat-syntax.tf <<'EOF'
resource "gcore_cloud_security_group" "test" {
  project_id  = 379987
  region_id   = 76
  name        = "test-flat-syntax"
  description = "Testing flat syntax"
}
EOF

terraform plan
terraform apply -auto-approve

# Test 2: Nested syntax (backward compatibility)
cat > test-nested-syntax.tf <<'EOF'
resource "gcore_cloud_security_group" "test" {
  project_id = 379987
  region_id  = 76
  security_group = {
    name        = "test-nested-syntax"
    description = "Testing nested syntax"
  }
}
EOF

terraform plan
terraform apply -auto-approve

# Test 3: Conflicting (should error)
cat > test-conflict.tf <<'EOF'
resource "gcore_cloud_security_group" "test" {
  project_id = 379987
  region_id  = 76
  name       = "test-conflict"  # Both set - should error
  security_group = {
    name = "test-conflict-nested"
  }
}
EOF

terraform plan  # Should show validation error
```

---

## Next Steps

**Immediate (Today):**
1. Choose Option 1 (Alias Fields) as the workaround
2. Implement the changes manually
3. Test thoroughly
4. Document the patch for future regenerations

**Short-term (This Week):**
1. Research Option 4 (Stainless customization)
2. Create a patch file or script for automation

**Long-term (Coordinate with Team):**
1. Update OpenAPI spec
2. Regenerate provider
3. Remove workaround code
4. Celebrate clean syntax! 🎉
