# Stainless Terraform Provider Post-Generation Fix Guide

## Overview

This guide documents common issues that arise from Stainless code generation for the Gcore Terraform provider and the custom code patterns needed to fix them. Based on analysis of 11 PRs authored by Pedro Oliveira and QA feedback from Kirill Tsaregorodtsev.

---

## Raw QA Feedback from Kirill (Jira Comments)

### GCLOUD2-21138 - Instance Resource

**Issue 1: Volume Extend forces replacement**
> "Extend Volume action tries to replace the instance. TF should use endpoint: `POST /v1/volumes/{project_id}/{region_id}/{volume_id}/extend`"

**Issue 2: Flavor change forces replacement**
> "Change Configuration (change the flavor) tries to replace the instance. TF should use endpoint: `POST /v1/instances/{project_id}/{region_id}/{instance_id}/changeflavor`"

**Issue 3: Wrong JSON field name**
> "VM instance creation fails when using an existing volume. Terraform sends an incorrect field in the request payload: `id` is sent instead of `volume_id`"

**Issue 4: Volume not deleted with instance**
> "When deleting a VM instance, the attached volume should be deleted automatically as well."
> "In the old Terraform provider, we passed only the volume_id to the instance creation resource, which meant that when the manifest was destroyed, all volumes were deleted as well. In the new Terraform provider, however, volumes are not deleted after the instance is removed."

**Issue 5: GET with request body**
> "endpoint GET /v1/instances/{project_id}/{region_id}/{instance_id} sends with payload (with req. body)."

**Issue 6: Documentation formatting**
> "Need to fix documentation cloud_instance.md Nested Schema for volumes section where volume types: standard, ssd_hiiops, and ssd_lowlatency look as attributes."

### GCLOUD2-20775 - Region Data Source

**Issue: Undocumented region_id field**
> "I checked the data source gcore_region and noticed an undescribed parameter in the response `region_id = 76`. Additionally, I couldn't find it in the API documentation. The data source gcore_cloud_regions doesn't contain region_id in the response."

### GCLOUD2-20774 - Project Data Source

**Issue: ID field naming inconsistency**
> "Same behavior as gcore_cloud_region - returns region_id"
> "The state in data sources will have both id and resource_id (region_id in this case)."

### GCLOUD2-22597 - CDN Origin Group

**Issue: Write-only/auth field handling**
> Pedro: "Left some comments in the PR... all related to the same handling of auth/write-only related attributes."

---

## Categorized Problems Summary

### Category 1: RequiresReplace When Should Update In-Place
- **Volume size changes** → Use `/extend` endpoint
- **Flavor changes** → Use `/changeflavor` endpoint
- **Interface changes** → Use `Attach`/`Detach` endpoints
- **Security group changes** → Use `Assign`/`Unassign` endpoints
- **Pool changes** → Use pool-specific CRUD endpoints

**Pattern**: Remove `RequiresReplace` plan modifier, add custom Update logic.

### Category 2: JSON Field Mapping Errors
- `volume_id` mapped incorrectly to `id` in request
- API returns `id` but model expects `volume_id`

**Pattern**: Check model JSON tags match API request/response fields.

### Category 3: Async Operations Not Polling
- Create/Update/Delete return immediately without waiting for task completion
- GET requests sending request body during polling

**Pattern**: Use `*AndPoll` methods, pass body via SDK params not `WithRequestBody`.

### Category 4: State Drift on Read
- `tasks`, `task_id` change every read
- Server-managed labels cause drift
- Computed fields showing as changes

**Pattern**: Remove drift-prone fields, filter server-managed data, use `no_refresh` tag.

### Category 5: Sensitive Data in State
- Passwords stored in plaintext state
- Passwords not returned by API causing drift

**Pattern**: Use write-only fields (`password_wo` + `password_wo_version`).

### Category 6: Nested Resource Lifecycle
- Volumes not deleted when parent instance deleted
- Pools/instances not managed individually

**Pattern**: Implement explicit delete logic for nested resources.

### Category 7: ID Field Inconsistencies
- `id` vs `resource_id` naming confusion
- Missing ID field setup breaks imports

**Pattern**: Set `data.ID = data.ResourceSpecificID` after operations.

### Category 8: Documentation Errors
- Enum values displayed as attributes
- Nested schema headers incorrect

**Pattern**: Review generated docs for formatting issues.

---

## Problem Categories and Solutions

### 1. Async Operations Not Polling (Critical)

**Problem**: Generated code calls API methods that return task IDs but doesn't wait for completion.

**Solution**: Replace standard methods with polling variants:

```go
// BEFORE (broken)
_, err = r.client.Cloud.Resource.New(ctx, params, ...)

// AFTER (working)
resource, err := r.client.Cloud.Resource.NewAndPoll(ctx, params, ...)
```

**Affected Methods**:
- `New` → `NewAndPoll`
- `Update` → `UpdateAndPoll`
- `Delete` → `DeleteAndPoll`
- `Action` → `ActionAndPoll`
- `Resize` → `ResizeAndPoll`
- `Upgrade` → `UpgradeAndPoll`

**Files**: `resource.go` in each service directory

**PRs**: #19, #34, #39, #46, #52, #63, #67, #69

---

### 2. Sensitive Write-Only Fields

**Problem**: Passwords and secrets are stored in state and cause drift on read (API doesn't return them).

**Solution**: Use Terraform 1.11+ write-only fields:

**Schema Changes** (`schema.go`):
```go
// BEFORE
"password": schema.StringAttribute{
    Required: true,
    PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
}

// AFTER
"password_wo": schema.StringAttribute{
    Description: "Password (write-only)",
    Required:    true,
    WriteOnly:   true,  // New in TF 1.11
},
"password_wo_version": schema.Int64Attribute{
    Description: "Version to trigger password updates",
    Required:    true,
},
```

**Model Changes** (`model.go`):
```go
// Change tfsdk tag and add no_refresh
Password          types.String `tfsdk:"password_wo" json:"password,optional,no_refresh"`
PasswordWoVersion types.Int64  `tfsdk:"password_wo_version"`
```

**Resource Changes** (`resource.go`):
```go
// Read from config explicitly (write-only fields not in plan)
resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("password_wo"), &data.Password)...)
```

**PRs**: #39, #52, #69

---

### 3. Computed Fields Causing State Drift

**Problem**: Fields like `tasks`, `task_id` change on every read, causing perpetual drift.

**Solution**: Remove from model or mark as `no_refresh`:

```go
// REMOVE these fields from model entirely:
// TaskID   types.String `tfsdk:"task_id" json:"task_id,computed"`
// Tasks    customfield.List[types.String] `tfsdk:"tasks" json:"tasks,computed,no_refresh"`

// OR mark with no_refresh if needed:
Tags *map[string]types.String `tfsdk:"tags" json:"tags,optional,no_refresh"`
```

**PRs**: #67, #69

---

### 4. Custom Update Logic

**Problem**: Single Update method can't handle all update scenarios (different API calls for different field changes).

**Solution**: Check what changed and call appropriate APIs:

```go
func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    // Get both state and plan
    var data, state *ResourceModel
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
    resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

    var stateHasChanged bool

    // 1. Check for name change
    if !data.Name.Equal(state.Name) {
        // Call Update API
        stateHasChanged = true
    }

    // 2. Check for tags change
    if !custom.TagsEqual(data.Tags, state.Tags) {
        // Call Action API with OfUpdateTags
        stateHasChanged = true
    }

    // 3. Check for size/count change
    if !data.Size.Equal(state.Size) {
        // Call Resize API
        stateHasChanged = true
    }

    if stateHasChanged {
        resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
    }
}
```

**Helper for tags comparison** (`internal/custom/tags.go`):
```go
func TagsEqual(a, b *map[string]types.String) bool {
    if a == nil && b == nil { return true }
    if a == nil || b == nil { return false }
    // Compare maps...
}
```

**PRs**: #34, #46, #67, #69

---

### 5. Nested Resource Management

**Problem**: Complex resources have nested entities (pools, instances, access rules) that need separate API calls.

**Solution**: Manage nested resources individually:

**Pattern for K8s Pools** (PR #67):
```go
// Compute differences
oldPools := poolsToMap(state.Pools)
newPools := poolsToMap(data.Pools)

// 1. Create new pools first (before deletes)
for poolName, newPool := range newPools {
    if _, exists := oldPools[poolName]; !exists {
        createPool(ctx, r.client, clusterName, projectID, regionID, newPool)
    }
}

// 2. Handle replacements and updates
for poolName, newPool := range newPools {
    oldPool, exists := oldPools[poolName]
    if !exists { continue }
    if poolNeedsReplace(oldPool, newPool) {
        deletePool(ctx, ...)
        createPool(ctx, ...)
    } else if poolNeedsUpdate(oldPool, newPool) {
        updatePool(ctx, ...)
    }
}

// 3. Delete removed pools last
for poolName := range oldPools {
    if _, exists := newPools[poolName]; !exists {
        deletePool(ctx, ...)
    }
}
```

**Pattern for Placement Group Instances** (PR #34):
```go
// Add instances
r.client.Cloud.Instances.AddToPlacementGroupAndPoll(ctx, instanceID, params)

// Remove instances
r.client.Cloud.Instances.RemoveFromPlacementGroupAndPoll(ctx, instanceID, params)
```

**PRs**: #34, #46, #67

---

### 6. ID Field Handling

**Problem**: Missing or incorrect ID field setup breaks imports and state management.

**Solution**: Standard ID pattern:

**Model**:
```go
type ResourceModel struct {
    ID        types.String `tfsdk:"id" json:"-,computed"`
    // Other fields...
}
```

**In Create/Read/Update**:
```go
data.ID = data.Name  // or data.ServergroupID, etc.
```

**ImportState**:
```go
func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    path_project_id := int64(0)
    path_region_id := int64(0)
    path_id := ""
    importpath.ParseImportID(req.ID, "<project_id>/<region_id>/<id>", &path_project_id, &path_region_id, &path_id)

    data.ProjectID = types.Int64Value(path_project_id)
    data.RegionID = types.Int64Value(path_region_id)
    data.ID = types.StringValue(path_id)
    // Read full resource...
}
```

**PRs**: #34, #67

---

### 7. Replace vs Update API Methods

**Problem**: Some resources require full replacement, not partial updates.

**Solution**: Use Replace method:

```go
// For resources that need full replacement (e.g., secrets, credentials)
_, err = r.client.Cloud.Inference.Secrets.Replace(
    ctx,
    data.Name.ValueString(),
    params,
    option.WithRequestBody("application/json", dataBytes),
)
```

**PRs**: #39, #52

---

### 8. API Method Naming Changes

**Problem**: Go SDK method names change (e.g., `Add`/`Remove` → `New`/`Delete`).

**Solution**: Update to new method names:

```go
// BEFORE
r.client.Cloud.LoadBalancers.Pools.Members.AddAndPoll(...)
r.client.Cloud.LoadBalancers.Pools.Members.RemoveAndPoll(...)

// AFTER
r.client.Cloud.LoadBalancers.Pools.Members.NewAndPoll(...)
r.client.Cloud.LoadBalancers.Pools.Members.DeleteAndPoll(...)
```

**PRs**: #63

---

### 9. Duplicate Provider Registrations

**Problem**: Same resource/datasource registered multiple times in provider.go.

**Solution**: Remove duplicates:

```go
// Check for duplicates in Resources() and DataSources() methods
func (p *Provider) Resources(ctx context.Context) []func() resource.Resource {
    return []func() resource.Resource{
        // Remove duplicate entries like:
        // dns_zone.NewResource,  // Keep only one
    }
}
```

**PRs**: #64

---

### 10. Custom Validators

**Problem**: Complex validation rules can't be expressed in schema alone.

**Solution**: Implement resource validators:

```go
// Validator for conditional requirements
type poolFlavorValidator struct{}

func (v *poolFlavorValidator) ValidateResource(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
    // Check if flavor is VM-based vs baremetal
    // Require different fields based on flavor type
}

// Register in schema:
func (r *Resource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
    return []resource.ConfigValidator{
        &slurmValidator{},
        &authenticationValidator{},
        &poolFlavorValidator{},
    }
}
```

**PRs**: #67

---

### 11. Custom Plan Modifiers

**Problem**: State/plan handling needs custom logic.

**Solution**: Implement plan modifiers:

```go
// Preserve null state values when plan is unknown
type useStateForUnknownIncludingNullStringModifier struct{}

func (m useStateForUnknownIncludingNullStringModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
    if !req.PlanValue.IsUnknown() {
        return
    }
    resp.PlanValue = req.StateValue  // Preserve even if null
}

// Reorder list to match state order
type poolsNormalizeOrderModifier struct{}
```

**PRs**: #67

---

### 12. Filter Server-Managed Data

**Problem**: Server adds data (labels, tags) that causes drift.

**Solution**: Filter server-managed data after reads:

```go
func (m *ResourceModel) FilterServerManagedLabels(ctx context.Context) {
    for _, pool := range *m.Pools {
        for k := range labelsMap {
            if strings.HasPrefix(k, "gcorecluster.x-k8s.io/") {
                delete(filteredLabels, k)
            }
        }
    }
}

// Call after every read
data.FilterServerManagedLabels(ctx)
```

**PRs**: #67

---

## Stainless Config (gcore-config)

Configuration in `openapi.stainless.yml` affects generation:

```yaml
# Set field configurability for Terraform
x-stainless:
  components:
    schemas:
      CreateSubnetSerializer:
        properties:
          gateway_ip:
            x-stainless-terraform-configurability: computed_optional
```

Values:
- `required` - Must be provided
- `optional` - Can be provided
- `computed_optional` - Computed by API if not provided
- `computed` - Always computed, never provided

---

## PR Review Checklist

When reviewing Stainless-generated Terraform PRs:

1. **Async Operations**: Are all CRUD operations using `*AndPoll` methods?
2. **Sensitive Fields**: Are passwords/secrets using write-only pattern?
3. **Drift Fields**: Are `tasks`/`task_id` removed or marked `no_refresh`?
4. **Update Logic**: Does Update handle all field changes correctly?
5. **Nested Resources**: Are nested entities managed separately?
6. **ID Field**: Is ID properly set and ImportState implemented?
7. **Duplicates**: Any duplicate registrations in provider.go?
8. **Validators**: Are complex rules validated properly?
9. **Plan Modifiers**: Any state preservation issues?
10. **Docs**: Are docs accurate for write-only fields?

---

## Files Commonly Modified

Per resource:
- `internal/services/<resource>/resource.go` - Main CRUD logic
- `internal/services/<resource>/schema.go` - Terraform schema
- `internal/services/<resource>/model.go` - Data models
- `internal/services/<resource>/data_source.go` - Data source
- `docs/resources/<resource>.md` - Documentation

Provider-wide:
- `internal/provider.go` - Resource registrations
- `internal/custom/` - Shared helper functions
- `go.mod` / `go.sum` - Dependencies

---

## Next Steps

To apply this knowledge to new fixes:

1. Check the Jira ticket for specific resource/issue
2. Find the corresponding PR in stainless-sdks/gcore-terraform
3. Review the diff against this guide
4. Apply relevant patterns from similar past fixes
5. Test with actual Terraform apply/plan cycles
