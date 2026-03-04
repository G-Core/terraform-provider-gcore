# LB Listener Drift Fix - Current Status

## Issue Summary
Load Balancer Listener resource was experiencing drift on second `terraform apply` due to optional fields being sent as `null` in PATCH requests.

## Affected Fields
- `sni_secret_id` (list)
- `timeout_client_data` (int64)
- `timeout_member_connect` (int64)
- `timeout_member_data` (int64)
- `user_list` (list of objects)

## Solution Progress

### ✅ Completed Steps

1. **Identified root cause**: Fields marked as `Optional: true` only, not `Computed: true`

2. **Found correct OpenAPI schema**: `PatchLbListenerSerializer` (V2 API update operation)

3. **Updated api-schemas config.yaml**:
   - File: `/Users/user/repos/api-schemas/scripts/config.yaml`
   - Added `PatchLbListenerSerializer` with all 5 fields marked as `computed_optional`
   - Commit: `833517e` - "Add computed_optional for LB Listener optional fields"

4. **Attempted manual fix**: Added `Computed: true` to schema.go
   - **Result**: FAILED with "Value Conversion Error"
   - **Reason**: List types need `types.List` not `*[]types.String` to handle unknown values
   - **Action**: Reverted manual changes

5. **Documented the issue**:
   - `PROPER_FIX.md` - Step-by-step guide for the proper CI/CD approach
   - `WHY_MANUAL_FIX_FAILED.md` - Detailed explanation of why manual edits don't work
   - `FIX_APPLIED.md` - Original temporary fix documentation

### ⏳ Waiting For

1. **CI Pipeline** to regenerate `openapi.yaml`:
   - Will add `x-stainless-terraform-configurability: computed_optional` to all 5 fields
   - Triggered by commit `833517e` in api-schemas

2. **Stainless Code Generation**:
   - Will generate proper schema.go with `Computed: true, Optional: true`
   - Will update model.go to use proper types (`types.List` instead of `*[]types.String`)
   - Will update JSON tags from `optional` to `computed_optional`
   - Will commit changes to `gcore-terraform` main branch

3. **Rebase this branch**:
   - After Stainless completes, rebase on main to get generated code
   - Generated code will properly handle computed_optional fields

## Next Actions

1. **Monitor api-schemas CI**: Check if `openapi.yaml` has been regenerated
   ```bash
   cd /Users/user/repos/api-schemas
   git pull
   grep -A 5 "PatchLbListenerSerializer:" openapi.yaml
   # Look for x-stainless-terraform-configurability: computed_optional
   ```

2. **Monitor gcore-terraform main branch**: Check for Stainless commits
   ```bash
   cd /Users/user/repos/gcore-terraform
   git fetch origin main
   git log origin/main --oneline | head -10
   # Look for Stainless automated commits
   ```

3. **Rebase after generation completes**:
   ```bash
   git rebase origin/main
   # Accept generated file changes from main
   ```

4. **Test the fix**:
   ```bash
   cd test-lblistener-drift
   terraform apply
   terraform apply  # Second apply should show "No changes"
   ```

## Why This Approach Is Necessary

Manual editing of generated files doesn't work because:
- Schema changes require coordinated updates across multiple files
- List/slice types need special handling for unknown values
- JSON tags must match schema definitions
- Code generator knows all the places that need updates

The proper approach using `config.yaml` ensures:
- Changes won't be overwritten by future code generation
- All necessary files are updated consistently
- Proper types and plan modifiers are used
- Single source of truth in api-schemas repository

## References

- Process documentation: `MAKE_FIELD_OPTIONAL_COMPUTED.md`
- OpenAPI operation: `LoadBalancerListenerInstanceViewSetV2.patch`
- INPUT schema: `PatchLbListenerSerializer`
- Config file: `/Users/user/repos/api-schemas/scripts/config.yaml`
- Commit: `833517e` in api-schemas repository
