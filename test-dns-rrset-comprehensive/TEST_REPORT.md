# Test Report: gcore_dns_zone_rrset

**Date**: 2025-12-01
**Resource**: `gcore_dns_zone_rrset`
**Branch**: `tf-dns-zone-record`
**JIRA Ticket**: GCLOUD2-22174

## Executive Summary

| Category | Status | Details |
|----------|--------|---------|
| Create | ⚠️ Partial | Resource created but computed fields fail |
| Read | ❌ FAIL | `updated_at` field not parsed correctly |
| Update | ❌ FAIL | Not implemented (empty method) |
| Delete | ✅ PASS | Works correctly |
| Drift | ❌ FAIL | Multiple fields missing after refresh |
| Import | ❌ NOT TESTED | Not implemented |

## Critical Issues Found

### Issue 1: `updated_at` Field Parsing Error (BLOCKER)

**Severity**: 🔴 Critical
**Location**: `model.go:24`, `resource.go:86`

**Problem**: The API returns `updated_at` as Unix nanosecond timestamp (`1764589445883638000`), but the model expects RFC3339 format (`timetypes.RFC3339`).

**API Response**:
```json
{
  "updated_at": 1764589445883638000
}
```

**Error**:
```
Error: Provider returned invalid result object after apply
After the apply operation, the provider still indicated an unknown value
for gcore_dns_zone_rrset.test_a.updated_at.
```

**Fix Required**: Update the model to use `types.Int64` or add custom deserialization for timestamp conversion.

---

### Issue 2: Update Method Not Implemented (BLOCKER)

**Severity**: 🔴 Critical
**Location**: `resource.go:95-97`

**Problem**: The `Update` method is a stub that does nothing:
```go
func (r *DNSZoneRrsetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    // Update is not supported for this resource
}
```

**Impact**: All attribute changes force resource recreation (destroy + create).

**Evidence**:
```
# gcore_dns_zone_rrset.test_a must be replaced
~ ttl = 300 -> 600 # forces replacement
```

**Fix Required**: Implement Update using `r.client.DNS.Zones.Rrsets.Replace()` SDK method.

---

### Issue 3: All Attributes Have RequiresReplace (BLOCKER)

**Severity**: 🔴 Critical
**Location**: `schema.go`

**Problem**: Every attribute has `planmodifier.RequiresReplace()`:
- `ttl` - Should support in-place update
- `resource_records` - Should support in-place update
- `pickers` - Should support in-place update
- `meta` - Should support in-place update

**Fix Required**: Remove `RequiresReplace()` from all updatable fields. Only `zone_name`, `rrset_name`, and `rrset_type` should force replacement (they're part of the resource identity).

---

### Issue 4: Missing `id` Field in ResourceRecords Model

**Severity**: 🟡 Medium
**Location**: `model.go:37-41`

**Problem**: API returns `id` for each resource record, but model doesn't capture it:

**API Response**:
```json
{
  "resource_records": [
    {
      "id": 143174902,
      "content": ["192.168.1.100"],
      "enabled": true,
      "meta": {}
    }
  ]
}
```

**Current Model**:
```go
type DNSZoneRrsetResourceRecordsModel struct {
    Content *[]jsontypes.Normalized `tfsdk:"content" json:"content,required"`
    Enabled types.Bool              `tfsdk:"enabled" json:"enabled,computed_optional"`
    Meta    *map[string]jsontypes.Normalized `tfsdk:"meta" json:"meta,optional"`
    // ID field missing!
}
```

**Fix Required**: Add `ID types.Int64 \`tfsdk:"id" json:"id,computed"\`` to the model and schema.

---

### Issue 5: `content` Field Uses `jsontypes.Normalized` (UX Issue)

**Severity**: 🟡 Medium
**Location**: `schema.go:45-49`

**Problem**: Users must use `jsonencode()` for content values, which is not intuitive:
```hcl
resource_records = [
  {
    content = [jsonencode("192.168.1.100")]  # Required!
    # content = ["192.168.1.100"]  # This fails!
  }
]
```

**Impact**: Poor developer experience, confusing error messages.

**Recommendation**: Consider using `types.String` elements with custom validation, or improve documentation.

---

### Issue 6: Missing ImportState Implementation

**Severity**: 🟡 Medium
**Location**: `resource.go`

**Problem**: No `ImportState` method. Users cannot import existing DNS records.

**Old Provider Import Format**: `zone:domain:type` (e.g., `maxima.lt:www.maxima.lt:A`)

**Fix Required**: Implement `ImportState` with format `zone_name/rrset_name/rrset_type`.

---

### Issue 7: Read Uses `apijson.Unmarshal` Instead of `UnmarshalComputed`

**Severity**: 🟡 Medium
**Location**: `resource.go:129`

**Problem**: Read operation uses `apijson.Unmarshal` which may not handle computed fields correctly:
```go
err = apijson.Unmarshal(bytes, &data)  // Should be UnmarshalComputed
```

**Fix Required**: Change to `apijson.UnmarshalComputed(bytes, &data)`.

---

### Issue 8: `filter_set_id` Null Handling

**Severity**: 🟢 Low
**Location**: `model.go:21`

**Problem**: API returns `"filter_set_id": null` but model uses `types.Int64`. May cause issues.

**API Response**: `"filter_set_id": null`

**Impact**: Field not appearing in state.

---

## Test Results

### Test 1: Create Basic A Record

| Step | Result | Notes |
|------|--------|-------|
| Plan | ✅ Pass | Correct plan generated |
| Apply | ⚠️ Partial | Resource created but error on computed fields |
| API Verification | ✅ Pass | Record exists in DNS API |
| State | ⚠️ Partial | Missing `updated_at`, `filter_set_id`, `warnings` |

### Test 2: Drift Detection

| Step | Result | Notes |
|------|--------|-------|
| Refresh | ⚠️ Partial | Some fields refreshed, `updated_at` missing |
| Second Plan | ❌ Fail | Exit code 2 (changes detected) |

### Test 3: Update TTL

| Step | Result | Notes |
|------|--------|-------|
| Plan | ❌ Fail | Shows "forces replacement" instead of update |
| Root Cause | `RequiresReplace()` on `ttl` field |

### Test 4: Delete

| Step | Result | Notes |
|------|--------|-------|
| Destroy | ✅ Pass | Resource deleted correctly |
| API Verification | ✅ Pass | Record no longer exists |

## API Calls Observed

```
POST /dns/v2/zones/maxima.lt/tf-test-a.maxima.lt/A  # Create
GET /dns/v2/zones/maxima.lt/tf-test-a.maxima.lt/A   # Read
DELETE /dns/v2/zones/maxima.lt/tf-test-a.maxima.lt/A # Delete
# PUT not called (Update not implemented)
```

## Recommended Fixes (Priority Order)

1. **[P0]** Fix `updated_at` timestamp parsing (convert Unix nanoseconds to RFC3339)
2. **[P0]** Implement `Update` method using `Rrsets.Replace()`
3. **[P0]** Remove `RequiresReplace()` from updatable fields (`ttl`, `resource_records`, `pickers`, `meta`)
4. **[P1]** Add `id` field to `ResourceRecordsModel`
5. **[P1]** Implement `ImportState` method
6. **[P1]** Fix Read to use `UnmarshalComputed`
7. **[P2]** Improve `content` field UX (documentation or type change)
8. **[P2]** Handle `filter_set_id` null properly

## Files to Modify

| File | Changes Required |
|------|------------------|
| `model.go` | Fix `updated_at` type, add `id` to ResourceRecords |
| `resource.go` | Implement Update, ImportState, fix Read |
| `schema.go` | Remove `RequiresReplace` from updatable fields, add `id` to schema |

## Evidence Artifacts

- `test-dns-rrset-comprehensive/main.tf` - Test configuration
- `terraform.tfstate` - State file (if preserved)
- API response showing correct data structure

## Conclusion

The `gcore_dns_zone_rrset` resource has **critical implementation gaps** that prevent basic functionality:

1. **Cannot update resources** - All changes force recreation
2. **Cannot track all state** - `updated_at` parsing failure
3. **No import support** - Cannot manage existing records

**Recommendation**: Fix P0 issues before releasing this resource.
