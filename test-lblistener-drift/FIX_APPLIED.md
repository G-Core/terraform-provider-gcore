# Fix Applied for LB Listener Drift Issue

## Problem
The Load Balancer Listener resource was experiencing drift detection on the second `terraform apply` run. The following optional fields were being sent as `null` in PATCH requests:
- `sni_secret_id`
- `timeout_client_data`
- `timeout_member_connect`
- `timeout_member_data`
- `user_list`

## Root Cause
These fields were marked as `Optional: true` only in the schema, but not `Computed: true`. This caused Terraform to not expect values from the API for these fields, leading to state drift when the API returned values.

## Solution
Added `Computed: true` to all five fields in the schema file:
`internal/services/cloud_load_balancer_listener/schema.go`

### Changes Made:
1. **timeout_client_data** (line 83-90): Added `Computed: true`
2. **timeout_member_connect** (line 91-98): Added `Computed: true`
3. **timeout_member_data** (line 99-106): Added `Computed: true`
4. **sni_secret_id** (line 112-117): Added `Computed: true`
5. **user_list** (line 118-132): Added `Computed: true`

## Verification
After applying the fix, the terraform plan now correctly shows these fields as `(known after apply)` instead of trying to set them to `null`:

```
# gcore_cloud_load_balancer_listener.ls will be created
resource "gcore_cloud_load_balancer_listener" "ls" {
  ...
  sni_secret_id          = (known after apply)
  timeout_client_data    = (known after apply)
  timeout_member_connect = (known after apply)
  timeout_member_data    = (known after apply)
  user_list              = (known after apply)
}
```

This indicates that Terraform now understands these fields are computed by the API and will not attempt to send null values for them on subsequent applies.

## Same Pattern as Previous Fix
This fix follows the same pattern as the fix applied to:
- Load Balancer Pool resource
- Load Balancer resource
- Router resource

All optional fields that can be computed by the API should be marked as both `Optional: true` and `Computed: true` to prevent drift detection.
