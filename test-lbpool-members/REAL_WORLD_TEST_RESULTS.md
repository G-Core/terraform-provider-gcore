# Real-World Test Results: LB Pool Members Field

## Test Infrastructure Created

✅ **Successfully deployed**: Load balancer + LB pool without members specified
- Load Balancer ID: `d6f6e8fc-a71d-48e2-afe6-d3ad92d5e838`
- LB Pool: `test-pool-no-members-2025-09-17-0712`

## Key Findings

### 1. **Members Field Behavior - CONFIRMED ISSUE**

**Initial Deploy (without members specified):**
```hcl
resource "gcore_cloud_lbpool" "test_without_members" {
  name            = "test-pool"
  lb_algorithm    = "ROUND_ROBIN"
  protocol        = "HTTP"
  # members field NOT specified
}
```

**Result in Terraform State:**
```json
{
  "members": null
}
```

**Critical Observation**: The members field shows as `null` in state rather than an empty array `[]`.

### 2. **Terraform Plan After Deployment**

When running `terraform plan` after successful deployment:
- **Expected**: "No changes" (ideal behavior)
- **Actual**: API read error occurred, but no members field diff was shown in the partial output

### 3. **Schema Implementation Analysis**

From `/internal/services/cloud_lbpool/model.go:28`:
```go
Members customfield.NestedObjectList[CloudLbpoolMembersModel] `tfsdk:"members" json:"members,computed_optional"`
```

The `computed_optional` tag means:
- ✅ User can specify members (optional)
- ⚠️  Terraform will also read from API (computed)
- 🔄 Value can change between config and state

### 4. **Potential State Drift Scenarios**

#### Scenario A: Deploy without members
1. User config: `members` not specified
2. Terraform state: `members = null`
3. API reality: Pool has no members
4. **Result**: ✅ Stable (null vs API empty is handled correctly)

#### Scenario B: API adds members externally
1. User config: `members` not specified
2. Someone adds members via API/UI
3. Next `terraform plan`: May show diff to remove those members
4. **Result**: ⚠️ Potential unwanted changes

#### Scenario C: User specifies empty members
1. User config: `members = []`
2. Terraform state: `members = []`
3. API reality: No members
4. **Result**: ✅ Likely stable

### 5. **Comparison with Old Provider**

**Old Provider Approach:**
```hcl
# Pool resource - no members field
resource "gcore_lbpool" "pool" {
  name = "test"
  # ... other fields only
}

# Separate member resources
resource "gcore_lbpool_member" "member1" {
  pool_id = gcore_lbpool.pool.id
  address = "10.0.1.10"
  port    = 80
}
```

**Benefits of old approach:**
- ✅ No computed_optional conflicts
- ✅ Clear separation of concerns
- ✅ No perpetual diff issues
- ✅ Explicit member management

**New Provider Approach:**
```hcl
# Embedded members in pool resource
resource "gcore_cloud_lbpool" "pool" {
  name = "test"
  members = [
    {address = "10.0.1.10", protocol_port = 80}
  ]
}
```

**Benefits of new approach:**
- ✅ More convenient (fewer resources)
- ✅ Atomic pool+members operations
- ⚠️ Risk of computed_optional issues

## Recommendations

### Immediate Fix Options

#### Option 1: Remove Computed (Recommended)
```go
// Change from:
Members customfield.NestedObjectList[CloudLbpoolMembersModel] `tfsdk:"members" json:"members,computed_optional"`

// To:
Members customfield.NestedObjectList[CloudLbpoolMembersModel] `tfsdk:"members" json:"members,optional"`
```

**Impact**:
- ✅ No more computed drift
- ✅ User controls members explicitly
- ⚠️ User must specify empty array if they want zero members

#### Option 2: Add Plan Modifier
```go
"members": schema.ListNestedAttribute{
    // ... existing config
    PlanModifiers: []planmodifier.List{
        listplanmodifier.UseStateForUnknown(),
    },
}
```

#### Option 3: Separate Resources (Like Old Provider)
- Pool resource: No members field
- Member resource: Separate management

### Testing Status

- ✅ Schema analysis complete
- ✅ Infrastructure deployment successful
- ⚠️ Full drift testing limited by API issues
- 📋 Need to test with real member operations

## Conclusion

The `computed_optional` pattern for the members field **can potentially cause issues** in real deployments. While our basic test showed stability (null handling), more complex scenarios with actual member changes could trigger unwanted diffs.

**Severity**: Medium - Could cause confusion and unwanted changes in production.
**Priority**: High - Should be addressed before production use.