# Schema Comparison: Old vs New Provider

## Resource Name
| Old Provider | New Provider |
|--------------|--------------|
| `gcore_dns_zone_record` | `gcore_dns_zone_rrset` |

## Import Format
| Old Provider | New Provider |
|--------------|--------------|
| `zone:domain:type` (colon-separated) | `zone_name/rrset_name/rrset_type` (slash-separated) |

Example:
- Old: `terraform import gcore_dns_zone_record.test maxima.lt:tf-test.maxima.lt:A`
- New: `terraform import gcore_dns_zone_rrset.test maxima.lt/tf-test.maxima.lt/A`

## Attribute Mapping

### Identity Fields
| Old Provider | New Provider |
|--------------|--------------|
| `zone` | `zone_name` |
| `domain` | `rrset_name` |
| `type` | `rrset_type` |

### Common Fields
| Old Provider | New Provider | Notes |
|--------------|--------------|-------|
| `ttl` | `ttl` | Same |
| `filter` | `pickers` | Renamed, same structure |
| `meta` (RRSet level) | `meta` | Map instead of Set |

### Resource Records
| Old Provider | New Provider | Notes |
|--------------|--------------|-------|
| `resource_record` (TypeSet) | `resource_records` (List) | Set → List |
| `content` (string) | `content` (list of JSON) | String → JSON array |
| `enabled` (bool) | `enabled` (bool) | Same |
| `meta` (nested set) | `meta` (map of JSON) | Structured → JSON map |
| N/A | `id` (computed) | New: API-assigned record ID |

### Computed Fields (New Provider Only)
| Field | Description |
|-------|-------------|
| `name` | Computed from API |
| `type` | Computed from API |
| `filter_set_id` | Filter set ID from API |
| `warning` | Deprecated, use `warnings` |
| `warnings` | List of warning objects |

## Content Format Differences

### A Record
```hcl
# Old Provider
resource_record {
  content = "192.168.1.100"
}

# New Provider
resource_records = [
  {
    content = ["\"192.168.1.100\""]
  }
]
```

### MX Record
```hcl
# Old Provider
resource_record {
  content = "10 mail.example.com."
}

# New Provider
resource_records = [
  {
    content = ["10", "\"mail.example.com.\""]
  }
]
```

### SRV Record
```hcl
# Old Provider
resource_record {
  content = "100 1 5061 sip.example.com."
}

# New Provider
resource_records = [
  {
    content = ["100", "1", "5061", "\"sip.example.com.\""]
  }
]
```

## Meta Format Differences

### Old Provider (Structured)
```hcl
resource_record {
  content = "192.168.1.100"
  meta {
    countries = ["US", "CA"]
    weight    = 50
  }
}
```

### New Provider (JSON)
```hcl
resource_records = [
  {
    content = ["\"192.168.1.100\""]
    meta = {
      countries = "[\"US\", \"CA\"]"
      weight    = "50"
    }
  }
]
```

## State Migration Notes

State migration from old to new provider requires:
1. Remove old resources from state
2. Update HCL to new schema
3. Import resources with new provider

There is no automatic state migration path.
