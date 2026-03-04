ROOT CAUSE ANALYSIS AND COMPLETE FIX

ISSUE
Second terraform apply detects drift even though nothing changed in configuration.

ROOT CAUSE IDENTIFIED
The listener Read method was using apijson.Unmarshal instead of apijson.UnmarshalComputed

File: internal/services/cloud_load_balancer_listener/resource.go
Line: 225 (in Read method)

WRONG CODE:
  err = apijson.Unmarshal(bytes, &data)

CORRECT CODE:
  err = apijson.UnmarshalComputed(bytes, &data)

WHY THIS CAUSES DRIFT

When using regular Unmarshal:
1. ALL fields from API response overwrite the state
2. Including computed_optional fields that were null in config
3. This creates a mismatch between config (null) and state (API values)
4. Second apply detects this as drift

When using UnmarshalComputed:
1. Checks JSON tags on model fields
2. For computed_optional fields: only updates if value was explicitly set
3. For optional fields: always updates
4. For computed fields: always updates
5. This respects the "don't send if not in config" behavior

COMPLETE FIX APPLIED

1. Schema Changes (schema.go):
   - timeout_client_data: Added Computed: true
   - timeout_member_connect: Added Computed: true
   - timeout_member_data: Added Computed: true
   - sni_secret_id: Added Computed: true + CustomType
   - user_list: Added Computed: true + CustomType

2. Model Changes (model.go):
   - TimeoutClientData: Changed to computed_optional tag
   - TimeoutMemberConnect: Changed to computed_optional tag
   - TimeoutMemberData: Changed to computed_optional tag
   - SniSecretID: Changed type to customfield.List + computed_optional tag
   - UserList: Changed type to customfield.NestedObjectList + computed_optional tag

3. Resource Method Fix (resource.go):
   - Read method line 225: Changed Unmarshal to UnmarshalComputed

WHY ALL THREE CHANGES ARE NEEDED

Schema Computed: true
- Tells Terraform these fields can be computed by API
- Shows (known after apply) in plan
- Prevents drift detection on known values

Model computed_optional tag
- Tells unmarshaler to use IfUnset behavior
- Only updates state if value was set in config
- Preserves null values when field not in config

UnmarshalComputed in Read
- Activates the computed_optional logic
- Respects the JSON tags
- Without this ALL changes above are ignored

COMPARISON WITH WORKING RESOURCES

Load Balancer resource (works correctly):
- Read method uses: apijson.UnmarshalComputed (line 223)
- Has computed_optional fields working properly

Load Balancer Pool resource (works correctly):
- Read method uses: apijson.UnmarshalComputed
- Has computed_optional fields working properly

Load Balancer Listener (was broken, now fixed):
- Read method was using: apijson.Unmarshal (WRONG)
- Changed to: apijson.UnmarshalComputed (CORRECT)

HOW TO VERIFY FIX

1. Source credentials:
   source ../.env

2. First apply (creates resources):
   terraform apply -auto-approve

3. Second apply (should show NO changes):
   terraform apply

Expected output:
  No changes. Your infrastructure matches the configuration.

If drift is still detected, check:
- Is Read method using UnmarshalComputed?
- Are all fields marked as Computed: true in schema?
- Are all JSON tags set to computed_optional in model?
- Are list types using customfield types?

LESSONS LEARNED

1. Computed + Optional fields require THREE coordinated changes:
   - Schema: Computed: true, Optional: true
   - Model: computed_optional JSON tag
   - Resource Read: apijson.UnmarshalComputed

2. List fields need special types:
   - Simple lists: customfield.List[T]
   - Nested lists: customfield.NestedObjectList[T]
   - Cannot use *[]T for computed_optional

3. Always check Read method implementation:
   - Must use UnmarshalComputed for resources with computed_optional fields
   - Regular Unmarshal ignores JSON tags and overwrites everything

4. The old SDK provider didn't have this issue because:
   - Used GetOkExists which could detect "not set"
   - Framework requires explicit Computed: true

WHEN STAINLESS REGENERATES

All these manual changes should match what Stainless generates from the x-stainless-terraform-configurability: computed_optional attributes in OpenAPI spec.

The Read method fix (UnmarshalComputed) should also be in Stainless template and generated correctly.
