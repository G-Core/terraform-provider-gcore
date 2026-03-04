# WAAP Resources Implementation Plan

Based on analysis of:
- Jira GCLOUD2-22875
- Stainless build diagnostics
- Old terraform provider (`terraform-provider-gcore`)
- Go SDK (`gcore-go`)
- TERRAFORM_GUIDE.md patterns

---

## Implementation Overview

### Resources to Implement

| Resource | Complexity | Special Requirements |
|----------|------------|---------------------|
| gcore_waap_domain | High | Lookup-only create, multi-endpoint update |
| gcore_waap_domain_api_path | Medium | UUID IDs, status management |
| gcore_waap_domain_custom_rule | High | 18+ condition types, complex actions |
| gcore_waap_domain_firewall_rule | Low | IP/IP-range conditions only |
| gcore_waap_domain_advanced_rule | High | CEL syntax source, complex actions |
| gcore_waap_custom_page_set | Medium | Multiple page types, domain assignments |
| gcore_waap_domain_insight_silence | Medium | UUID IDs, RFC3339 datetime |
| gcore_waap_policy | Medium | Toggle-based, validation required |

---

## Phase 1: Fix Compiler Errors

### Issue: DomainID in ListParams structs

The generated data source models incorrectly include `DomainID` in params structs.

**Files to fix:**
- `internal/services/waap_domain_api_path/data_source_model.go`
- `internal/services/waap_domain_custom_rule/data_source_model.go`
- `internal/services/waap_domain_insight_silence/data_source_model.go`
- `internal/services/waap_domain_advanced_rule/data_source_model.go`
- `internal/services/waap_domain_firewall_rule/data_source_model.go`

**Pattern to apply:**

```go
// BEFORE (broken)
func (m *WaapDomainAPIPathDataSourceModel) toListParams(_ context.Context) (params waap.DomainAPIPathListParams, diags diag.Diagnostics) {
    params = waap.DomainAPIPathListParams{
        DomainID: m.DomainID.ValueInt64(),  // WRONG!
    }
    return
}

// AFTER (fixed)
func (m *WaapDomainAPIPathDataSourceModel) toListParams(_ context.Context) (params waap.DomainAPIPathListParams, diags diag.Diagnostics) {
    params = waap.DomainAPIPathListParams{
        // DomainID removed - it's a path param, not query param
    }
    return
}

// Also fix in data_source.go - pass domainID as separate argument
res, err := r.client.Waap.DomainAPIPaths.List(ctx, data.DomainID.ValueInt64(), params)
```

---

## Phase 2: waap_domain Resource

### Special Behavior
From old provider analysis, WAAP domains:
1. **Cannot be created** - must look up existing domain by name
2. **Cannot be deleted** - delete is no-op
3. **Updates require 3 API calls**

### Implementation

```go
// resource.go

func (r *WaapDomainResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data WaapDomainResourceModel
    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

    // 1. List domains and find by name (NOT create!)
    domains, err := r.client.Waap.Domains.List(ctx, waap.DomainListParams{
        Name: param.NewOpt(data.Name.ValueString()),
    })

    // 2. Find matching domain
    domain := findDomainByName(domains.Results, data.Name.ValueString())
    if domain == nil {
        resp.Diagnostics.AddError("Domain not found",
            fmt.Sprintf("Domain '%s' must exist before it can be managed", data.Name.ValueString()))
        return
    }

    // 3. Update status if needed
    if !data.Status.IsNull() && data.Status.ValueString() != string(domain.Status) {
        r.client.Waap.Domains.Update(ctx, domain.ID, waap.DomainUpdateParams{
            Status: waap.DomainUpdateParamsStatus(data.Status.ValueString()),
        })
    }

    // 4. Update settings if specified
    if !data.Settings.IsNull() {
        r.updateDomainSettings(ctx, domain.ID, data.Settings)
    }

    // 5. Update API discovery settings if specified
    if !data.APIDiscoverySettings.IsNull() {
        r.updateAPIDiscoverySettings(ctx, domain.ID, data.APIDiscoverySettings)
    }

    data.ID = types.Int64Value(domain.ID)
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *WaapDomainResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    // No-op - domains are managed externally
}
```

### Schema Updates (schema.go)

```go
"name": schema.StringAttribute{
    Required:    true,
    Description: "Name of the domain to manage",
    PlanModifiers: []planmodifier.String{
        stringplanmodifier.RequiresReplace(),
    },
},
"status": schema.StringAttribute{
    Optional:    true,
    Computed:    true,
    Description: "Domain status: active or monitor",
    Validators: []validator.String{
        stringvalidator.OneOf("active", "monitor"),
    },
},
"settings": schema.SingleNestedAttribute{
    Optional: true,
    Computed: true,
    Attributes: map[string]schema.Attribute{
        "ddos": schema.SingleNestedAttribute{...},
        "api": schema.SingleNestedAttribute{...},
    },
},
"api_discovery_settings": schema.SingleNestedAttribute{
    Optional: true,
    Computed: true,
    Attributes: map[string]schema.Attribute{
        "description_file_location": schema.StringAttribute{Required: true},
        "description_file_scan_enabled": schema.BoolAttribute{Optional: true, Computed: true},
        "description_file_scan_interval_hours": schema.Int64Attribute{Optional: true, Computed: true},
        "traffic_scan_enabled": schema.BoolAttribute{Optional: true, Computed: true},
        "traffic_scan_interval_hours": schema.Int64Attribute{Optional: true, Computed: true},
    },
},
```

---

## Phase 3: Rule Resources (custom_rule, advanced_rule, firewall_rule)

### Action Schema Pattern

From TERRAFORM_GUIDE.md, actions with empty object types should use `jsontypes.Normalized`:

```go
"action": schema.SingleNestedAttribute{
    Required: true,
    Attributes: map[string]schema.Attribute{
        "allow": schema.StringAttribute{
            CustomType:  jsontypes.NormalizedType{},
            Optional:    true,
            Description: "Allow action (empty object)",
        },
        "block": schema.SingleNestedAttribute{
            Optional: true,
            Attributes: map[string]schema.Attribute{
                "status_code": schema.Int64Attribute{
                    Optional: true,
                    Validators: []validator.Int64{
                        int64validator.OneOf(403, 405, 418, 429),
                    },
                },
                "action_duration": schema.StringAttribute{
                    Optional:    true,
                    Description: "Duration like '5m', '12h'",
                },
            },
        },
        "captcha": schema.StringAttribute{
            CustomType: jsontypes.NormalizedType{},
            Optional:   true,
        },
        "handshake": schema.StringAttribute{
            CustomType: jsontypes.NormalizedType{},
            Optional:   true,
        },
        "monitor": schema.StringAttribute{
            CustomType: jsontypes.NormalizedType{},
            Optional:   true,
        },
        "tag": schema.SingleNestedAttribute{
            Optional: true,
            Attributes: map[string]schema.Attribute{
                "tags": schema.ListAttribute{
                    ElementType: types.StringType,
                    Required:    true,
                },
            },
        },
    },
},
```

### Condition Serialization (custom_rule)

Conditions require custom serialization for 18+ types:

```go
// helpers.go

func conditionsToAPIPayload(conditions []ConditionModel) []waap.Condition {
    var result []waap.Condition
    for _, c := range conditions {
        if c.IP != nil {
            result = append(result, waap.Condition{
                IP: &waap.IPCondition{
                    IPAddress: c.IP.IPAddress.ValueString(),
                    Negation:  c.IP.Negation.ValueBool(),
                },
            })
        }
        if c.IPRange != nil {
            result = append(result, waap.Condition{
                IPRange: &waap.IPRangeCondition{
                    LowerBound: c.IPRange.LowerBound.ValueString(),
                    UpperBound: c.IPRange.UpperBound.ValueString(),
                    Negation:   c.IPRange.Negation.ValueBool(),
                },
            })
        }
        // ... handle all 18+ condition types
    }
    return result
}

func conditionsFromAPIResponse(conditions []waap.Condition) []ConditionModel {
    // Reverse mapping
}
```

---

## Phase 4: Import Support

### Import ID Format

| Resource | Import Format |
|----------|---------------|
| waap_domain | `{domain_id}` |
| waap_domain_api_path | `{domain_id}:{path_id}` |
| waap_domain_custom_rule | `{domain_id}:{rule_id}` |
| waap_domain_firewall_rule | `{domain_id}:{rule_id}` |
| waap_domain_advanced_rule | `{domain_id}:{rule_id}` |
| waap_custom_page_set | `{set_id}` |
| waap_domain_insight_silence | `{domain_id}:{silence_id}` |

### Implementation Pattern

```go
func (r *WaapDomainCustomRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    parts := strings.Split(req.ID, ":")
    if len(parts) != 2 {
        resp.Diagnostics.AddError("Invalid import ID",
            "Expected format: domain_id:rule_id")
        return
    }

    domainID, err := strconv.ParseInt(parts[0], 10, 64)
    if err != nil {
        resp.Diagnostics.AddError("Invalid domain_id", err.Error())
        return
    }

    ruleID, err := strconv.ParseInt(parts[1], 10, 64)
    if err != nil {
        resp.Diagnostics.AddError("Invalid rule_id", err.Error())
        return
    }

    resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("domain_id"), domainID)...)
    resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), ruleID)...)
}
```

---

## Phase 5: Testing

### Test Manifests

```hcl
# test-waap-domain/main.tf
resource "gcore_waap_domain" "test" {
  name   = "example.com"
  status = "active"

  settings {
    ddos {
      global_threshold = 1000
      burst_threshold  = 100
    }
    api {
      is_api   = false
      api_urls = ["/api/*"]
    }
  }

  api_discovery_settings {
    description_file_location          = "https://example.com/openapi.json"
    description_file_scan_enabled      = true
    description_file_scan_interval_hours = 24
    traffic_scan_enabled               = true
    traffic_scan_interval_hours        = 12
  }
}

# test-waap-custom-rule/main.tf
resource "gcore_waap_domain_custom_rule" "test" {
  domain_id   = gcore_waap_domain.test.id
  name        = "Block bad IPs"
  description = "Block known malicious IPs"
  enabled     = true

  action {
    block {
      status_code     = 403
      action_duration = "1h"
    }
  }

  conditions {
    ip {
      ip_address = "192.168.1.1"
      negation   = false
    }
  }
}
```

### Verification Checklist

For each resource:
- [ ] `terraform plan` - no errors
- [ ] `terraform apply` - creates/updates successfully
- [ ] `terraform plan` (after apply) - no drift
- [ ] `terraform import` - imports correctly
- [ ] `terraform state show` - state matches API
- [ ] `terraform destroy` - cleans up (or no-op for domains)

---

## Files Changed Summary

### New/Modified Files

```
internal/services/waap_domain/
├── resource.go         # Custom create (lookup), delete (no-op), update
├── schema.go           # Add settings, api_discovery_settings
├── model.go            # Add nested models
└── data_source_model.go # Fix DomainID handling

internal/services/waap_domain_*/
├── resource.go         # Custom action/condition handling
├── data_source.go      # Fix path param passing
├── data_source_model.go # Remove DomainID from params
└── helpers.go          # Condition serialization (new file)

internal/custom/
└── waap.go            # Shared WAAP helpers (new file)
```

---

## Timeline Estimate

| Phase | Tasks | Effort |
|-------|-------|--------|
| 1 | Fix compiler errors | 2h |
| 2 | waap_domain resource | 4h |
| 3 | Rule resources | 8h |
| 4 | Import support | 2h |
| 5 | Testing | 4h |
| **Total** | | **20h** |

---

## Dependencies

1. **OAS fixes** must be applied first (see WAAP_OAS_REQUIREMENTS.md)
2. **Go SDK** merge conflict PR #182 must be resolved
3. **Stainless rebuild** after OAS fixes
