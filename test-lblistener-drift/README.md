# LB Listener Drift Issue Reproduction

## Issue Description

After applying the Terraform configuration successfully, running `terraform apply` a second time shows that the LB Listener wants to update with null values for optional fields.

## Problem

The listener attempts to send a PATCH request with:
```json
{
    "sni_secret_id": null,
    "timeout_client_data": null,
    "timeout_member_connect": null,
    "timeout_member_data": null,
    "user_list": null
}
```

## Fields Involved

- `sni_secret_id` (optional)
- `timeout_client_data` (optional, default: 50000)
- `timeout_member_connect` (optional, default: 5000)
- `timeout_member_data` (optional, default: 50000)
- `user_list` (optional)

## Test Steps

1. Source environment variables:
   ```bash
   source ../set_env.sh
   ```

2. Initialize Terraform:
   ```bash
   terraform init
   ```

3. Apply configuration (first time):
   ```bash
   terraform apply
   ```

4. Apply configuration again (second time) - this should show drift:
   ```bash
   terraform apply
   ```

## Expected Behavior

Second apply should show "No changes" - infrastructure is up-to-date.

## Actual Behavior

Second apply shows planned updates for both LB and Listener resources with `(known after apply)` values and attempts to send null values for optional fields.
