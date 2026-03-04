# Proper Fix for LB Listener Drift Issue

## Current Situation

We manually edited the schema file to add `Computed: true` to these fields:
- `sni_secret_id`
- `timeout_client_data`
- `timeout_member_connect`
- `timeout_member_data`
- `user_list`

**PROBLEM:** This manual fix will be overwritten the next time Stainless regenerates the code.

## Proper Fix Process

According to `MAKE_FIELD_OPTIONAL_COMPUTED.md`, we should use the CI/CD pipeline:

### Step 1: Update api-schemas config.yaml

Add the following configuration to `/Users/user/repos/api-schemas/scripts/config.yaml`:

```yaml
specs:
  - url: https://api.gcore.com/cloud/docs/openapi.yaml
    product: cloud
    title: Cloud API
    extra_terraform:
      components:
        schemas:
          # ... existing configurations ...

          # Add this new configuration:
          PatchLbListenerSerializer:
            properties:
              sni_secret_id: computed_optional
              timeout_client_data: computed_optional
              timeout_member_connect: computed_optional
              timeout_member_data: computed_optional
              user_list: computed_optional
```

### Step 2: Submit PR to api-schemas Repository

1. Create a branch in the `api-schemas` repository
2. Add the configuration above to `scripts/config.yaml`
3. Submit a PR with title: "Add computed_optional for LB Listener timeout and optional fields"
4. Description:
   ```
   This PR adds computed_optional configuration for Load Balancer Listener fields
   that are causing drift detection when not specified by the user.

   Fields marked as computed_optional:
   - sni_secret_id: Can be null or array of UUIDs
   - timeout_client_data: Frontend client inactivity timeout
   - timeout_member_connect: Backend member connection timeout
   - timeout_member_data: Backend member inactivity timeout
   - user_list: Basic authentication user list

   These fields are optional in the user configuration but can have values
   computed/returned by the API, causing Terraform to detect drift on subsequent
   applies.
   ```

### Step 3: CI Pipeline Automatically Updates OpenAPI

After the PR is merged:
1. CI regenerates `openapi.yaml`
2. Adds `x-stainless-terraform-configurability: computed_optional` to each field:
   ```yaml
   PatchLbListenerSerializer:
     properties:
       sni_secret_id:
         # ... existing field definition ...
         x-stainless-terraform-configurability: computed_optional
       timeout_client_data:
         # ... existing field definition ...
         x-stainless-terraform-configurability: computed_optional
       # ... etc
   ```

### Step 4: Stainless Codegen Runs Automatically

When the OpenAPI spec is updated:
1. Stainless detects the change
2. Regenerates Terraform provider code
3. Commits to `gcore-terraform` main branch
4. Generated code will include:
   - Schema with `Computed: true, Optional: true`
   - Model with `computed_optional` JSON tag

### Step 5: Rebase Branch on Main

After Stainless codegen completes:
```bash
cd /Users/user/repos/gcore-terraform
git fetch origin main
git rebase origin/main

# If there are conflicts in generated files, prefer main's version:
git checkout --theirs internal/services/cloud_load_balancer_listener/schema.go
git checkout --theirs internal/services/cloud_load_balancer_listener/model.go

git rebase --continue
```

## Why This is Better Than Manual Edits

1. **Permanent Fix**: Changes won't be overwritten by code generation
2. **Correct Implementation**: Stainless generates the proper plan modifiers and JSON tags
3. **Single Source of Truth**: Configuration lives in api-schemas, not scattered in code
4. **Consistent Pattern**: Follows the same process used for other resources (Subnet, Router, FloatingIP, etc.)

## Current Status

✅ Manual fix applied to test the solution (temporary)
❌ Proper fix not yet implemented (needs api-schemas PR)

## Action Items

1. [ ] Create branch in `api-schemas` repository
2. [ ] Add `PatchLbListenerSerializer` configuration to `scripts/config.yaml`
3. [ ] Submit PR to `api-schemas`
4. [ ] Wait for CI to regenerate OpenAPI spec
5. [ ] Wait for Stainless to regenerate Terraform code
6. [ ] Rebase this branch on main
7. [ ] Remove manual edits (they will be replaced by generated code)
8. [ ] Test the final generated code

## Related Documentation

- `MAKE_FIELD_OPTIONAL_COMPUTED.md` - Complete workflow guide
- `api-schemas/scripts/config.yaml` - Configuration file location
- OpenAPI operation: `LoadBalancerListenerInstanceViewSetV2.patch`
- INPUT schema: `PatchLbListenerSerializer`
