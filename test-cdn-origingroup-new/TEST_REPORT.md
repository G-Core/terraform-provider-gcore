# Test Report: gcore_cdn_origin_group (GCLOUD2-22597)

## Summary

Analyzed the `gcore_cdn_origin_group` resource in the Stainless-generated provider and compared it with the old `gcore_cdn_origingroup` implementation. Identified and fixed three critical issues.

**Test Status**: Blocked by CDN service access (403 error). Fixes implemented based on code analysis.

## Code Analysis

### Old Provider (`gcore_cdn_origingroup`)
- **File**: `old_terraform_provider/gcore/resource_gcore_cdn_origin_group.go`
- **Key Features**:
  1. S3 credentials (`s3_access_key_id`, `s3_secret_access_key`) marked as Sensitive
  2. Credentials preserved in Read() operation - API doesn't return them
  3. CustomizeDiff validation for mutual exclusivity between `origin` and `auth`

### New Provider (`gcore_cdn_origin_group`)
- **Directory**: `internal/services/cdn_origin_group/`
- **Initial Issues**:
  1. S3 credentials NOT marked as Sensitive
  2. S3 credentials NOT preserved in Read() - would be lost on refresh
  3. ConfigValidators was empty - no validation

## Issues Found & Fixes Applied

### Issue 1: S3 Credentials Not Marked Sensitive

**Problem**: `s3_access_key_id` and `s3_secret_access_key` were not marked as sensitive, causing them to appear in plan output and logs.

**Fix**: Added `Sensitive: true` to both attributes in `schema.go`:
```go
"s3_access_key_id": schema.StringAttribute{
    // ...
    Sensitive: true,
},
"s3_secret_access_key": schema.StringAttribute{
    // ...
    Sensitive: true,
},
```

### Issue 2: S3 Credentials Lost on Read

**Problem**: The CDN API does NOT return `s3_access_key_id` and `s3_secret_access_key` in GET responses. The new provider would overwrite the state with the API response, losing the credentials.

**Fix**: Modified Create(), Update(), and Read() functions in `resource.go` to:
1. Save credentials before API call
2. Restore credentials after unmarshalling API response

```go
// Preserve S3 credentials from state - API doesn't return them
var savedS3AccessKeyID, savedS3SecretAccessKey types.String
if data.Auth != nil {
    savedS3AccessKeyID = data.Auth.S3AccessKeyID
    savedS3SecretAccessKey = data.Auth.S3SecretAccessKey
}

// ... API call ...

// Restore S3 credentials from prior state since API doesn't return them
if data.Auth != nil && !savedS3AccessKeyID.IsNull() {
    data.Auth.S3AccessKeyID = savedS3AccessKeyID
    data.Auth.S3SecretAccessKey = savedS3SecretAccessKey
}
```

### Issue 3: Missing Validation

**Problem**: ConfigValidators returned an empty slice. Users could specify both `sources` and `auth`, or neither.

**Fix**: Added ExactlyOneOf validator in `schema.go`:
```go
func (r *CdnOriginGroupResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
    return []resource.ConfigValidator{
        resourcevalidator.ExactlyOneOf(
            path.MatchRoot("sources"),
            path.MatchRoot("auth"),
        ),
    }
}
```

## Files Modified

1. `internal/services/cdn_origin_group/schema.go`
   - Added `Sensitive: true` to S3 credential fields
   - Added imports for `resourcevalidator` and `path`
   - Implemented ConfigValidators with ExactlyOneOf validator

2. `internal/services/cdn_origin_group/resource.go`
   - Modified Create() to preserve S3 credentials
   - Modified Update() to preserve S3 credentials
   - Modified Read() to preserve S3 credentials

## Build Verification

- Go build: PASSED
- Lint script: PASSED

## Testing Status

### CDN Access Investigation (2026-01-07)

CDN service access was requested and partially enabled. Investigation revealed:

**Read access**: WORKS
```bash
# GET /cdn/origin_groups returns successfully
curl -H "Authorization: APIKey $GCORE_API_KEY" "https://api.gcore.com/cdn/origin_groups"
# Returns: []
```

**Write access**: BLOCKED (403)
```bash
# POST /cdn/origin_groups fails
curl -X POST -H "Authorization: APIKey $GCORE_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"name": "test", "sources": [{"source": "example.com", "enabled": true}]}' \
  "https://api.gcore.com/cdn/origin_groups"
# Returns: {"message":"CDN service is stopped. You do not have permission to perform this action."}
```

**Conclusion**: CDN list/read operations are enabled, but create/update/delete operations require additional permissions. Full CDN write access needs to be requested to complete infrastructure testing.

### Next Steps to Complete Testing

1. Request full CDN write access (create/update/delete permissions)
2. Once granted, run comprehensive tests:
   - Create with public origins
   - Create with S3 auth (both amazon and other types)
   - Update operations
   - Import operation
   - Drift detection (verify credentials preserved)
   - Delete operations

## Test Configurations Created

Test configurations created in `test-cdn-origingroup-new/`:
- `main.tf` - Test cases for public origins and Gcore S3 auth

## Artifacts

- Branch: `fix/cdn-origin-group-GCLOUD2-22597`
- Commit: `fix(cdn_origin_group): preserve S3 credentials and add validation`
- Test config: `test-cdn-origingroup-new/main.tf`
- Old provider example: `test-cdn-origingroup-old/main.tf`

## Recommendations

1. **Testing**: Once CDN access is restored, run comprehensive tests covering:
   - Create with public origins
   - Create with S3 auth (both amazon and other types)
   - Update operations
   - Import operation
   - Drift detection (verify credentials preserved)

2. **Additional Validation**: Consider adding validators for:
   - If `auth.s3_type = "amazon"`, require `s3_region`
   - If `auth.s3_type = "other"`, require `s3_storage_hostname`
