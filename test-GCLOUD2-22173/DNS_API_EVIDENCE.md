# DNS API Response Mismatch Evidence

## Issue
The `GET /dns/v2/zones/{name}` endpoint response structure does NOT match the OpenAPI specification.

## OpenAPI Spec (gcore-config/openapi.yml)

### Expected Response Structure (line 113218-113223)
```yaml
ItemZoneResponse:
  type: object
  properties:
    Zone:
      $ref: '#/components/schemas/OutputZone'
  description: Complete zone info with all records included
```

**Expected Response:**
```json
{
  "Zone": {
    "name": "...",
    "contact": "...",
    ... OutputZone fields ...
  }
}
```

### OutputZone Schema (line 113720-113782)
Does NOT include `enabled` field. Fields in OutputZone:
- client_id, contact, dnssec_enabled, expiry, id, meta, name, nx_ttl
- primary_server, records, refresh, retry, rrsets_amount, serial, status

## Actual API Response

Captured from `GET /dns/v2/zones/tf-test-gcloud2-22173.com`:

```json
{
  "name": "tf-test-gcloud2-22173.com",
  "meta": null,
  "enabled": true,
  "nx_ttl": 300,
  "retry": 3600,
  "refresh": 0,
  "expiry": 1209600,
  "contact": "support@gcore.com",
  "serial": 1765963470,
  "primary_server": "ns1.gcorelabs.net",
  "records": [...],
  "rrsets_amount": {...},
  "status": "non-delegated",
  "dnssec_enabled": false
}
```

## Discrepancies

| Aspect | OpenAPI Spec | Actual API |
|--------|--------------|------------|
| Response wrapper | `{ "Zone": { ... } }` | Flat JSON at root |
| `enabled` field | NOT in OutputZone | Present in response |

## Impact on Terraform Provider

1. **Data source `zone` attribute returns null** - Model expects `json:"Zone,computed"` but API doesn't return `Zone` wrapper
2. **Generated code mismatch** - Stainless generates code based on OpenAPI spec, but actual API differs

## Recommendation

DNS team should update OpenAPI spec to match actual API behavior:
1. Change `ItemZoneResponse` to return flat OutputZone (no Zone wrapper)
2. Add `enabled` field to OutputZone schema
