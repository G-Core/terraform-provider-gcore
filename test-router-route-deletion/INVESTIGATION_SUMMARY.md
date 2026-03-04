# Router Route Deletion Bug Investigation

## Issue Summary
**JIRA**: GCLOUD2-21144
**Reporter**: Kirill Tsaregorodtsev
**Date**: November 3, 2025

When a route is removed from a `gcore_cloud_network_router` configuration, Terraform reports "Apply complete! Resources: 0 added, 1 changed, 0 destroyed" but the route is **NOT actually deleted** from the router on the GCore API.

## Reproduction Steps

1. Create a router with an interface:
```hcl
resource "gcore_cloud_network_router" "router" {
  project_id = local.project_id[0]
  region_id  = data.gcore_cloud_region.rg.id
  name       = "qa-terr-rename"
  external_gateway_info = {
    enable_snat = true
  }
  interfaces = [{
    subnet_id = gcore_cloud_network_subnet.sb.id
    type      = "subnet"
  }]
}
```

2. Add a route to the router:
```hcl
resource "gcore_cloud_network_router" "router" {
  # ... same as above ...
  routes = [{
    destination = "10.0.3.0/24"
    nexthop     = "192.168.0.1"
  }]
}
```

3. Remove the `routes` block from the configuration

4. Run `terraform apply`

**Expected**: Route should be deleted from the router
**Actual**: Terraform says "Apply complete" but route remains on the router

## Code Analysis

### File: internal/services/cloud_network_router/resource.go

The `Update` function (lines 109-298) handles router updates:

```go
func (r *CloudNetworkRouterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    // ... get plan and state ...

    // Handle interface attach/detach operations (lines 128-222)
    // Interfaces use special attach/detach API endpoints

    // Update other router attributes (name, routes, external_gateway_info)
    needsUpdate := !data.Name.Equal(state.Name) ||
        !data.Routes.Equal(state.Routes) ||
        !data.ExternalGatewayInfo.Equal(state.ExternalGatewayInfo)

    if needsUpdate {
        // Send PATCH request with updated fields
        dataBytes, err = data.MarshalJSONForUpdate(*state)
        // ... PATCH /v1/routers/{project_id}/{region_id}/{router_id} ...
    }
}
```

### Key Finding

**Interfaces** are handled via special API endpoints:
- Attach: `POST /v1/routers/{project_id}/{region_id}/{router_id}/attach`
- Detach: `POST /v1/routers/{project_id}/{region_id}/{router_id}/detach`

**Routes** are handled via PATCH request:
- `PATCH /v1/routers/{project_id}/{region_id}/{router_id}` with routes in the body

### MarshalJSONForUpdate Behavior

File: `internal/services/cloud_network_router/model.go:34-40`

```go
func (m CloudNetworkRouterModel) MarshalJSONForUpdate(state CloudNetworkRouterModel) (data []byte, err error) {
    // Create a copy of the model to marshal, but force interfaces to equal state
    // so they're not included in the PATCH request (interfaces are managed via attach/detach)
    mCopy := m
    mCopy.Interfaces = state.Interfaces
    return apijson.MarshalForPatch(mCopy, state)
}
```

This marshals the update as a JSON PATCH (RFC 7386) - only changed fields are included.

### Array Encoding in Patch Mode

File: `internal/apijson/encoder.go:294-340`

The `newArrayTypeEncoder` function handles array serialization:

```go
func (e *encoder) newArrayTypeEncoder(t reflect.Type) encoderFunc {
    return func(plan reflect.Value, state reflect.Value) ([]byte, error) {
        stateNil := !state.IsValid() || state.IsNil()
        planNil := !plan.IsValid() || plan.IsNil()

        if stateNil && planNil {
            return nil, nil  // Both nil: omit from output
        } else if planNil {
            return explicitJsonNull, nil  // Plan is null: send null to unset
        } else if !stateNil && arrayPatch && reflect.DeepEqual(plan.Interface(), state.Interface()) {
            return nil, nil  // Arrays are equal: omit from output
        }

        // Arrays are different: serialize the plan array
        json := []byte("[]")
        for i := 0; i < plan.Len(); i++ {
            // ... add items to array ...
        }
        return json, nil
    }
}
```

**When routes go from `[{...}]` to `[]`:**
- plan is `[]` (empty array, NOT nil)
- state is `[{...}]` (has routes)
- `DeepEqual` returns false
- The function serializes plan as `[]`
- PATCH request includes: `{"routes": []}`

## Root Cause Hypothesis

The Terraform provider correctly sends `{"routes": []}` in the PATCH request, but there are two possible issues:

1. **API doesn't support route deletion via empty array**: The GCore API might not interpret `routes: []` as "delete all routes". It might require a special endpoint like:
   - `DELETE /v1/routers/{project_id}/{region_id}/{router_id}/routes/{route_id}`
   - Or a specific `null` value instead of empty array

2. **API requires explicit route deletion**: Similar to how interfaces use attach/detach endpoints, routes might need:
   - Add route: `POST /v1/routers/{project_id}/{region_id}/{router_id}/add_route`
   - Delete route: `POST /v1/routers/{project_id}/{region_id}/{router_id}/remove_route`

## Recommended Fix

Check the GCore API documentation for routers:
- Does `PATCH` with `routes: []` delete all routes?
- Does `PATCH` with `routes: null` delete all routes?
- Is there a dedicated endpoint for adding/removing routes?

If the API requires special endpoints for route management (like interfaces do), then the fix is to:

1. Add route add/remove API calls similar to attach/detach for interfaces
2. Update the `Update` function to detect route changes and call the appropriate endpoints
3. Exclude routes from the PATCH request body (similar to how interfaces are excluded)

Example fix in `resource.go`:

```go
// Handle route add/delete operations (similar to interfaces)
if !data.Routes.Equal(state.Routes) {
    oldRoutes, _ := state.Routes.AsStructSliceT(ctx)
    newRoutes, _ := data.Routes.AsStructSliceT(ctx)

    // Find routes to add
    // Find routes to delete
    // Call appropriate API endpoints
}

// Update: exclude routes from needsUpdate if API has dedicated endpoints
needsUpdate := !data.Name.Equal(state.Name) ||
    !data.ExternalGatewayInfo.Equal(state.ExternalGatewayInfo)
```

## Test Files

- `/Users/user/repos/gcore-terraform/test-router-route-deletion/main.tf` - Terraform manifest
- `/Users/user/repos/gcore-terraform/test-router-route-deletion/test_route_deletion.sh` - Automated reproduction script

Run the test script to reproduce the issue:
```bash
cd test-router-route-deletion
chmod +x test_route_deletion.sh
./test_route_deletion.sh
```
