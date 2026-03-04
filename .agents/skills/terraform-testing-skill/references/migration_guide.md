# Migration Guide: Old SDK to Stainless Framework

## Quick Decision Tree

```
Does old provider use WaitTaskAndReturn?
├─ YES → Use *AndPoll methods (async)
└─ NO → Use direct SDK methods (sync)
```

## Step-by-Step Migration Process

### 1. Analyze Old Provider

**Location**: `old_terraform_provider/gcore/resource_gcore_*.go`

**Key patterns to identify**:

```go
// ASYNC PATTERN - uses tasks
results, err := listeners.Create(client, opts).Extract()
taskID := results.Tasks[0]
listenerID, err := tasks.WaitTaskAndReturnResult(client, taskID, ...)

// SYNC PATTERN - direct response
kp, err := keypairs.Create(client, opts).Extract()
d.SetId(kp.ID)
```

### 2. Map to New SDK Methods

#### Async Resources (Use *AndPoll)

**Create:**
```go
// Old
results, err := routers.Create(client, opts).Extract()
taskID := results.Tasks[0]
routerID, err := tasks.WaitTaskAndReturnResult(client, taskID)

// New
router, err := r.client.Cloud.Routers.NewAndPoll(ctx, params,
    option.WithRequestBody("application/json", dataBytes),
    option.WithMiddleware(logging.Middleware(ctx)),
)
err = apijson.UnmarshalComputed([]byte(router.RawJSON()), &data)
```

**Update:**
```go
// Old
results, err := routers.Update(client, id, opts).Extract()
taskID := results.Tasks[0]
_, err = tasks.WaitTaskAndReturnResult(client, taskID)

// New
router, err := r.client.Cloud.Routers.UpdateAndPoll(ctx, id, params,
    option.WithRequestBody("application/json", dataBytes),
    option.WithMiddleware(logging.Middleware(ctx)),
)
err = apijson.UnmarshalComputed([]byte(router.RawJSON()), &data)
```

**Delete:**
```go
// Old
results, err := routers.Delete(client, id).Extract()
taskID := results.Tasks[0]
_, err = tasks.WaitTaskAndReturnResult(client, taskID)

// New
err := r.client.Cloud.Routers.DeleteAndPoll(ctx, id, params,
    option.WithMiddleware(logging.Middleware(ctx)),
)
```

#### Sync Resources (Direct Methods)

**Create:**
```go
// Old
key, err := keypairs.Create(client, opts).Extract()

// New
key, err := r.client.Cloud.SSHKeys.New(ctx, params,
    option.WithRequestBody("application/json", dataBytes),
)
err = apijson.UnmarshalComputed([]byte(key.RawJSON()), &data)
```

### 3. Handle Special Update Logic

#### Pattern: Update with Unset

Some resources need to clear fields explicitly:

```go
// Check if old provider has Unset logic
if toUnset {
    results, err := pools.Unset(clientV2, d.Id(), unsetOpts)
}
```

**Migration approach**:
1. Check if new SDK has Unset methods
2. If not, use PATCH with explicit null/empty values
3. Test with real infrastructure to verify behavior

#### Pattern: Complex Field Updates

**Routes example**:
```go
// Old - might have special handling
if d.HasChange("routes") {
    old, new := d.GetChange("routes")
    // Complex logic to determine add/remove
}

// New - simplified
if data.Routes.IsNull() {
    // Keep existing routes
} else if len(data.Routes.Elements()) == 0 {
    // Clear all routes - send empty array
    updateData["routes"] = []interface{}{}
} else {
    // Update with new routes
    updateData["routes"] = routesData
}
```

### 4. Model and Schema Updates

#### Model File

```go
type RouterResourceModel struct {
    // Basic fields
    ID   types.String `tfsdk:"id" json:"id,computed"`
    Name types.String `tfsdk:"name" json:"name,required"`
    
    // Computed optional fields (API can compute defaults)
    HTTPMethod types.String `json:"http_method,computed_optional"`
    
    // Optional fields
    Description types.String `json:"description,optional"`
    
    // Complex nested structures
    Routes     types.List `tfsdk:"routes" json:"routes,optional"`
    Interfaces types.List `tfsdk:"interfaces" json:"interfaces,optional"`
}
```

#### Schema File

```go
func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                Computed: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.UseStateForUnknown(),
                },
            },
            "name": schema.StringAttribute{
                Required: true,
            },
            "http_method": schema.StringAttribute{
                Computed: true,  // CRITICAL for computed_optional
                Optional: true,
            },
        },
    }
}
```

### 5. Critical Implementation Points

#### Always Use UnmarshalComputed

```go
// In Read method
pool, err := r.client.Cloud.LoadBalancers.Pools.Get(ctx, poolID, params)
bytes := []byte(pool.RawJSON())
err = apijson.UnmarshalComputed(bytes, &data)  // NOT Unmarshal

// In ImportState
err = apijson.UnmarshalComputed(bytes, &data)  // NOT Unmarshal
```

#### Handle Task Errors

```go
resource, err := r.client.Cloud.Resources.NewAndPoll(ctx, params, options...)
if err != nil {
    // Check if it's a task error
    if strings.Contains(err.Error(), "task failed") {
        resp.Diagnostics.AddError("Resource creation failed", 
            "The creation task failed. Check quotas and parameters.")
    } else {
        resp.Diagnostics.AddError("API error", err.Error())
    }
    return
}
```

### 6. Testing Migration

#### Minimal Test Config

```hcl
resource "gcore_cloud_router" "test" {
  name = "migration-test"
  # Minimal fields to test basic functionality
}
```

#### Drift Test

```bash
# Must pass without drift
terraform apply -auto-approve
terraform plan -detailed-exitcode
# Exit code must be 0
```

#### Update Test

```bash
# Save initial ID
initial_id=$(terraform show -json | jq -r '.values.root_module.resources[0].values.id')

# Update
terraform apply -var="name=updated-name"

# Check ID unchanged (PATCH worked)
new_id=$(terraform show -json | jq -r '.values.root_module.resources[0].values.id')
[ "$initial_id" = "$new_id" ] && echo "PASS" || echo "FAIL"
```

### 7. Common Migration Issues

| Issue | Symptom | Solution |
|-------|---------|----------|
| Wrong polling method | Timeout or immediate return | Use *AndPoll for async |
| Drift on computed fields | Plan always shows changes | Use UnmarshalComputed |
| Update recreates | Resource ID changes | Implement Update method |
| Missing task handling | Random failures | Add proper error handling |
| Field not updating | Changes ignored | Check JSON tags in model |

### 8. Validation Checklist

- [ ] Identify async vs sync operations
- [ ] Map old methods to new SDK
- [ ] Update model with correct JSON tags
- [ ] Update schema with Computed flags
- [ ] Use UnmarshalComputed everywhere
- [ ] Handle special update logic
- [ ] Test with real infrastructure
- [ ] Verify no drift
- [ ] Verify updates use PATCH
- [ ] Document any limitations