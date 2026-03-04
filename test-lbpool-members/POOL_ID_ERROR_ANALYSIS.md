# LB Pool "missing required pool_id parameter" Error Analysis

## Executive Summary

The API error "missing required pool_id parameter" occurs during `terraform plan` refresh after creating an LB pool resource. Investigation reveals the ID field is not being properly saved to Terraform state during resource creation, causing subsequent Read operations to fail.

## Error Context

**When it occurs:**
- After successfully creating LB pool resource
- During `terraform plan` when Terraform tries to refresh state
- The Read method receives an empty pool ID from state

**Error location:**
- SDK: `/Users/user/go/pkg/mod/github.com/!g-!core/gcore-go@v0.12.0/cloud/loadbalancerpool.go:171`
- Provider: `/Users/user/repos/gcore-terraform/internal/services/cloud_lbpool/resource.go:178`

## Root Cause Analysis

### 1. SDK Validation (Working Correctly)

The SDK properly validates the pool ID:
```go
// loadbalancerpool.go:171
if poolID == "" {
    err = errors.New("missing required pool_id parameter")
    return
}
```

The LoadBalancerPool struct confirms ID is required:
```go
type LoadBalancerPool struct {
    ID string `json:"id,required" format:"uuid4"`
    // ... other fields
}
```

### 2. SDK's NewAndPoll Method (Working Correctly)

The SDK correctly returns the pool with ID:
```go
func (r *LoadBalancerPoolService) NewAndPoll(...) (*LoadBalancerPool, error) {
    // Creates pool and waits for task
    resourceID := task.CreatedResources.Pools[0]
    // Fetches full pool details with ID
    return r.Get(ctx, resourceID, getParams, opts...)
}
```

### 3. Provider's Create Method (PROBLEM IDENTIFIED)

The new provider's Create method:
```go
pool, err := r.client.Cloud.LoadBalancers.Pools.NewAndPoll(...)
// Uses apijson.UnmarshalComputed to populate fields
err = apijson.UnmarshalComputed([]byte(pool.RawJSON()), &data)
// Saves to state
resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
```

**Issue:** The ID field may not be properly unmarshalled by `apijson.UnmarshalComputed` into the `types.String` field.

### 4. Evidence from Testing

During our infrastructure test:
- Pool creation reported success: "Creation complete after 3s"
- But `terraform show` output showed NO ID field for the pool resource
- State inspection revealed: `"members": null` but no ID field
- This confirms ID was never saved to state

## Comparison with Old Provider

### Old Provider (Working)
```go
// Explicitly extracts ID from task
lbPoolID, err := lbpools.ExtractPoolIDFromTask(taskInfo)
// Directly sets the ID
d.SetId(lbPoolID.(string))
// Immediately reads to verify
resourceLBPoolRead(ctx, d, m)
```

### New Provider (Problematic)
```go
// Relies on automatic unmarshalling
err = apijson.UnmarshalComputed([]byte(pool.RawJSON()), &data)
// Assumes ID is populated
resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
```

## Resource Leakage Risk

**WARNING:** If we simply remove the resource from state when ID is missing:
- ❌ Resource exists in cloud but Terraform loses track
- ❌ Next `terraform apply` creates duplicate resource
- ❌ Original resource becomes orphaned (leaked)
- ❌ User pays for unused resources

## Recommended Solutions

### Solution 1: Fix Create Method (RECOMMENDED)

Explicitly set the ID after creation, matching old provider pattern:

```go
func (r *CloudLbpoolResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    // ... existing code ...

    pool, err := r.client.Cloud.LoadBalancers.Pools.NewAndPoll(...)
    if err != nil {
        resp.Diagnostics.AddError("failed to make http request", err.Error())
        return
    }

    // EXPLICITLY SET ID (like old provider)
    if pool.ID == "" {
        resp.Diagnostics.AddError("Provider error", "Pool created but ID not returned by API")
        return
    }
    data.ID = types.StringValue(pool.ID)

    // Then unmarshal other computed fields
    err = apijson.UnmarshalComputed([]byte(pool.RawJSON()), &data)
    if err != nil {
        resp.Diagnostics.AddError("failed to deserialize http request", err.Error())
        return
    }

    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
```

**Benefits:**
- ✅ Minimal code change
- ✅ Matches proven pattern from old provider
- ✅ Ensures ID is always saved
- ✅ Prevents resource leakage
- ✅ Fixes root cause, not symptoms

### Solution 2: Defensive Check in Read (ADDITIONAL SAFETY)

Add error handling without removing from state:

```go
func (r *CloudLbpoolResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var data *CloudLbpoolModel

    resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

    if resp.Diagnostics.HasError() {
        return
    }

    // Defensive check - should never happen
    if data.ID.IsNull() || data.ID.IsUnknown() || data.ID.ValueString() == "" {
        resp.Diagnostics.AddError(
            "Resource configuration error",
            "Pool resource ID is missing from state. This indicates a provider bug during resource creation. "+
            "The pool may exist in the cloud but Terraform has lost track of it. "+
            "Please check your cloud console and manually import the resource if it exists.",
        )
        return  // DO NOT remove from state to prevent leakage
    }

    // Continue normal read...
}
```

## Implementation Priority

1. **High Priority:** Fix Create method to explicitly set ID
2. **Medium Priority:** Add defensive check in Read method
3. **Low Priority:** Investigate why `apijson.UnmarshalComputed` doesn't set ID field

## Testing Requirements

After implementing fix:
1. Create new LB pool resource
2. Verify ID appears in `terraform show` output
3. Run `terraform plan` - should show no changes
4. Run `terraform refresh` - should complete without errors
5. Verify resource can be updated and deleted

## Related Issues to Check

- Similar pattern may affect other resources (LoadBalancer, Listener, etc.)
- Check if other resources rely on `apijson.UnmarshalComputed` for ID
- Consider standardizing ID handling across all resources

## Conclusion

The error is caused by the ID not being saved to Terraform state during resource creation. The new provider's reliance on `apijson.UnmarshalComputed` appears to be failing to properly set the ID field. The safest fix is to explicitly set the ID after creation, matching the pattern used in the old provider. This prevents resource leakage and ensures Terraform maintains proper state tracking.