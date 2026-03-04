# Test Report: GCLOUD2-20778 Tags Fix Verification

## Summary

| Test | Description | Result |
|------|-------------|--------|
| 1 | Create LB without tags | PASS |
| 2 | Add tags to existing LB (reproduces bug) | PASS |
| 3 | Drift check after adding tags | PASS |
| 4 | Modify tags | PASS |
| 5 | Remove all tags | PASS |
| 6 | Drift check after removing tags | PASS |

**Overall Result: ALL TESTS PASSED**

## Issue Details

**JIRA Ticket**: GCLOUD2-20778
**Error Before Fix**:
```
Error: Provider produced inconsistent result after apply

When applying changes to gcore_cloud_load_balancer.lb, provider
"provider[\"local.gcore.com/repo/gcore\"]" produced an unexpected new value: .tags_v2: new
element 0 has appeared.
```

## Root Cause Analysis

The `tags_v2` field had `UseStateForUnknown()` plan modifier, which told Terraform to expect `tags_v2` to remain unchanged when other fields update. However, `tags_v2` is directly derived from the `tags` input field, so when `tags` changes, `tags_v2` must also change.

This created an inconsistency:
- Plan expected: `tags_v2` stays `[]`
- API returned: `tags_v2` populated with new tags
- Result: Inconsistency error

## Fix Applied

Commit: `0b2dab99f4125f9888398dc9911b92947c0892c9`

```diff
-				PlanModifiers: []planmodifier.List{
-					listplanmodifier.UseStateForUnknown(),
-				},
+				// NOTE: Removed UseStateForUnknown() to fix GCLOUD2-20778 tags inconsistency error
+				// tags_v2 is derived from tags input, so it should always refresh from API
+				PlanModifiers: []planmodifier.List{},
```

## Test Execution

### Test 1: Create LB without tags

```hcl
resource "gcore_cloud_load_balancer" "lb" {
  name           = "qa-lb-tags-test-20778"
  flavor         = "lb1-2-4"
  vip_network_id = gcore_cloud_network.test.id
  vip_subnet_id  = gcore_cloud_network_subnet.test.id
  tags           = {}
}
```

**Result**: Created successfully
**LB ID**: `670f9edb-6e75-4f70-89a1-57d118f5c8aa`

### Test 2: Add tags to existing LB (Critical - reproduces GCLOUD2-20778)

```bash
terraform apply -var='lb_tags={"qa"="load-balancer"}'
```

**Result**: SUCCESS - No inconsistency error!

**Output**:
```
lb_tags = tomap({
  "qa" = "load-balancer"
})
lb_tags_v2 = tolist([
  {
    "key" = "qa"
    "read_only" = false
    "value" = "load-balancer"
  },
])
```

### Test 3: Drift check after adding tags

```bash
terraform plan -detailed-exitcode
```

**Result**: Exit code 0 (No changes detected)

### Test 4: Modify tags

```bash
terraform apply -var='lb_tags={"qa"="modified", "env"="test"}'
```

**Result**: SUCCESS - Tags modified without error

### Test 5: Remove all tags

```bash
terraform apply -var='lb_tags={}'
```

**Result**: SUCCESS - Tags cleared without error

### Test 6: Drift check after removing tags

```bash
terraform plan -detailed-exitcode
```

**Result**: Exit code 0 (No changes detected)

## Cleanup

All test resources destroyed successfully:
- Load Balancer: `670f9edb-6e75-4f70-89a1-57d118f5c8aa`
- Subnet: `b11fcb02-02a6-4497-bbe7-9891b9b0f6fb`
- Network: `1cdc6ddd-b309-4d0b-a617-410e794d84ac`

## Conclusion

The fix for GCLOUD2-20778 has been verified:

1. The inconsistency error no longer occurs when adding tags to an existing load balancer
2. Tags can be added, modified, and removed without errors
3. No drift detected after any tag operations
4. The fix correctly allows `tags_v2` to be refreshed from the API instead of preserving stale state

**Tested on**: 2025-11-25
**Provider Version**: Development build from branch `bugfix/terraform-lbpool`
**Commit**: `0b2dab9`
