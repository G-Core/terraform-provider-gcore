# WAAP OpenAPI Spec (OAS) Requirements for Terraform Support

This document outlines required changes to the Gcore WAAP OpenAPI specification to fix Stainless Terraform generation issues.

## Current Issues Summary

| Issue Type | Count | Impact |
|------------|-------|--------|
| Missing response schemas | 12 endpoints | Resources cannot be inferred |
| Undefined property types | 12 properties | Generates raw JSON fields |
| Compiler errors | 5 files | Build fails |

---

## 1. Missing Response Schemas

### Problem
Stainless reports: `"Request and response types do not match. Response type: undefined"`

The create (POST) and update (PATCH) endpoints don't return response bodies with proper schemas, so Terraform cannot infer the resource state after operations.

### Required Changes

#### 1.1 API Paths
```yaml
/waap/v1/domains/{domain_id}/api-paths:
  post:
    responses:
      '201':
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/WaapAPIPath'  # ADD THIS

/waap/v1/domains/{domain_id}/api-paths/{path_id}:
  patch:
    responses:
      '204':
        description: No Content  # Keep as-is, or add 200 with body
```

#### 1.2 Custom Rules
```yaml
/waap/v1/domains/{domain_id}/custom-rules:
  post:
    responses:
      '201':
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/CustomRule'  # ADD THIS

/waap/v1/domains/{domain_id}/custom-rules/{rule_id}:
  patch:
    responses:
      '204':
        description: No Content
```

#### 1.3 Firewall Rules
```yaml
/waap/v1/domains/{domain_id}/firewall-rules:
  post:
    responses:
      '201':
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/FirewallRule'  # ADD THIS
```

#### 1.4 Advanced Rules
```yaml
/waap/v1/domains/{domain_id}/advanced-rules:
  post:
    responses:
      '201':
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/AdvancedRule'  # ADD THIS
```

#### 1.5 Custom Page Sets
```yaml
/waap/v1/custom-page-sets:
  post:
    responses:
      '201':
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/CustomPageSet'  # ADD THIS
```

#### 1.6 Insight Silences
```yaml
/waap/v1/domains/{domain_id}/insight-silences:
  post:
    responses:
      '200':
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/InsightSilence'  # ADD THIS
```

---

## 2. Undefined Property Types (Ambiguous Properties)

### Problem
Stainless reports: `"resource has properties with an undefined type, which results in a raw JSON field"`

Action properties like `allow`, `captcha`, `handshake`, `monitor` are defined without explicit types.

### Current (Broken)
```yaml
CustomerRuleAction-Input:
  type: object
  properties:
    allow: {}  # No type defined!
    captcha: {}
    handshake: {}
    monitor: {}
    block:
      $ref: '#/components/schemas/BlockActionInput'
    tag:
      $ref: '#/components/schemas/TagActionInput'
```

### Required Changes

#### 2.1 CustomerRuleAction-Input (for custom_rule and advanced_rule)
```yaml
CustomerRuleAction-Input:
  type: object
  properties:
    allow:
      type: object
      nullable: true
      description: "Allow action - empty object to allow request"
    captcha:
      type: object
      nullable: true
      description: "Captcha challenge action"
    handshake:
      type: object
      nullable: true
      description: "Handshake verification action"
    monitor:
      type: object
      nullable: true
      description: "Monitor action - log without blocking"
    block:
      $ref: '#/components/schemas/BlockActionInput'
    tag:
      $ref: '#/components/schemas/TagActionInput'
```

#### 2.2 CustomerRuleAction-Output
```yaml
CustomerRuleAction-Output:
  type: object
  properties:
    allow:
      type: object
      nullable: true
    captcha:
      type: object
      nullable: true
    handshake:
      type: object
      nullable: true
    monitor:
      type: object
      nullable: true
    block:
      $ref: '#/components/schemas/BlockActionOutput'
    tag:
      $ref: '#/components/schemas/TagActionOutput'
```

#### 2.3 FirewallRuleAction-Input
```yaml
FirewallRuleAction-Input:
  type: object
  properties:
    allow:
      type: object
      nullable: true
      description: "Allow action for firewall rule"
    block:
      $ref: '#/components/schemas/BlockActionInput'
```

#### 2.4 FirewallRuleAction-Output
```yaml
FirewallRuleAction-Output:
  type: object
  properties:
    allow:
      type: object
      nullable: true
    block:
      $ref: '#/components/schemas/BlockActionOutput'
```

---

## 3. Stainless Terraform Configuration

### Problem
Some WAAP resources have special lifecycle behaviors that Stainless cannot infer.

### Required Changes in `openapi.stainless.yml`

```yaml
resources:
  waap:
    subresources:
      domains:
        terraform:
          # Domain is lookup-only, not created via API
          data_source: true
          resource:
            methods:
              create:
                # Use list + filter instead of create
                behavior: lookup_by_name
              delete:
                # Domains cannot be deleted
                behavior: noop
```

---

## 4. Complete List of Schema References to Fix

### Files in api-schemas/openapi.yml

| Schema Path | Issue | Fix |
|-------------|-------|-----|
| `#/components/schemas/CustomerRuleAction-Input/properties/allow` | undefined | Add `type: object` |
| `#/components/schemas/CustomerRuleAction-Input/properties/captcha` | undefined | Add `type: object` |
| `#/components/schemas/CustomerRuleAction-Input/properties/handshake` | undefined | Add `type: object` |
| `#/components/schemas/CustomerRuleAction-Input/properties/monitor` | undefined | Add `type: object` |
| `#/components/schemas/CustomerRuleAction-Output/properties/allow` | undefined | Add `type: object` |
| `#/components/schemas/CustomerRuleAction-Output/properties/captcha` | undefined | Add `type: object` |
| `#/components/schemas/CustomerRuleAction-Output/properties/handshake` | undefined | Add `type: object` |
| `#/components/schemas/CustomerRuleAction-Output/properties/monitor` | undefined | Add `type: object` |
| `#/components/schemas/FirewallRuleAction-Input/properties/allow` | undefined | Add `type: object` |
| `#/components/schemas/FirewallRuleAction-Output/properties/allow` | undefined | Add `type: object` |

---

## 5. Endpoints Requiring Response Schema

| Endpoint | Method | Current Response | Required Response |
|----------|--------|-----------------|-------------------|
| `/waap/v1/domains/{domain_id}/api-paths` | POST | 201 (no body) | 201 with WaapAPIPath |
| `/waap/v1/domains/{domain_id}/api-paths/{path_id}` | PATCH | 204 | 204 (ok) or 200 with body |
| `/waap/v1/domains/{domain_id}/custom-rules` | POST | 201 (no body) | 201 with CustomRule |
| `/waap/v1/domains/{domain_id}/custom-rules/{rule_id}` | PATCH | 204 | 204 (ok) |
| `/waap/v1/domains/{domain_id}/firewall-rules` | POST | 201 (no body) | 201 with FirewallRule |
| `/waap/v1/domains/{domain_id}/firewall-rules/{rule_id}` | PATCH | 204 | 204 (ok) |
| `/waap/v1/domains/{domain_id}/advanced-rules` | POST | 201 (no body) | 201 with AdvancedRule |
| `/waap/v1/domains/{domain_id}/advanced-rules/{rule_id}` | PATCH | 204 | 204 (ok) |
| `/waap/v1/custom-page-sets` | POST | 201 (no body) | 201 with CustomPageSet |
| `/waap/v1/custom-page-sets/{set_id}` | PATCH | 204 | 204 (ok) |
| `/waap/v1/domains/{domain_id}/insight-silences` | POST | 200 (no body) | 200 with InsightSilence |
| `/waap/v1/domains/{domain_id}/insight-silences/{silence_id}` | PATCH | 200 (no body) | 200 with InsightSilence |

---

## 6. Verification Steps

After making OAS changes:

1. Run Stainless build: `stainless build --project gcore --target terraform`
2. Check diagnostics: `stainless diagnostics --project gcore --target terraform`
3. Verify no "undefined type" errors
4. Verify no "cannot infer response" errors
5. Run `go build` on generated Terraform code

---

## 7. Related Jira Tickets

- **GCLOUD2-22875**: Support gcore_waap_domain in Stainless terraform provider
- **Sprint 48 Goal**: Discover the problems of Gcore products and their OAS

---

## 8. Priority

| Priority | Issue | Impact |
|----------|-------|--------|
| P0 | Missing response schemas | Blocks resource generation |
| P0 | Undefined action types | Generates unusable JSON fields |
| P1 | Path parameter handling | Compiler errors in data sources |
| P2 | Stainless config for special behavior | Custom logic needed anyway |
