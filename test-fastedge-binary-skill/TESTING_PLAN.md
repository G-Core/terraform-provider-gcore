# Testing Plan: gcore_fastedge_binary

## JIRA Ticket: GCLOUD2-22876

## Phase 1: Analysis Summary

### Resource Overview

**Resource Type**: `gcore_fastedge_binary`
**Purpose**: Upload and manage FastEdge WebAssembly binaries
**Immutability**: Binaries are immutable - any content change triggers replacement

### OpenAPI → SDK Alignment

| API Operation | SDK Method | Async? | Notes |
|---------------|------------|--------|-------|
| POST /fastedge/v1/binaries/raw | `Binaries.New()` | No | Raw binary upload (application/octet-stream) |
| GET /fastedge/v1/binaries | `Binaries.List()` | No | List all binaries |
| GET /fastedge/v1/binaries/{id} | `Binaries.Get()` | No | Get binary details |
| DELETE /fastedge/v1/binaries/{id} | `Binaries.Delete()` | No | Delete binary |

**Notes**:
- No UPDATE endpoint - binaries are immutable
- Create operation uses raw binary upload (not JSON)
- No async operations (no tasks/polling)

### Legacy Provider Comparison

| Feature | Old Provider | New Provider | Parity? |
|---------|-------------|--------------|---------|
| Create | Direct API call | `Binaries.New()` | ✅ |
| Read | Direct API call | `Binaries.Get()` | ✅ |
| Delete | Direct API call | `Binaries.Delete()` | ✅ |
| Update | N/A (immutable) | N/A (immutable) | ✅ |
| Checksum calculation | CustomizeDiff | ModifyPlan | ✅ |
| 404 on Read | Warning + remove | Warning + remove | ✅ |
| 409 on Delete | Warning (in use) | Warning (in use) | ✅ |
| ForceNew/Replace | `ForceNew: true` | `RequiresReplace()` | ✅ |

### Schema Comparison

| Field | Old Provider | New Provider | Match? |
|-------|-------------|--------------|--------|
| `id` | Computed | Computed, RequiresReplace, UseStateForUnknown | ✅ |
| `filename` | Required, ForceNew | Required, RequiresReplace | ✅ |
| `checksum` | Computed, ForceNew | Computed, RequiresReplace | ✅ |
| `api_type` | - | Computed | New |
| `source` | - | Computed | New |
| `status` | - | Computed | New |
| `unref_since` | - | Computed | New |

### Migration Nuances

- [x] No special endpoints (unset, attach, detach) - resource is simple
- [x] File path handling preserved (local file → binary upload)
- [x] Checksum verification preserved
- [x] Error handling patterns matched (404 warning, 409 graceful)

### JIRA/Known Issues

- No known issues specific to this resource
- File upload is a special pattern - uses `application/octet-stream`

---

## Phase 2: Resource Coverage Matrix

### Union Type Coverage

**Not applicable** - This resource has no union types. Simple scalar fields only.

### Field Coverage Matrix

| Field | Type | Tag | Default | Test # | Notes |
|-------|------|-----|---------|--------|-------|
| `id` | Int64 | computed | API-assigned | All | UseStateForUnknown + RequiresReplace |
| `filename` | String | required | - | 1,3,4 | RequiresReplace on change |
| `checksum` | String | computed | MD5 hash | 1,2,3 | RequiresReplace, calculated in ModifyPlan |
| `api_type` | String | computed | "reactr" | 1 | API determines WASM type |
| `source` | Int64 | computed | 2 (JS) | 1 | 0=unknown, 1=Rust, 2=JavaScript |
| `status` | Int64 | computed | 1 (compiled) | 1 | Compilation status |
| `unref_since` | String | computed | null | 1 | UTC timestamp if unreferenced |

### Test Case Matrix

| Test # | Description | Focus | Expected Behavior |
|--------|-------------|-------|-------------------|
| 1 | Create + Drift Test | Basic lifecycle | No drift on 2nd plan |
| 2 | File Change Detection | RequiresReplace | Replacement triggered when file changes |
| 3 | Filename Change | RequiresReplace | Replacement triggered when path changes |
| 4 | Import Test | State import | Import works, warns about filename |
| 5 | Delete with 409 | Error handling | Warning when binary in use |
| 6 | Delete normal | Clean deletion | Resource removed successfully |

### QA Checklist Mapping

| QA Category | Test # | Expected Behavior |
|-------------|--------|-------------------|
| Drift detection | 1 | No changes on 2nd plan |
| RequiresReplace | 2,3 | Forces replacement on content/path change |
| Computed fields | 1 | All computed fields populated correctly |
| Import | 4 | State imported, filename warning shown |
| Error handling | 5 | 409 Conflict handled gracefully |
| Delete | 6 | Clean deletion works |

---

## Phase 3: Test Execution Plan

| Test # | Description | Config | Verification | Dependencies |
|--------|-------------|--------|--------------|--------------|
| 1 | Create + Drift | minimal.wasm | `plan -detailed-exitcode` returns 0 | None |
| 2 | File change | modified.wasm | Plan shows "forces replacement" | Test 1 |
| 3 | Path change | same content, different path | Plan shows "forces replacement" | Test 1 |
| 4 | Import | Import existing ID | Warning about filename | Existing binary |
| 5 | Delete conflict | Binary in use by app | Warning shown | FastEdge app |
| 6 | Clean delete | Normal destroy | Resource removed | Test 1 |

### Test Files Required

1. `minimal.wasm` - Valid WASM binary compiled with @gcoredev/fastedge-sdk-js
2. `modified.wasm` - Same structure, different content (for change detection)
3. `minimal_copy.wasm` - Copy of minimal.wasm (for path change test)

### Dependencies

- **Network**: No
- **Subnet**: No
- **Pre-created resources**: No (self-contained)
- **FastEdge app**: Only for Test 5 (optional)

### Cleanup Plan

Resources to destroy:
- All test binaries (automatic via terraform destroy)
- API cleanup via Gcore MCP if needed

---

## Reviewer Sign-off

- [x] All fields mapped to tests
- [x] All error cases covered
- [x] RequiresReplace behavior tested
- [x] Import test included
- [x] Dependencies identified
- [x] Cleanup plan reviewed

**Resource has no union types - simple file upload pattern**

---

## Phase 4: Execution Notes

### Pre-requisites Completed

1. ✅ Provider built: `go build -o terraform-provider-gcore`
2. ✅ .terraformrc configured for dev overrides
3. ✅ Credentials loaded from .env
4. ✅ Test WASM binary available: `minimal.wasm`

### Previous Test Results (from earlier session)

Tests already performed successfully:
- Create: Binary ID 2848 created
- Drift: No changes on 2nd plan
- File change detection: RequiresReplace triggered correctly
- Delete: Clean deletion worked
- Import: State imported with filename warning

### Evidence from jira-evidence/

- `new_main.tf` - New provider configuration
- `new_terraform.tfstate` - State file from new provider
- `old_main.tf` - Old provider configuration
- `old_terraform.tfstate` - State file from old provider (Binary ID 2849)

---

## Test Execution

### Test 1: Create + Drift Detection

```bash
cd test-fastedge-binary-skill
terraform init
terraform apply -auto-approve
terraform plan -detailed-exitcode
# Expected: Exit code 0 (no changes)
```

### Test 2: File Change Detection

```bash
# Modify minimal.wasm or use different file
# Edit main.tf to point to modified file
terraform plan
# Expected: Shows "forces replacement" for checksum change
```

### Test 3: Import Test

```bash
# After creating a binary, get the ID
# Remove from state
terraform state rm gcore_fastedge_binary.test

# Import
terraform import gcore_fastedge_binary.test <binary_id>
# Expected: Warning about filename required
```

### Test 4: Clean Delete

```bash
terraform destroy -auto-approve
# Expected: Resource deleted successfully
```
