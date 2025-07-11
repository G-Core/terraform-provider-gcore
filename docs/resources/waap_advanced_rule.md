---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "gcore_waap_advanced_rule Resource - terraform-provider-gcore"
subcategory: ""
description: |-
  Represent Advanced Rules for a specific WAAP domain
---

# gcore_waap_advanced_rule (Resource)

Represent Advanced Rules for a specific WAAP domain

## Example Usage

```terraform
provider gcore {
  permanent_api_token = "768660$.............a43f91f"
}

resource "gcore_cdn_resource" "cdn_resource" {
  cname  = "api.example.com"
  origin = "origin.example.com"
  options {
    waap { value = true }
  }
}

resource "gcore_waap_domain" "domain" {
  name   = gcore_cdn_resource.cdn_resource.cname
  status = "monitor"
}

resource "gcore_waap_advanced_rule" "advanced_rule" {
  domain_id = gcore_waap_domain.domain.id
  name = "Advanced Rule"
  enabled = true
  action {
    block {
      status_code = 403
      action_duration = "5m"
    }
  }
  source = "request.ip == '117.20.32.55'"
  description = "Description of the advanced rule"
  phase = "access"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `action` (Block List, Min: 1, Max: 1) The action that the rule takes when triggered. (see [below for nested schema](#nestedblock--action))
- `domain_id` (Number) The WAAP domain ID for which the Advanced Rule is configured.
- `enabled` (Boolean) Whether the rule is enabled.
- `name` (String) The name assigned to the rule.
- `source` (String) A CEL syntax expression that contains the rule's conditions. Allowed objects are: request, whois, session, response, tags, user_defined_tags, user_agent, client_data. More info can be found here: https://gcore.com/docs/waap/waap-rules/advanced-rules

### Optional

- `description` (String) The description assigned to the rule.
- `phase` (String) The WAAP request/response phase for applying the rule. The 'access' phase is responsible for modifying the request before it is sent to the origin server. The 'header_filter' phase is responsible for modifying the HTTP headers of a response before they are sent back to the client.The 'body_filter' phase is responsible for modifying the body of a response before it is sent back to the client. Default is 'access'.

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--action"></a>
### Nested Schema for `action`

Optional:

- `allow` (Boolean) The WAAP allows the request.
- `block` (Block List, Max: 1) The WAAP blocks the request. (see [below for nested schema](#nestedblock--action--block))
- `captcha` (Boolean) The WAAP requires the user to solve a CAPTCHA challenge.
- `handshake` (Boolean) The WAAP performs automatic browser validation.
- `monitor` (Boolean) The WAAP monitors the request but took no action.
- `tag` (Block List, Max: 1) The WAAP tags the request. (see [below for nested schema](#nestedblock--action--tag))

<a id="nestedblock--action--block"></a>
### Nested Schema for `action.block`

Optional:

- `action_duration` (String) How long a rule's block action will apply to subsequent requests. Can be specified in seconds or by using a numeral followed by 's', 'm', 'h', or 'd' to represent time format (seconds, minutes, hours, or days). Example: 12h. Must match the pattern ^[0-9]*[smhd]?$
- `status_code` (Number) A custom HTTP status code that the WAAP returns if a rule blocks a request. It must be one of these values {403, 405, 418, 429}. Default is 403.


<a id="nestedblock--action--tag"></a>
### Nested Schema for `action.tag`

Required:

- `tags` (Set of String) The list of user defined tags to tag the request with.

## Import

Import is supported using the following syntax:

The [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import) can be used, for example:

```shell
# import using <domain_id>:<rule_id>
terraform import gcore_waap_advanced_rule.advanced_rule 10029:98347
```
