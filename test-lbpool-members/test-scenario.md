# Testing Members Field Optional+Computed Behavior

## Test Scenario

The `members` field in the LB Pool resource is defined as both `Optional` and `Computed`:

```go
"members": schema.ListNestedAttribute{
    Description: "Pool members",
    Computed:    true,
    Optional:    true,
    // ...
}
```

## Potential Issues

### 1. Perpetual Diff Problem
When a field is both Optional and Computed:
- User doesn't specify members in config → Terraform computes the value from API
- On next `terraform plan` → Terraform compares computed value with empty config
- This can cause perpetual diff if not handled correctly

### 2. Test Cases to Run

#### Case 1: Deploy without members, then run plan
```hcl
resource "gcore_cloud_lbpool" "test" {
  name            = "test-pool"
  lb_algorithm    = "ROUND_ROBIN"
  protocol        = "HTTP"
  loadbalancer_id = var.lb_id
  # members field is NOT specified
}
```

**Expected behavior**:
- First apply: succeeds
- Second plan: should show "No changes"
- **Problem**: If API returns members and Terraform tries to manage them

#### Case 2: Deploy with empty members, then run plan
```hcl
resource "gcore_cloud_lbpool" "test" {
  name            = "test-pool"
  lb_algorithm    = "ROUND_ROBIN"
  protocol        = "HTTP"
  loadbalancer_id = var.lb_id
  members         = []  # Explicitly empty
}
```

#### Case 3: Deploy with members, then run plan
```hcl
resource "gcore_cloud_lbpool" "test" {
  name            = "test-pool"
  lb_algorithm    = "ROUND_ROBIN"
  protocol        = "HTTP"
  loadbalancer_id = var.lb_id
  members = [{
    address       = "10.0.1.10"
    protocol_port = 80
  }]
}
```

## Analysis from Schema

Looking at the current schema in `internal/services/cloud_lbpool/schema.go:215-267`:

The members field uses:
- `Computed: true` - Terraform will read value from API
- `Optional: true` - User can specify value
- `CustomType: customfield.NewNestedObjectListType[CloudLbpoolMembersModel](ctx)`

This combination suggests the provider attempts to handle the optional+computed pattern, but we need to verify it works correctly in practice.

## Recommended Fix

If issues are found, consider:
1. Making members purely `Optional` (remove `Computed`)
2. Using separate data source for reading existing members
3. Implementing proper ignore/suppress logic for computed changes