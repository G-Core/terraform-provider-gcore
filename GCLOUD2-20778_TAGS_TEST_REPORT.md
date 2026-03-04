# GCLOUD2-20778 Tags Issue Test Report

**Date:** 2025-11-24
**Ticket:** [GCLOUD2-20778](https://jira.gcore.lu/browse/GCLOUD2-20778)
**Branch:** bugfix/terraform-lbpool
**Status:** ❌ **CANNOT REPRODUCE** - Likely Already Fixed

---

## Executive Summary

Attempted to reproduce the tags inconsistency error reported in GCLOUD2-20778:
> "An inconsistency error appears after adding tags to the Load Balancer"
> `Error: Provider produced inconsistent result after apply`
> `.tags_v2: new element 0 has appeared.`

**Result:** The error **could NOT be reproduced** with the current code. Tags work correctly for:
- ✅ Adding tags to LB without tags
- ✅ Modifying existing tags
- ✅ No drift detection after tag operations

---

## Reported Issue

**Error Message:**
```
Error: Provider produced inconsistent result after apply

When applying changes to gcore_cloud_load_balancer.lb, provider
"provider["local.gcore.com/repo/gcore"]" produced an unexpected new value: .tags_v2: new
element 0 has appeared.

This is a bug in the provider, which should be reported in the provider's own issue
tracker.
```

**User Configuration:**
```hcl
resource "gcore_cloud_load_balancer" "lb" {
  project_id = local.project_id[0]
  region_id  = data.gcore_cloud_region.rg.id
  flavor = "lb1-2-4"
  name = "qa-lb-name"
  tags = {
    "qa" = "load-balancer"
  }
}
```

---

## Test Scenario Executed

### Test Environment
- **Directory:** `/Users/user/repos/gcore-terraform/test-lb-tags-issue`
- **Provider:** Latest build from `bugfix/terraform-lbpool` branch
- **Region:** Luxembourg-2 (76)
- **Project:** 379987

### Test Steps

#### Step 1: Create LB Without Tags ✅
```bash
terraform apply -auto-approve
```

**Result:** LB created successfully
- `tags`: null
- `tags_v2`: []

#### Step 2: Add Tags to Configuration ✅
Modified main.tf to add:
```hcl
tags = {
  "qa" = "load-balancer"
}
```

**Result:** Applied successfully, no error
- `tags`: {"qa": "load-balancer"}
- `tags_v2`: [{"key": "qa", "read_only": false, "value": "load-balancer"}]

#### Step 3: Check for Drift ✅
```bash
terraform plan
```

**Result:** "No changes. Your infrastructure matches the configuration."

#### Step 4: Modify Tags ✅
Changed tags to:
```hcl
tags = {
  "qa"          = "load-balancer"
  "environment" = "test"
}
```

**Result:** Applied successfully, no error
- `tags`: {"qa": "load-balancer", "environment": "test"}
- `tags_v2`: [2 elements with both tags]

---

## Schema Analysis

### Tags Configuration in Load Balancer Schema

**File:** `internal/services/cloud_load_balancer/schema.go`

```go
// Line 448 - User input (map)
"tags": schema.MapAttribute{
    Description: "Key-value tags to associate with the resource...",
    Optional:    true,
    ElementType: types.StringType,
},

// Line 895 - API output (list of nested objects)
"tags_v2": schema.ListNestedAttribute{
    Description: "List of key-value tags associated with the resource...",
    Computed:    true,  // Read-only
    CustomType:  customfield.NewNestedObjectListType[CloudLoadBalancerTagsV2Model](ctx),
    PlanModifiers: []planmodifier.List{
        listplanmodifier.UseStateForUnknown(),  // ✅ Prevents drift
    },
    NestedObject: schema.NestedAttributeObject{
        Attributes: map[string]schema.Attribute{
            "key": schema.StringAttribute{
                Description: "Tag key. The maximum size for a key is 255 bytes.",
                Computed:    true,
            },
            "read_only": schema.BoolAttribute{
                Description: "Indicates if the tag is read-only",
                Computed:    true,
            },
            "value": schema.StringAttribute{
                Description: "Tag value. The maximum size for a value is 255 bytes.",
                Computed:    true,
            },
        },
    },
},
```

**Key Points:**
1. ✅ `tags` is `Optional: true` - user can set it
2. ✅ `tags_v2` is `Computed: true` - read-only, API response
3. ✅ `tags_v2` has `UseStateForUnknown()` plan modifier - prevents drift
4. ✅ User sets `tags` (map), API returns `tags_v2` (list)

---

## Why The Bug Cannot Be Reproduced

### Likely Already Fixed

Commit `216b73a` (2025-11-04) fixed drift issues in load balancer resources:

```
fix: resolve configuration drift in load balancer and listener resources

Change Read methods from apijson.Unmarshal to apijson.UnmarshalComputed
to properly handle computed_optional fields. This prevents false drift
detection on second terraform apply.

Changes:
- Load balancer resource: Use UnmarshalComputed in Read method
- Listener resource: Use UnmarshalComputed in Read method
```

**What Changed:**
```go
// OLD - internal/services/cloud_load_balancer/resource.go (before 216b73a)
err = apijson.Unmarshal(bytes, &data)

// NEW - internal/services/cloud_load_balancer/resource.go (after 216b73a)
err = apijson.UnmarshalComputed(bytes, &data)
```

**Why This Fixes Tags:**

- `Unmarshal`: Overwrites all fields from API response, ignoring JSON tags
- `UnmarshalComputed`: Respects `computed_optional` tags, only updates when explicitly set

This prevents `tags_v2` from appearing unexpectedly when it wasn't in the previous state.

---

## Additional Relevant Commits

### f999c1c (2025-11-18)
```
feat(cloud): enhance load balancer resource schema with state management modifiers

Added UseStateForUnknown() plan modifiers to all computed fields in LB,
Listener, Pool, and Pool Member schemas
```

This added the `UseStateForUnknown()` modifier to `tags_v2`, preventing drift.

### 0a16c2d
```
fix(cloud): remove duplicate tags in fips schema
```

Fixed tag handling in other resources.

---

## Current Behavior (Working Correctly)

### Tag Addition Flow

1. **User sets `tags` in config:**
   ```hcl
   tags = { "qa" = "load-balancer" }
   ```

2. **Terraform sends PATCH request:**
   ```json
   {"tags": {"qa": "load-balancer"}}
   ```

3. **API responds with full LB object:**
   ```json
   {
     "tags": [{"key": "qa", "value": "load-balancer", "read_only": false}],
     ...
   }
   ```

4. **Provider processes response:**
   - `UnmarshalComputed` correctly maps API's `tags` array to `tags_v2` field
   - `tags` field (map) remains as user configured it
   - No inconsistency detected

5. **State contains both:**
   - `tags`: Map for user input
   - `tags_v2`: List from API (computed, read-only)

6. **Next plan:**
   - Compares user's `tags` input with state
   - Ignores `tags_v2` changes (has `UseStateForUnknown()`)
   - No drift detected ✅

---

## Conclusion

### Bug Status: LIKELY FIXED ✅

The tags inconsistency error **cannot be reproduced** with the current code on the `bugfix/terraform-lbpool` branch.

**Root cause was fixed by:**
1. Commit `216b73a` (2025-11-04) - Changed to `UnmarshalComputed` in Read method
2. Commit `f999c1c` (2025-11-18) - Added `UseStateForUnknown()` to `tags_v2`

**Test Results:**
- ✅ Adding tags: Works
- ✅ Modifying tags: Works
- ✅ No drift detection: Works
- ✅ No inconsistency errors: Works

### Recommendation

1. **For QA:** Test with the latest provider build from `bugfix/terraform-lbpool` branch
2. **If still seeing errors:** Need more specific reproduction steps:
   - Exact sequence of operations
   - Provider version/commit being tested
   - Whether error occurs on first apply or subsequent apply
   - Any specific tag keys/values that trigger it

3. **For Jira:** This specific tags issue appears to be resolved and can be marked as fixed (verify with QA first)

---

## Artifacts

- **Test configuration:** `test-lb-tags-issue/main.tf`
- **Test output:** `test-lb-tags-issue/apply_with_tags.log`
- **State file:** `test-lb-tags-issue/terraform.tfstate`
- **This report:** `GCLOUD2-20778_TAGS_TEST_REPORT.md`

---

## Related Issues Still Open

While the tags issue appears fixed, the **drift issue** from GCLOUD2-20778 is still reproducible:
- See `GCLOUD2-20778_REPRODUCTION_REPORT.md` for details
- After adding pool, listener shows drift in `stats` and other computed fields
- Multiple resources show "(known after apply)" for fields that shouldn't change
