---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "gcore_waap_custom_rule Resource - terraform-provider-gcore"
subcategory: ""
description: |-
  Represent Custom Rules for a specific WAAP domain
---

# gcore_waap_custom_rule (Resource)

Represent Custom Rules for a specific WAAP domain

## Example Usage

```terraform
provider gcore {
    permanent_api_token = "768660$.............a43f91f"
}

data "gcore_waap_tag" "proxy_network" {
  name = "Proxy Network"
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

resource "gcore_waap_custom_rule" "custom_rule" {
    domain_id = gcore_waap_domain.domain.id
    name = "Custom Rule"
    enabled = true
    action {
        block {
            status_code = 403
            action_duration = "5m"
        }
    }
    conditions {
        ip {
            negation    = true
            ip_address = "192.168.0.6"
        }
        http_method {
            http_method  = "POST"
        }
        tags {
            tags = [data.gcore_waap_tag.proxy_network.id]
        }
    }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `action` (Block List, Min: 1, Max: 1) The action that the rule takes when triggered. (see [below for nested schema](#nestedblock--action))
- `conditions` (Block List, Min: 1, Max: 1) The conditions required for the WAAP engine to trigger the rule. Rules may have between 1 and 5 conditions. All conditions must pass for the rule to trigger. (see [below for nested schema](#nestedblock--conditions))
- `domain_id` (Number) The WAAP domain ID for which the Custom Rule is configured.
- `enabled` (Boolean) Whether the rule is enabled.
- `name` (String) The name assigned to the rule.

### Optional

- `description` (String) The description assigned to the rule.

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



<a id="nestedblock--conditions"></a>
### Nested Schema for `conditions`

Optional:

- `content_type` (Block List) Content type condition. This condition matches the content type of the request. (see [below for nested schema](#nestedblock--conditions--content_type))
- `country` (Block List) Country condition. This condition matches the country of the request based on the source IP address. (see [below for nested schema](#nestedblock--conditions--country))
- `file_extension` (Block List) File extension condition. This condition matches the file extension of the request. (see [below for nested schema](#nestedblock--conditions--file_extension))
- `header` (Block List) Request header condition. This condition matches a request header and its value. (see [below for nested schema](#nestedblock--conditions--header))
- `header_exists` (Block List) Request header exists condition. This condition checks if a request header exists. (see [below for nested schema](#nestedblock--conditions--header_exists))
- `http_method` (Block List) HTTP method condition. This condition matches the HTTP method of the request. (see [below for nested schema](#nestedblock--conditions--http_method))
- `ip` (Block List) IP address condition. This condition matches a single IP address. (see [below for nested schema](#nestedblock--conditions--ip))
- `ip_range` (Block List) IP range condition. This condition matches a range of IP addresses. (see [below for nested schema](#nestedblock--conditions--ip_range))
- `organization` (Block List) Organization condition. This condition matches the organization of the request based on the source IP address. (see [below for nested schema](#nestedblock--conditions--organization))
- `owner_types` (Block List) (see [below for nested schema](#nestedblock--conditions--owner_types))
- `request_rate` (Block List) Request rate condition. This condition matches the request rate. (see [below for nested schema](#nestedblock--conditions--request_rate))
- `response_header` (Block List) (see [below for nested schema](#nestedblock--conditions--response_header))
- `response_header_exists` (Block List) Response header exists condition. This condition checks if a response header exists. (see [below for nested schema](#nestedblock--conditions--response_header_exists))
- `session_request_count` (Block List) Session request count condition. This condition matches the number of dynamic requests in the session. (see [below for nested schema](#nestedblock--conditions--session_request_count))
- `tags` (Block List) Tags condition. This condition matches the request tags. (see [below for nested schema](#nestedblock--conditions--tags))
- `url` (Block List) URL condition. This condition matches a URL path. (see [below for nested schema](#nestedblock--conditions--url))
- `user_agent` (Block List) User agent condition. This condition matches the user agent of the request. (see [below for nested schema](#nestedblock--conditions--user_agent))
- `user_defined_tags` (Block List) (see [below for nested schema](#nestedblock--conditions--user_defined_tags))

<a id="nestedblock--conditions--content_type"></a>
### Nested Schema for `conditions.content_type`

Required:

- `content_type` (Set of String) The list of content types to match against.

Optional:

- `negation` (Boolean) Whether or not to apply a boolean NOT operation to the rule's condition.


<a id="nestedblock--conditions--country"></a>
### Nested Schema for `conditions.country`

Required:

- `country_code` (Set of String) A list of ISO 3166-1 alpha-2 formatted strings representing the countries to match against.

Optional:

- `negation` (Boolean) Whether or not to apply a boolean NOT operation to the rule's condition.


<a id="nestedblock--conditions--file_extension"></a>
### Nested Schema for `conditions.file_extension`

Required:

- `file_extension` (Set of String) The list of file extensions to match against.

Optional:

- `negation` (Boolean) Whether or not to apply a boolean NOT operation to the rule's condition.


<a id="nestedblock--conditions--header"></a>
### Nested Schema for `conditions.header`

Required:

- `header` (String) The request header name.
- `value` (String) The request header value.

Optional:

- `match_type` (String) The type of matching condition. Valid values are 'Exact', 'Contains'. Default is 'Contains'.
- `negation` (Boolean) Whether or not to apply a boolean NOT operation to the rule's condition.


<a id="nestedblock--conditions--header_exists"></a>
### Nested Schema for `conditions.header_exists`

Required:

- `header` (String) The request header name.

Optional:

- `negation` (Boolean) Whether or not to apply a boolean NOT operation to the rule's condition.


<a id="nestedblock--conditions--http_method"></a>
### Nested Schema for `conditions.http_method`

Required:

- `http_method` (String) The HTTP method to match against. Valid values are 'CONNECT', 'DELETE', 'GET', 'HEAD', 'OPTIONS', 'PATCH', 'POST', 'PUT', and 'TRACE'.

Optional:

- `negation` (Boolean) Whether or not to apply a boolean NOT operation to the rule's condition.


<a id="nestedblock--conditions--ip"></a>
### Nested Schema for `conditions.ip`

Required:

- `ip_address` (String) A single IPv4 or IPv6 address

Optional:

- `negation` (Boolean) Whether or not to apply a boolean NOT operation to the rule's condition.


<a id="nestedblock--conditions--ip_range"></a>
### Nested Schema for `conditions.ip_range`

Required:

- `lower_bound` (String) The lower bound IPv4 or IPv6 address to match against.
- `upper_bound` (String) The upper bound IPv4 or IPv6 address to match against.

Optional:

- `negation` (Boolean) Whether or not to apply a boolean NOT operation to the rule's condition.


<a id="nestedblock--conditions--organization"></a>
### Nested Schema for `conditions.organization`

Required:

- `organization` (String) The organization to match against.

Optional:

- `negation` (Boolean) Whether or not to apply a boolean NOT operation to the rule's condition.


<a id="nestedblock--conditions--owner_types"></a>
### Nested Schema for `conditions.owner_types`

Required:

- `owner_types` (Set of String) Match the type of organization that owns the IP address making an incoming request. Valid values are 'COMMERCIAL', 'EDUCATIONAL', 'GOVERNMENT', 'HOSTING_SERVICES', 'ISP', 'MOBILE_NETWORK', 'NETWORK', and 'RESERVED'.

Optional:

- `negation` (Boolean) Whether or not to apply a boolean NOT operation to the rule's condition.


<a id="nestedblock--conditions--request_rate"></a>
### Nested Schema for `conditions.request_rate`

Required:

- `path_pattern` (String) A regular expression matching the URL path of the incoming request.
- `requests` (Number) The number of incoming requests over the given time that can trigger a request rate condition.
- `time` (Number) The number of seconds that the WAAP measures incoming requests over before triggering a request rate condition.

Optional:

- `http_methods` (Set of String) Possible HTTP request methods that can trigger a request rate condition. Valid values are 'CONNECT', 'DELETE', 'GET', 'HEAD', 'OPTIONS', 'PATCH', 'POST', 'PUT', and 'TRACE'.
- `ips` (Set of String) A list of source IPs that can trigger a request rate condition.
- `user_defined_tag` (String) A user-defined tag that can be included in incoming requests and used to trigger a request rate condition.


<a id="nestedblock--conditions--response_header"></a>
### Nested Schema for `conditions.response_header`

Required:

- `header` (String) The request header name.
- `value` (String) The request header value.

Optional:

- `match_type` (String) The type of matching condition. Valid values are 'Exact', 'Contains'. Default is 'Contains'.
- `negation` (Boolean) Whether or not to apply a boolean NOT operation to the rule's condition.


<a id="nestedblock--conditions--response_header_exists"></a>
### Nested Schema for `conditions.response_header_exists`

Required:

- `header` (String) The request header name.

Optional:

- `negation` (Boolean) Whether or not to apply a boolean NOT operation to the rule's condition.


<a id="nestedblock--conditions--session_request_count"></a>
### Nested Schema for `conditions.session_request_count`

Required:

- `request_count` (Number) The number of dynamic requests in the session.

Optional:

- `negation` (Boolean) Whether or not to apply a boolean NOT operation to the rule's condition.


<a id="nestedblock--conditions--tags"></a>
### Nested Schema for `conditions.tags`

Required:

- `tags` (Set of String) A list of tags to match against the request tags. Tags can be obtained from the API endpoint /v1/tags or you can use the gcore_waap_tag data source.

Optional:

- `negation` (Boolean) Whether or not to apply a boolean NOT operation to the rule's condition.


<a id="nestedblock--conditions--url"></a>
### Nested Schema for `conditions.url`

Required:

- `url` (String) The URL to match.

Optional:

- `match_type` (String) The type of matching condition. Valid values are 'Exact', 'Contains', and 'Regex'. Default is 'Contains'.
- `negation` (Boolean) Whether or not to apply a boolean NOT operation to the rule's condition.


<a id="nestedblock--conditions--user_agent"></a>
### Nested Schema for `conditions.user_agent`

Required:

- `user_agent` (String) The user agent value to match.

Optional:

- `match_type` (String) The type of matching condition. Valid values are 'Exact', 'Contains'. Default is 'Contains'.
- `negation` (Boolean) Whether or not to apply a boolean NOT operation to the rule's condition.


<a id="nestedblock--conditions--user_defined_tags"></a>
### Nested Schema for `conditions.user_defined_tags`

Required:

- `tags` (Set of String) A list of user-defined tags to match against the request tags.

Optional:

- `negation` (Boolean) Whether or not to apply a boolean NOT operation to the rule's condition.

## Import

Import is supported using the following syntax:

The [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import) can be used, for example:

```shell
# import using <domain_id>:<rule_id>
terraform import gcore_waap_custom_rule.custom_rule 10029:98347
```
