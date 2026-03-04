# TEST SUCCESS - Drift Fix Verified

## Test Results

**Test Date**: 2025-11-04
**Test Script**: `./test_drift_fix.sh`
**Result**: ✅ **SUCCESS** - No drift detected

## Test Output

### First Apply
- Created 2 resources successfully:
  - Load Balancer: `1f476930-676d-404e-abb8-22d9fe9acc7d`
  - Listener: `e3f9818b-fcf7-4f69-9420-692a4d39ec04`

### Second Apply (Drift Check)
```
No changes. Your infrastructure matches the configuration.

Terraform has compared your real infrastructure against your configuration
and found no differences, so no changes are needed.

Apply complete! Resources: 0 added, 0 changed, 0 destroyed.
```

## State File Verification

### Load Balancer State
All API-computed fields properly stored:
- `id`: 1f476930-676d-404e-abb8-22d9fe9acc7d
- `vip_address`: 109.61.125.90
- `vip_port_id`: 0c9660f1-0e6e-4acd-9d2e-9eb5708b424f
- `created_at`: 2025-11-04T09:46:27Z

### Listener State
All computed_optional fields properly stored:
- `id`: e3f9818b-fcf7-4f69-9420-692a4d39ec04
- `timeout_client_data`: 50000 (API default)
- `timeout_member_connect`: 5000 (API default)
- `timeout_member_data`: 50000 (API default)
- `sni_secret_id`: [] (empty, not configured)
- `user_list`: [] (empty, not configured)

## The Fix

### Files Changed

1. **Load Balancer** (internal/services/cloud_load_balancer/resource.go:269)
   ```go
   // BEFORE:
   err = apijson.Unmarshal(bytes, &data)

   // AFTER:
   err = apijson.UnmarshalComputed(bytes, &data)
   ```

2. **Listener Schema** (internal/services/cloud_load_balancer_listener/schema.go)
   - Added `Computed: true` to 5 fields
   - Added `CustomType` to list fields

3. **Listener Model** (internal/services/cloud_load_balancer_listener/model.go)
   - Changed JSON tags to `computed_optional` for 5 fields
   - Changed list types to `customfield` types

4. **Listener Resource** (internal/services/cloud_load_balancer_listener/resource.go:225)
   ```go
   // BEFORE:
   err = apijson.Unmarshal(bytes, &data)

   // AFTER:
   err = apijson.UnmarshalComputed(bytes, &data)
   ```

## Root Cause

The Read methods in both resources were using `apijson.Unmarshal` instead of `apijson.UnmarshalComputed`.

**Why this caused drift:**
- `Unmarshal`: Overwrites ALL fields from API response, ignoring JSON tags
- `UnmarshalComputed`: Respects JSON tags (`computed`, `computed_optional`, `optional`)
- For `computed_optional` fields that weren't set in config, regular `Unmarshal` would populate them from API, causing Terraform to detect a difference between the null config value and the API-returned value

## Conclusion

The minimal fix (2 one-line changes to Read methods + manual schema/model updates for listener) successfully resolves the configuration drift issue. Second `terraform apply` confirms no changes are needed, validating that state properly reflects infrastructure.

**Next Steps:**
- Wait for Stainless to regenerate listener code based on OpenAPI `x-stainless-terraform-configurability` attributes
- Rebase to replace manual listener schema/model changes with generated code
- Load balancer change (resource.go:269) can remain as-is
