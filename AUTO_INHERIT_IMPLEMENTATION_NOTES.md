# Auto-Inherit project_id/region_id - Implementation Notes

## User Request
Remove repetition of project_id/region_id in security group rules by auto-inheriting from parent security group.

## Problem
Current usage requires repetition:
```hcl
resource "gcore_cloud_security_group_rule" "tcp_range" {
  group_id   = gcore_cloud_security_group.test.id
  project_id = 379987  # Repetition!
  region_id  = 76       # Repetition!
}
```

## Solution: Three Approaches

### ✅ Approach 1: Terraform Reference (RECOMMENDED)
**Status**: Works TODAY, no code changes needed

```hcl
resource "gcore_cloud_security_group_rule" "tcp_range" {
  group_id   = gcore_cloud_security_group.test.id
  project_id = gcore_cloud_security_group.test.project_id  # Reference parent
  region_id  = gcore_cloud_security_group.test.region_id   # Reference parent
}
```

**Pros:**
- ✅ Works immediately
- ✅ Explicit and clear
- ✅ No runtime API calls needed
- ✅ DRY principle - values defined once in parent

**Cons:**
- Still requires typing the reference (but much better than hardcoded values)

---

### ✅ Approach 2: Provider-Level Defaults
**Status**: Already supported via schema

```hcl
# Set environment variables:
# export GCORE_CLOUD_PROJECT_ID=379987
# export GCORE_CLOUD_REGION_ID=76

resource "gcore_cloud_security_group_rule" "tcp_range" {
  group_id = gcore_cloud_security_group.test.id
  # project_id and region_id automatically use env vars
}
```

**How it works:**
- Schema marks fields as `Optional+Computed`
- Provider reads `GCORE_CLOUD_PROJECT_ID` and `GCORE_CLOUD_REGION_ID` from environment
- Values are automatically applied if not specified in config

**Pros:**
- ✅ Cleanest configuration
- ✅ Works for all resources in the project
- ✅ No code changes needed

**Cons:**
- Requires environment variables to be set
- Less explicit (values not visible in config)

---

### ❌ Approach 3: Runtime Auto-Fetch (ATTEMPTED, FLAWED)
**Status**: Implemented but has chicken-and-egg problem

**Attempted implementation:**
```go
// In Create method:
if data.ProjectID.IsNull() || data.RegionID.IsNull() {
    // Fetch parent security group to get project_id/region_id
    if err := r.ensureProjectAndRegion(ctx, data); err != nil {
        return err
    }
}
```

**Problem:**
To fetch the parent security group, we need to call:
```go
r.client.Cloud.SecurityGroups.Get(ctx, groupID, params)
```

But `params` requires project_id and region_id!
```go
params := cloud.SecurityGroupGetParams{
    ProjectID: param.NewOpt(???),  // We don't have this!
    RegionID:  param.NewOpt(???),   // We don't have this!
}
```

This creates a chicken-and-egg problem:
- Need project_id/region_id to fetch the security group
- But we're fetching the security group to GET the project_id/region_id!

**Result**: API call defaults to 0/0, causing 400 Bad Request error.

---

## Recommendations

### For Users
**Best Practice**: Use Terraform references (Approach 1)
```hcl
resource "gcore_cloud_security_group_rule" "https" {
  group_id   = gcore_cloud_security_group.test.id
  project_id = gcore_cloud_security_group.test.project_id
  region_id  = gcore_cloud_security_group.test.region_id

  direction  = "ingress"
  protocol   = "tcp"
  # ...
}
```

**Alternative**: Use environment variables (Approach 2) for cleaner config if all resources use same project/region.

### For Provider Developers
**Action Required**: **REVERT** the runtime auto-fetch code I added

The attempted implementation in Create/Update/Delete should be removed because:
1. It doesn't solve the actual problem
2. It creates API errors (400 Bad Request with 0/0 values)
3. Terraform already provides better solutions (references and env vars)

**Files to revert:**
- `/Users/user/repos/gcore-terraform/internal/services/cloud_security_group_rule/resource.go`
  - Lines 66-75 in Create method
  - Lines 132-141 in Update method
  - Lines 266-275 in Delete method

---

## Security Group Schema Flattening

### Current Structure (Nested)
```hcl
resource "gcore_cloud_security_group" "test" {
  security_group = {
    name        = "my-sg"
    description = "Description"
  }
}
```

### Desired Structure (Flat)
```hcl
resource "gcore_cloud_security_group" "test" {
  name        = "my-sg"
  description = "Description"
}
```

### Root Cause
The API request schema (`CreateSecurityGroupSerializer`) has a nested structure:
```json
{
  "security_group": {
    "name": "my-sg",
    "description": "Description"
  }
}
```

While the response schema (`SecurityGroupSerializer`) is flat:
```json
{
  "id": "...",
  "name": "my-sg",
  "description": "Description"
}
```

### Solution Options

#### Option A: Change OpenAPI Spec (Recommended)
Update `openapi.yml` to use `SingleCreateSecurityGroupSerializer` directly instead of nesting it:

```yaml
# Current:
requestBody:
  content:
    application/json:
      schema:
        $ref: '#/components/schemas/CreateSecurityGroupSerializer'  # Has nesting

# Change to:
requestBody:
  content:
    application/json:
      schema:
        $ref: '#/components/schemas/SingleCreateSecurityGroupSerializer'  # Flat
```

Then update Stainless config to match:
```yaml
# In openapi.stainless.yml
resources:
  cloud:
    subresources:
      security_groups:
        models:
          security_group: '#/components/schemas/SecurityGroupSerializer'  # Response
        methods:
          create:
            endpoint: post /cloud/v1/securitygroups/{project_id}/{region_id}
            request_schema: '#/components/schemas/SingleCreateSecurityGroupSerializer'  # Flat request
```

#### Option B: Stainless Schema Override (If Available)
Check if Stainless supports flattening nested request bodies via config. Documentation research found limited options for this.

#### Option C: Manual Schema Override (Not Recommended)
Manually edit generated schema files, but this will be overwritten on regeneration.

---

## Summary

1. **Auto-inherit project/region**: Use Terraform references or environment variables (no code changes needed)
2. **Flatten security_group block**: Requires OpenAPI spec change or Stainless config update
3. **Action needed**: Revert the runtime auto-fetch implementation I added (it's flawed)

