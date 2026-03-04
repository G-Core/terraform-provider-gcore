# How to Apply Flat Syntax Workaround (Manual Steps)

## Overview

This workaround allows users to write **flat syntax** for security groups:

```hcl
# ✅ New syntax (flat)
resource "gcore_cloud_security_group" "test" {
  project_id  = 379987
  region_id   = 76
  name        = "my-sg"
  description = "Description"
}
```

Instead of the current **nested syntax**:

```hcl
# ⚠️ Legacy syntax (nested)
resource "gcore_cloud_security_group" "test" {
  project_id = 379987
  region_id  = 76
  security_group = {
    name        = "my-sg"
    description = "Description"
  }
}
```

Both syntaxes will work and are mutually exclusive (validation prevents using both).

---

## Files to Modify

1. `internal/services/cloud_security_group/model.go`
2. `internal/services/cloud_security_group/resource.go`

---

## Step 1: Backup Files

```bash
cd /Users/user/repos/gcore-terraform

cp internal/services/cloud_security_group/model.go \
   internal/services/cloud_security_group/model.go.backup

cp internal/services/cloud_security_group/resource.go \
   internal/services/cloud_security_group/resource.go.backup
```

---

## Step 2: Modify `model.go`

### Change 2.1: Update `Name` field tag (line ~19)

**Find:**
```go
Name               types.String                                                            `tfsdk:"name" json:"name,optional"`
```

**Replace with:**
```go
Name               types.String                                                            `tfsdk:"name" json:"-"` // WORKAROUND: Flat alternative to security_group.name
```

**Why:** Change `json:"name,optional"` to `json:"-"` to prevent direct JSON serialization (we'll handle it manually).

---

### Change 2.2: Update `Description` field tag (line ~23)

**Find:**
```go
Description        types.String                                                            `tfsdk:"description" json:"description,computed"`
```

**Replace with:**
```go
Description        types.String                                                            `tfsdk:"description" json:"-"` // WORKAROUND: Flat alternative to security_group.description
```

**Why:** Same reason - we'll manually construct the nested structure.

---

### Change 2.3: Replace `MarshalJSON` method (line ~31)

**Find:**
```go
func (m CloudSecurityGroupModel) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(m)
}
```

**Replace with:**
```go
func (m CloudSecurityGroupModel) MarshalJSON() (data []byte, err error) {
	// WORKAROUND: Support both nested and flat syntax
	//
	// Flat syntax:  name = "...", description = "..."
	// Nested syntax: security_group = { name = "...", description = "..." }
	//
	// API always expects nested structure, so convert flat to nested here.

	model := m

	// If user provided flat fields, construct nested structure for API
	if model.SecurityGroup == nil && !model.Name.IsNull() {
		model.SecurityGroup = &CloudSecurityGroupSecurityGroupModel{
			Name:        model.Name,
			Description: model.Description,
		}
	}

	return apijson.MarshalRoot(model)
}
```

---

### Change 2.4: Replace `MarshalJSONForUpdate` method (line ~35)

**Find:**
```go
func (m CloudSecurityGroupModel) MarshalJSONForUpdate(state CloudSecurityGroupModel) (data []byte, err error) {
	return apijson.MarshalForPatch(m, state)
}
```

**Replace with:**
```go
func (m CloudSecurityGroupModel) MarshalJSONForUpdate(state CloudSecurityGroupModel) (data []byte, err error) {
	// WORKAROUND: Same logic for updates
	model := m

	// Convert flat to nested if needed
	if model.SecurityGroup == nil && !model.Name.IsNull() {
		model.SecurityGroup = &CloudSecurityGroupSecurityGroupModel{
			Name:        model.Name,
			Description: model.Description,
		}
	}

	return apijson.MarshalForPatch(model, state)
}
```

---

## Step 3: Modify `resource.go`

### Change 3.1: Replace `ModifyPlan` method (line ~286)

**Find:**
```go
func (r *CloudSecurityGroupResource) ModifyPlan(_ context.Context, _ resource.ModifyPlanRequest, _ *resource.ModifyPlanResponse) {

}
```

**Replace with:**
```go
func (r *CloudSecurityGroupResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var plan *CloudSecurityGroupModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() || plan == nil {
		return
	}

	// WORKAROUND: Validate mutually exclusive syntax options
	//
	// User can choose ONE of these syntaxes:
	//   1. Nested: security_group = { name = "...", description = "..." }
	//   2. Flat:   name = "...", description = "..."

	hasNested := plan.SecurityGroup != nil && !plan.SecurityGroup.Name.IsNull()
	hasFlat := !plan.Name.IsNull()

	if hasNested && hasFlat {
		resp.Diagnostics.AddError(
			"Conflicting security group configuration",
			"Cannot use both 'security_group' block and flat 'name'/'description' attributes.\n\n"+
			"Choose ONE syntax:\n\n"+
			"Option 1 (Nested - Legacy):\n"+
			"  security_group = {\n"+
			"    name        = \"my-sg\"\n"+
			"    description = \"Description\"\n"+
			"  }\n\n"+
			"Option 2 (Flat - Recommended):\n"+
			"  name        = \"my-sg\"\n"+
			"  description = \"Description\"",
		)
		return
	}

	if !hasNested && !hasFlat {
		resp.Diagnostics.AddError(
			"Missing required security group name",
			"Must specify security group name using ONE of these syntaxes:\n\n"+
			"  security_group = { name = \"my-sg\" }\n\n"+
			"OR\n\n"+
			"  name = \"my-sg\"",
		)
		return
	}

	// project_id and region_id will use provider-level defaults automatically
}
```

---

## Step 4: Build and Test

```bash
cd /Users/user/repos/gcore-terraform

# Build provider
go build -o terraform-provider-gcore

# Test with flat syntax
cd test-secgroup-rule-edgecases

# Create test file
cat > test-flat-syntax.tf <<'EOF'
terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {
  # Uses environment variables
}

resource "gcore_cloud_security_group" "test_flat" {
  project_id  = 379987
  region_id   = 76
  name        = "test-flat-syntax-workaround"
  description = "Testing flat syntax workaround"
}
EOF

# Set environment
export TF_CLI_CONFIG_FILE="/Users/user/repos/gcore-terraform/.terraformrc"
set -o allexport
source /Users/user/repos/gcore-terraform/.env
set +o allexport

# Test
terraform init
terraform plan

# Should show:
# + resource "gcore_cloud_security_group" "test_flat" {
#     + id          = (known after apply)
#     + name        = "test-flat-syntax-workaround"
#     + description = "Testing flat syntax workaround"
#     + project_id  = 379987
#     + region_id   = 76
#   }
```

---

## Step 5: Test Backward Compatibility

```bash
# Test nested syntax still works
cat > test-nested-syntax.tf <<'EOF'
resource "gcore_cloud_security_group" "test_nested" {
  project_id = 379987
  region_id  = 76
  security_group = {
    name        = "test-nested-syntax-workaround"
    description = "Testing nested syntax still works"
  }
}
EOF

terraform plan
# Should work without errors
```

---

## Step 6: Test Validation

```bash
# Test conflict detection
cat > test-conflict.tf <<'EOF'
resource "gcore_cloud_security_group" "test_conflict" {
  project_id = 379987
  region_id  = 76
  name       = "test-flat"
  security_group = {
    name = "test-nested"
  }
}
EOF

terraform plan
# Should show error:
# Error: Conflicting security group configuration
# Cannot use both 'security_group' block and flat 'name'/'description' attributes.
```

---

## Maintenance

**After running `stainless generate`:**
1. Repeat Steps 1-3 (the changes will be overwritten)
2. Consider creating a git patch file for easier reapplication:
   ```bash
   git diff > patches/flat-syntax-workaround.patch
   # Next time: git apply patches/flat-syntax-workaround.patch
   ```

**Long-term solution:**
- Update `/Users/user/repos/gcore-config/openapi.yml`
- Change `CreateSecurityGroupSerializer` to use `SingleCreateSecurityGroupSerializer`
- Regenerate provider with Stainless
- Remove this workaround

---

## Summary of Changes

| File | Change | Lines |
|------|--------|-------|
| `model.go` | Change `Name` field tag | ~19 |
| `model.go` | Change `Description` field tag | ~23 |
| `model.go` | Replace `MarshalJSON` method | ~31 |
| `model.go` | Replace `MarshalJSONForUpdate` method | ~35 |
| `resource.go` | Replace `ModifyPlan` method | ~286 |

**Total code changes:** ~70 lines added/modified

**Benefit:** Users can now use clean flat syntax! 🎉
