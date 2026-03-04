# Security Group Rule PUT Endpoint - ID Change Issue

## Summary

The Gcore Cloud API's PUT endpoint for updating security group rules (`PUT /cloud/v1/securitygrouprules/{project_id}/{region_id}/{rule_id}`) **replaces the rule with a new one** (delete old + create new with different ID) instead of updating it in-place as documented.

## Test Date

2025-11-13

## API Documentation Claims

According to https://gcore.com/docs/api-reference/cloud/security-groups/update-security-group-rule:

- **Endpoint**: `PUT /cloud/v1/securitygrouprules/{project_id}/{region_id}/{rule_id}`
- **Expected Behavior**: "Update the configuration of an existing security group rule"
- **Expected Response**: Same rule with updated `updated_at` timestamp and incremented `revision_number`
- **Expected ID Behavior**: Rule ID should remain unchanged (in-place modification)

## Actual Behavior (Tested with curl)

### Test Setup
- **Security Group ID**: `5969403e-c10f-4864-9085-70df4ba18af7`
- **Original Rule ID**: `0f028d9f-7eba-4697-a0d8-3351b34bcd54`
- **Original Description**: "HTTP web traffic - TEST UPDATE 2"
- **Original Revision**: 0

### PUT Request
```bash
curl -X PUT \
  -H "Authorization: APIKey $GCORE_API_KEY" \
  -H "Content-Type: application/json" \
  "https://api.gcore.com/cloud/v1/securitygrouprules/$GCORE_CLOUD_PROJECT_ID/$GCORE_CLOUD_REGION_ID/$RULE_ID" \
  -d '{
    "security_group_id": "'$SG_ID'",
    "direction": "ingress",
    "ethertype": "IPv4",
    "protocol": "tcp",
    "port_range_min": 80,
    "port_range_max": 80,
    "remote_ip_prefix": "0.0.0.0/0",
    "description": "UPDATED VIA CURL PUT"
  }'
```

### PUT Response
```json
{
  "id": "9156daf0-0f63-477c-8ca0-c9397bb70d62",  ← NEW ID!
  "description": "UPDATED VIA CURL PUT",
  "direction": "ingress",
  "protocol": "tcp",
  "port_range_min": 80,
  "port_range_max": 80,
  "security_group_id": "5969403e-c10f-4864-9085-70df4ba18af7",
  "revision_number": 0,
  "updated_at": "2025-11-13T09:08:28+0000"
}
```

### Verification via Security Group GET

After the PUT request:
- **Original rule** (`0f028d9f-7eba-4697-a0d8-3351b34bcd54`): ✗ **NOT FOUND** in security group
- **New rule** (`9156daf0-0f63-477c-8ca0-c9397bb70d62`): ✓ **FOUND** in security group with updated description

## Key Findings

| Aspect | Expected (per docs) | Actual Behavior |
|--------|-------------------|----------------|
| **Rule ID** | Unchanged | **CHANGED** from `0f028d9f...` to `9156daf0...` |
| **Operation Type** | In-place update | **Replace** (delete old + create new) |
| **Original Rule** | Should remain with updated fields | **Deleted** from security group |
| **Revision Number** | Should increment | Unchanged (both 0) |

## Impact on Terraform Provider

### Problem
When Terraform calls the PUT endpoint to update a rule attribute (e.g., description), the API:
1. Deletes the old rule
2. Creates a new rule with a different ID
3. Returns the new ID in the response

This causes:
- **State inconsistency**: Terraform expects the resource ID to remain stable
- **Error**: "Provider produced inconsistent result after apply" because ID changed
- **Drift**: On next `terraform plan`, the old ID is not found and Terraform wants to recreate it

### Current Workaround

In `/Users/user/repos/gcore-terraform/internal/services/cloud_security_group_rule/resource.go:158-159`:

```go
// Preserve the original ID from state - API may return a different ID but the resource identity must not change
data.ID = state.ID
```

We manually override the returned ID with the original ID from state. This prevents Terraform from detecting the ID change, maintaining resource identity in state.

### Limitations of Workaround

1. **Read() method will fail on next refresh**: The old ID no longer exists in the API, so subsequent `terraform plan` or `terraform refresh` will fail to find the rule
2. **State becomes stale**: The state contains an ID that doesn't exist in the API
3. **Import/drift detection broken**: Cannot reliably import or detect drift for rules that have been updated

## Solution Implemented ✅

### Option 1: Use RequiresReplace for all fields (IMPLEMENTED)

We marked all mutable rule attributes with `RequiresReplace` so any change triggers a full delete + create cycle, matching the API's actual behavior.

**Implementation Changes:**

1. **schema.go** - Added `RequiresReplace` plan modifier to all optional fields:
   - `description`
   - `direction`
   - `ethertype`
   - `port_range_max`
   - `port_range_min`
   - `protocol`
   - `remote_group_id`
   - `remote_ip_prefix`

2. **resource.go** - Simplified Update() method to return an error since it should never be called

**Result:**
- ✅ Any rule attribute change triggers: destroy old rule → create new rule
- ✅ Matches the API's actual replace behavior
- ✅ State stays consistent - no stale IDs
- ✅ No drift detection issues
- ✅ Users see clear "forces replacement" messages in terraform plan

**Trade-off:** Changes to rules are disruptive (brief recreation), but this accurately reflects what the API does anyway.

## Alternative Options (Not Implemented)

### Option 2: Report API Bug to Gcore
The PUT endpoint should update in-place as documented, not replace with a new ID.

**Pros**: Would fix the root cause
**Cons**: Requires Gcore API team to fix their endpoint

### Option 3: Track by Content, Not ID (Complex)
Use a content-based key (security_group_id + direction + protocol + ports + CIDR) to identify rules instead of relying on API-provided IDs.

**Pros**: Works around API limitations
**Cons**: Very complex, requires major provider refactoring

## Test Script

A complete test script is available at:
`/Users/user/repos/gcore-terraform/test-secgroup-manual/test_rule_update_api.sh`

Run with:
```bash
cd /Users/user/repos/gcore-terraform/test-secgroup-manual
./test_rule_update_api.sh
```

## Test Results Files

- `/tmp/rule_before.json` - Rule state before PUT request
- `/tmp/rule_put_response.json` - PUT endpoint response showing new ID
- `/tmp/sg_after_update.json` - Security group state after update (original ID gone, new ID present)

## Conclusion

The Gcore Cloud API's PUT endpoint for security group rules does **NOT** update rules in-place as documented. Instead, it performs a **replace operation** (delete + create) that results in a new rule ID. This contradicts the API documentation and breaks Terraform's resource identity model.

The current Terraform provider workaround (preserving the old ID in state) prevents immediate errors but causes subsequent Read() operations to fail since the old ID no longer exists in the API.

**Action Required**: Either the API behavior needs to be fixed by Gcore, or the Terraform provider needs to explicitly treat all rule updates as replace operations.
