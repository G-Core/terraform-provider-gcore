Complete Manual Fix Summary - All 5 Fields

MANUAL CHANGES APPLIED

All 5 fields now have proper computed_optional configuration:

1. Schema Changes (schema.go):

   timeout_client_data (line 83-90):
   - Added: Computed: true
   - Already had: Optional: true

   timeout_member_connect (line 91-98):
   - Added: Computed: true
   - Already had: Optional: true

   timeout_member_data (line 99-106):
   - Added: Computed: true
   - Already had: Optional: true

   sni_secret_id (line 112-118):
   - Added: Computed: true
   - Added: CustomType: customfield.NewListType[types.String](ctx)
   - Already had: Optional: true

   user_list (line 119-134):
   - Added: Computed: true
   - Added: CustomType: customfield.NewNestedObjectListType[CloudLoadBalancerListenerUserListModel](ctx)
   - Already had: Optional: true

2. Model Changes (model.go):

   TimeoutClientData (line 22):
   - Changed JSON tag: optional -> computed_optional
   - Type unchanged: types.Int64

   TimeoutMemberConnect (line 23):
   - Changed JSON tag: optional -> computed_optional
   - Type unchanged: types.Int64

   TimeoutMemberData (line 24):
   - Changed JSON tag: optional -> computed_optional
   - Type unchanged: types.Int64

   SniSecretID (line 26):
   - Changed JSON tag: optional -> computed_optional
   - Changed type: *[]types.String -> customfield.List[types.String]

   UserList (line 27):
   - Changed JSON tag: optional -> computed_optional
   - Changed type: *[]*CloudLoadBalancerListenerUserListModel -> customfield.NestedObjectList[CloudLoadBalancerListenerUserListModel]

TERRAFORM PLAN VERIFICATION

All 5 fields now correctly show (known after apply):
  timeout_client_data    = (known after apply)
  timeout_member_connect = (known after apply)
  timeout_member_data    = (known after apply)
  sni_secret_id          = (known after apply)
  user_list              = (known after apply)

This means Terraform understands these fields can be computed by the API.

COMPARISON WITH OLD PROVIDER

Old provider (old_terraform_provider/gcore/resource_gcore_lblistener.go):
- All fields were Optional: true only (no Computed)
- Used SDK-based schema (schema.TypeInt, schema.TypeList)
- Used GetOkExists to detect unset values
- No computed_optional concept

New provider:
- Uses Terraform Plugin Framework
- Requires Computed: true for API-computed fields
- Uses customfield types for lists that can be unknown
- Uses computed_optional JSON tags for proper unmarshaling

WHY THE OLD PROVIDER DIDN'T HAVE THIS ISSUE

The old provider used:
  d.GetOkExists("timeout_client_data")

This could distinguish between:
- Field not set in config (don't send to API)
- Field set to 0 in config (send 0 to API)

The new framework doesn't have GetOkExists. Instead it uses:
- Optional only: Always send the value (even if null)
- Computed + Optional: Only send if explicitly set in config

NEXT STEPS TO TEST

1. Source credentials:
   source ../.env

2. Run first apply:
   terraform apply -auto-approve

3. Run second apply to verify NO drift:
   terraform apply

Expected result: No changes detected

WHEN STAINLESS REGENERATES

These manual changes will be replaced by properly generated code that includes:
- Correct schema with Computed: true and CustomType
- Correct model with customfield types and computed_optional tags
- Proper plan modifiers (may include RequiresReplaceIfConfigured)

The manual changes match what Stainless should generate.
