# attachment_tag Comprehensive Test Report

## Test Date: 2025-12-03

## Summary
All tests passed for the `attachment_tag` field implementation in the `gcore_cloud_instance` volumes schema.

## Test Results

| Test # | Description | Result | Evidence |
|--------|-------------|--------|----------|
| 1 | Create instance WITH attachment_tag="boot" | PASS | API payload includes `"attachment_tag":"boot"` |
| 2 | Drift detection (terraform plan -detailed-exitcode) | PASS | Exit code 0 - no drift |
| 3 | Create instance WITHOUT attachment_tag | PASS | API payload correctly omits `attachment_tag` field |
| 4 | Cleanup - destroy all resources | PASS | All resources destroyed |

## API Payload Evidence

### Test 1: With attachment_tag
```json
{
  "volumes": [{
    "attachment_tag": "boot",
    "boot_index": 0,
    "delete_on_termination": false,
    "source": "existing-volume",
    "volume_id": "6e700677-4080-4ea3-9fbd-02d98e851e2e"
  }]
}
```

### Test 3: Without attachment_tag
```json
{
  "volumes": [{
    "boot_index": 0,
    "delete_on_termination": false,
    "source": "existing-volume",
    "volume_id": "fca11b2e-a3ce-46b8-ac5e-f04cb74cbe2e"
  }]
}
```
Note: `attachment_tag` field correctly omitted when not specified in config.

## Implementation Details

### Files Modified:
1. `internal/services/cloud_instance/model.go` - Added `AttachmentTag` field and updated `MarshalJSONWithState`
2. `internal/services/cloud_instance/schema.go` - Added `attachment_tag` schema attribute
3. `examples/resources/gcore_cloud_instance/resource.tf` - Updated example

### Key Code (MarshalJSONWithState):
```go
func (m CloudInstanceVolumesModel) MarshalJSONWithState(plan any, state any) ([]byte, error) {
    payload := map[string]interface{}{
        "source":                "existing-volume",
        "volume_id":             m.VolumeID.ValueString(),
        "boot_index":            m.BootIndex.ValueInt64(),
        "delete_on_termination": false,
    }

    // Only include attachment_tag if it's set
    if !m.AttachmentTag.IsNull() && !m.AttachmentTag.IsUnknown() && m.AttachmentTag.ValueString() != "" {
        payload["attachment_tag"] = m.AttachmentTag.ValueString()
    }

    return json.Marshal(payload)
}
```

## Conclusion
The `attachment_tag` field implementation is complete and working correctly:
- Field is included in API payload when specified by user
- Field is omitted from API payload when not specified (no empty string sent)
- No drift detected after creation
- Resources can be created and destroyed successfully
