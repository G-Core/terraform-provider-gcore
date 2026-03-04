# Test Report: gcore_dns_zone Resource (GCLOUD2-22173)

## Summary

Testing of the `gcore_dns_zone` and `gcore_dns_zone_rrset` resources for the new Stainless-generated Terraform provider, with comparison to the old G-Core provider.

**Test Date**: 2025-12-08
**JIRA Ticket**: GCLOUD2-22173
**Related PR**: stainless-sdks/gcore-terraform#60 (dns_zone_rrset improvements)

## Key Findings

### 1. Critical Bug Found and Fixed

**Issue**: `dns_zone` resource `Read` and `Delete` methods used integer `ID` but DNS API expects zone **name** (string).

**Location**: `internal/services/dns_zone/resource.go:111,145`

**Fix Applied**:
```go
// Before (broken):
r.client.DNS.Zones.Get(ctx, strconv.FormatInt(data.ID.ValueInt64(), 10), ...)

// After (fixed):
r.client.DNS.Zones.Get(ctx, data.Name.ValueString(), ...)
```

### 2. DNS Operations Are Synchronous

**Verified**: DNS Zone SDK methods (`New`, `Get`, `Delete`, `Replace`) do **NOT** use task polling. They are synchronous operations - no `*AndPoll` methods required.

### 3. Compatibility Issues Between Old and New Providers

| Feature | Old Provider (G-Core/gcore) | New Provider (gcore/gcore) | Compatibility |
|---------|---------------------------|---------------------------|---------------|
| DNSSEC toggle | `dnssec = true/false` | Not supported | :x: Breaking |
| Zone enable/disable | `enabled` with special API | `enabled` forces replacement | :warning: Different |
| Update SOA fields | In-place update | All fields force replacement | :x: Breaking |
| Zone ID | Uses zone name as ID | Uses integer ID | :warning: Different |
| RRSet Update | Full CRUD support | PR#60 adds Replace API | :white_check_mark: With PR#60 |

### 4. dns_zone_rrset Issues (Fixed in PR#60)

The current `main` branch has these issues that PR#60 fixes:

1. **`updated_at` field** - Causes "unknown value after apply" error because API doesn't return this field
2. **No Update support** - Cannot modify TTL/records without recreate
3. **No Import support** - Cannot import existing RRSets

## Test Execution Results

### Old Provider Tests

| Test | Result | Notes |
|------|--------|-------|
| Create zone | :white_check_mark: PASS | Zone created successfully |
| Create A record | :white_check_mark: PASS | |
| Create TXT record | :white_check_mark: PASS | |
| Create CNAME record | :white_check_mark: PASS | |
| Create MX record | :white_check_mark: PASS | |
| Drift detection | :warning: PARTIAL | API returns computed defaults (contact, serial, primary_server) not in config |
| Zone update | :x: BLOCKED | Tariff restriction: "contact editing is prohibited" |
| Destroy | :white_check_mark: PASS | All resources cleaned up |

### New Provider Tests

| Test | Result | Notes |
|------|--------|-------|
| Create zone | :white_check_mark: PASS | Zone ID=983629 |
| Create A record | :white_check_mark: PASS | |
| Create TXT record | :white_check_mark: PASS | |
| Create CNAME record | :white_check_mark: PASS | |
| Create MX record | :white_check_mark: PASS | |
| Read after create | :x: FAIL (before fix) | "Resource not found" due to ID vs Name bug |
| Read after create | :white_check_mark: PASS (after fix) | |
| dns_zone_rrset apply | :x: FAIL | "updated_at" unknown value error |
| Destroy | :white_check_mark: PASS | All resources cleaned up |

## API Observations

### DNS Zone API Response Structure

```json
{
  "zone": null,
  "name": "tf-test-22173-new.com",
  "meta": null,
  "enabled": true,
  "nx_ttl": 300,
  "retry": 3600,
  "refresh": 0,
  "expiry": 1209600,
  "contact": "support@gcore.com",
  "serial": 1765185088,
  "primary_server": "ns1.gcorelabs.net",
  "records": [...],
  "status": "pending",
  "dnssec_enabled": false
}
```

**Key Notes**:
- `id` is NOT at the top level - it's inside `zone.id` (which is null in this response)
- API uses zone **name** for all operations, not integer ID
- `zone` object is null in GET response but may contain details in other contexts

## Recommendations

### Immediate Actions Required

1. **Merge PR#60** - Fixes dns_zone_rrset `updated_at` issue and adds Update/Import support
2. **Apply Read/Delete fix** - The fix for using Name instead of ID in `dns_zone/resource.go`

### Future Enhancements

1. **Add Update support to dns_zone** - Currently all fields trigger replacement
2. **Add DNSSEC support** - Old provider has this, new provider doesn't
3. **Review ID handling** - Consider using zone name as resource ID like old provider

## Test Artifacts

- `evidence/old_provider_state.json` - Terraform state from old provider
- `evidence/new_provider_state.json` - Terraform state from new provider
- `old_provider/main.tf` - Test configuration for old provider
- `new_provider/main.tf` - Test configuration for new provider

## Zones Created During Testing

| Zone Name | Provider | Status |
|-----------|----------|--------|
| tf-test-22173-gcloud.com | Old | Destroyed |
| tf-test-22173-new.com | New | Destroyed |
