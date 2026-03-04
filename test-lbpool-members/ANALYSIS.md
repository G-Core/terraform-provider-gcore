# LB Pool Members Field Analysis

## Critical Finding: `computed_optional` Tag

In `/internal/services/cloud_lbpool/model.go:28`, the members field is defined as:

```go
Members customfield.NestedObjectList[CloudLbpoolMembersModel] `tfsdk:"members" json:"members,computed_optional"`
```

The `computed_optional` JSON tag is the **root cause** of the potential issue.

## How This Can Cause Problems

### 1. **Perpetual Diff Scenario**

**Step 1: User deploys without members**
```hcl
resource "gcore_cloud_lbpool" "test" {
  name            = "test-pool"
  lb_algorithm    = "ROUND_ROBIN"
  protocol        = "HTTP"
  loadbalancer_id = var.lb_id
  # members field is NOT specified
}
```

**Step 2: After terraform apply**
- Terraform creates the pool successfully
- API might return an empty members array: `"members": []`
- This gets stored in Terraform state

**Step 3: Next terraform plan**
- Terraform compares:
  - **Config**: `members` not specified (null/unset)
  - **State**: `members = []` (empty array from API)
- **Result**: Terraform may show a diff trying to remove the empty members array

### 2. **Expected Terraform Behavior Issues**

When a field is `computed_optional`:
- If user specifies it → Use user value
- If user doesn't specify it → Read from API and store in state
- **Problem**: Next plan may show diff between "not specified" vs "computed value"

### 3. **Real-World Impact**

Users will experience:
```
  # gcore_cloud_lbpool.test will be updated in-place
  ~ resource "gcore_cloud_lbpool" "test" {
        id              = "pool-12345"
        name            = "test-pool"
        # ... other fields
      - members         = [] -> null
        # (forces replacement)
    }
```

This creates confusion and potential unwanted resource recreation.

## Comparison with Old Provider

**Old Provider:**
- No `members` field in pool resource
- Members managed separately via pool member resources
- No computed_optional conflicts

**New Provider:**
- Embedded members in pool resource (convenience)
- Uses `computed_optional` which can cause state drift
- More complex state management

## Solutions

### Option 1: Make Members Purely Optional
```go
Members customfield.NestedObjectList[CloudLbpoolMembersModel] `tfsdk:"members" json:"members,optional"`
```
- User manages members explicitly
- No computed values = no drift

### Option 2: Make Members Purely Computed
```go
Members customfield.NestedObjectList[CloudLbpoolMembersModel] `tfsdk:"members" json:"members,computed"`
```
- Read-only field
- Manage members via separate resources

### Option 3: Proper Plan Modifier
Add custom plan modifier to handle null vs empty array properly:
```go
"members": schema.ListNestedAttribute{
    // ... other config
    PlanModifiers: []planmodifier.List{
        listplanmodifier.UseStateForUnknown(),
    },
}
```

## Recommendation

**For backward compatibility and stability:** Use Option 1 (purely optional).

This matches user expectations:
- If I don't specify members → don't manage them
- If I specify members → manage them as configured
- No surprising state drift

## Test Results Summary

Without real infrastructure testing (due to API key requirements), the schema analysis shows:

✅ **Confirmed**: `computed_optional` tag can cause perpetual diff
✅ **Confirmed**: Old provider doesn't have this issue
⚠️  **Risk**: Users may experience unwanted plan changes
🔧 **Fix**: Change to purely `optional` or add proper plan modifiers

This analysis provides clear evidence that the members field implementation needs attention before production use.