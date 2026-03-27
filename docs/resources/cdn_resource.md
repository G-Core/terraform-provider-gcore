---
page_title: "gcore_cdn_resource Resource - Gcore"
subcategory: ""
description: |-
  
---

# gcore_cdn_resource (Resource)



## Example Usage

### Basic CDN resource

Create a CDN resource with an origin group and common caching and security options.

```terraform
resource "gcore_cdn_origin_group" "example" {
  name     = "origin_group_1"
  use_next = true
  sources = [{
    source  = "example.com"
    enabled = true
  }]
}

resource "gcore_cdn_resource" "example" {
  cname               = "cdn.example.com"
  origin_group        = gcore_cdn_origin_group.example.id
  origin_protocol     = "MATCH"
  secondary_hostnames = ["cdn2.example.com"]

  options = {
    edge_cache_settings = {
      enabled = true
      default = "8d"
    }
    browser_cache_settings = {
      enabled = true
      value   = "1d"
    }
    redirect_http_to_https = {
      enabled = true
      value   = true
    }
    gzip_on = {
      enabled = true
      value   = true
    }
    cors = {
      enabled = true
      value   = ["*"]
    }
    rewrite = {
      enabled = true
      body    = "/(.*) /$1"
    }
    tls_versions = {
      enabled = true
      value   = ["TLSv1.2"]
    }
    force_return = {
      enabled = true
      code    = 200
      body    = "OK"
    }
    request_limiter = {
      enabled   = true
      rate_unit = "r/s"
      rate      = 5
    }
  }
}
```

### Advanced CDN resource

Create a CDN resource with a comprehensive set of options including caching, security, access control, compression, request/response manipulation, rate limiting, and more.

```terraform
resource "gcore_cdn_origin_group" "example" {
  name     = "origin_group_1"
  use_next = true
  sources = [{
    source  = "example.com"
    enabled = true
  }]
}

resource "gcore_cdn_resource" "example" {
  cname           = "cdn.example.com"
  origin_group    = gcore_cdn_origin_group.example.id
  origin_protocol = "HTTPS"
  ssl_enabled     = true
  active          = true
  description     = "CDN resource with advanced options"

  secondary_hostnames = ["cdn2.example.com", "cdn3.example.com"]

  options = {
    # Caching
    edge_cache_settings = {
      enabled = true
      value   = "43200s"
      custom_values = {
        "100" = "400s"
        "101" = "400s"
      }
    }
    browser_cache_settings = {
      enabled = true
      value   = "3600s"
    }
    ignore_cookie = {
      enabled = true
      value   = true
    }
    ignore_query_string = {
      enabled = true
      value   = false
    }
    slice = {
      enabled = true
      value   = true
    }
    stale = {
      enabled = true
      value   = ["http_404", "http_500"]
    }

    # Security
    redirect_http_to_https = {
      enabled = true
      value   = true
    }
    tls_versions = {
      enabled = true
      value   = ["TLSv1.2"]
    }
    secure_key = {
      enabled = true
      key     = "secret"
      type    = 2
    }
    cors = {
      enabled = true
      value   = ["*"]
      always  = true
    }

    # Access control
    country_acl = {
      enabled         = true
      policy_type     = "allow"
      excepted_values = ["GB", "DE"]
    }
    ip_address_acl = {
      enabled         = true
      policy_type     = "deny"
      excepted_values = ["192.168.1.100/32"]
    }
    referrer_acl = {
      enabled         = true
      policy_type     = "deny"
      excepted_values = ["*.google.com"]
    }
    user_agent_acl = {
      enabled         = true
      policy_type     = "allow"
      excepted_values = ["UserAgent"]
    }

    # Compression
    gzip_on = {
      enabled = true
      value   = true
    }
    brotli_compression = {
      enabled = true
      value   = ["text/html", "text/plain"]
    }
    fetch_compressed = {
      enabled = true
      value   = false
    }

    # Origin settings
    host_header = {
      enabled = true
      value   = "host.com"
    }
    forward_host_header = {
      enabled = true
      value   = false
    }
    sni = {
      enabled         = true
      sni_type        = "custom"
      custom_hostname = "custom.example.com"
    }

    # Request / response manipulation
    rewrite = {
      enabled = true
      body    = "/(.*) /additional_path/$1"
      flag    = "break"
    }
    static_request_headers = {
      enabled = true
      value = {
        "X-Custom" = "X-Request"
      }
    }
    static_response_headers = {
      enabled = true
      value = [{
        name   = "X-Custom1"
        value  = ["Value1", "Value2"]
        always = false
        }, {
        name   = "X-Custom2"
        value  = ["CDN"]
        always = true
      }]
    }
    response_headers_hiding_policy = {
      enabled  = true
      mode     = "hide"
      excepted = ["my-header"]
    }

    # Rate limiting
    request_limiter = {
      enabled   = true
      rate      = 5
      rate_unit = "r/s"
    }
    limit_bandwidth = {
      enabled    = true
      limit_type = "static"
      speed      = 100
      buffer     = 200
    }

    # Proxy settings
    proxy_cache_methods_set = {
      enabled = true
      value   = false
    }
    proxy_connect_timeout = {
      enabled = true
      value   = "4s"
    }
    proxy_read_timeout = {
      enabled = true
      value   = "10s"
    }

    # Other
    allowed_http_methods = {
      enabled = true
      value   = ["GET", "POST"]
    }
    follow_origin_redirect = {
      enabled = true
      codes   = [301, 302]
    }
    websockets = {
      enabled = true
      value   = true
    }
    http3_enabled = {
      enabled = true
      value   = true
    }
    image_stack = {
      enabled      = true
      quality      = 80
      avif_enabled = true
      webp_enabled = false
      png_lossless = true
    }
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `cname` (String) Delivery domains that will be used for content delivery through a CDN.

Delivery domains should be added to your DNS settings.

### Optional

- `active` (Boolean) Enables or disables a CDN resource.

Possible values:
- **true** - CDN resource is active. Content is being delivered.
- **false** - CDN resource is deactivated. Content is not being delivered.
- `description` (String) Optional comment describing the CDN resource.
- `name` (String) CDN resource name.
- `options` (Attributes) List of options that can be configured for the CDN resource.

In case of `null` value the option is not added to the CDN resource.
Option may inherit its value from the global account settings. (see [below for nested schema](#nestedatt--options))
- `origin` (String) IP address or domain name of the origin and the port, if custom port is used.

You can use either the `origin` or `originGroup` parameter in the request.
- `origin_group` (Number) Origin group ID with which the CDN resource is associated.

You can use either the `origin` or `originGroup` parameter in the request.
- `origin_protocol` (String) Protocol used by CDN servers to request content from an origin source.

Possible values:
- **HTTPS** - CDN servers will connect to the origin via HTTPS.
- **HTTP** - CDN servers will connect to the origin via HTTP.
- **MATCH** - connection protocol will be chosen automatically (content on the origin source should be available for the CDN both through HTTP and HTTPS).

If protocol is not specified, HTTP is used to connect to an origin server.
Available values: "HTTP", "HTTPS", "MATCH".
- `primary_resource` (Number) ID of the main CDN resource which has a shared caching zone with a reserve CDN resource.

If the parameter is not empty, then the current CDN resource is the reserve.
You cannot change some options, create rules, set up origin shielding, or use the reserve CDN resource for Streaming.
- `proxy_ssl_ca` (Number) ID of the trusted CA certificate used to verify an origin.

It can be used only with `"proxy_ssl_enabled": true`.
- `proxy_ssl_data` (Number) ID of the SSL certificate used to verify an origin.

It can be used only with `"proxy_ssl_enabled": true`.
- `proxy_ssl_enabled` (Boolean) Enables or disables SSL certificate validation of the origin server before completing any connection.

Possible values:
- **true** - Origin SSL certificate validation is enabled.
- **false** - Origin SSL certificate validation is disabled.
- `secondary_hostnames` (Set of String) Additional delivery domains (CNAMEs) that will be used to deliver content via the CDN.

Up to ten additional CNAMEs are possible.
- `ssl_data` (Number) ID of the SSL certificate linked to the CDN resource.

Can be used only with `"sslEnabled": true`.
- `ssl_enabled` (Boolean) Defines whether the HTTPS protocol enabled for content delivery.

Possible values:
- **true** - HTTPS is enabled.
- **false** - HTTPS is disabled.
- `waap_api_domain_enabled` (Boolean) Defines whether the associated WAAP Domain is identified as an API Domain.

Possible values:
- **true** - The associated WAAP Domain is designated as an API Domain.
- **false** - The associated WAAP Domain is not designated as an API Domain.

### Read-Only

- `can_purge_by_urls` (Boolean) Defines whether the CDN resource can be used for purge by URLs feature.

It's available only in case the CDN resource has enabled `ignore_vary_header` option.
- `client` (Number) ID of an account to which the CDN resource belongs.
- `created` (String) Date of CDN resource creation.
- `full_custom_enabled` (Boolean) Defines whether the CDN resource has a custom configuration.

Possible values:
- **true** - CDN resource has a custom configuration. You cannot change resource settings, except for the SSL certificate. To change other settings, contact technical support.
- **false** - CDN resource has a regular configuration. You can change CDN resource settings.
- `id` (Number) CDN resource ID.
- `is_primary` (Boolean) Defines whether a CDN resource has a cache zone shared with other CDN resources.

Possible values:
- **true** - CDN resource is main and has a shared caching zone with other CDN resources, which are called reserve.
- **false** - CDN resource is reserve and it has a shared caching zone with the main CDN resource. You cannot change some options, create rules, set up origin shielding and use the reserve resource for Streaming.
- **null** - CDN resource does not have a shared cache zone.

The main CDN resource is specified in the `primary_resource` field. It cannot be suspended unless all related reserve CDN resources are suspended.
- `origin_group_name` (String) Origin group name.
- `preset_applied` (Boolean) Defines whether the CDN resource has a preset applied.

Possible values:
- **true** - CDN resource has a preset applied. CDN resource options included in the preset cannot be edited.
- **false** - CDN resource does not have a preset applied.
- `rules` (List of String) Rules configured for the CDN resource.
- `shield_dc` (String) Name of the origin shielding location data center.

Parameter returns **null** if origin shielding is disabled.
- `shield_enabled` (Boolean) Defines whether origin shield is active and working for the CDN resource.

Possible values:
- **true** - Origin shield is active.
- **false** - Origin shield is not active.
- `shield_routing_map` (Number) Defines whether the origin shield with a dynamic location is enabled for the CDN resource.

To manage origin shielding, you must contact customer support.
- `shielded` (Boolean) Defines whether origin shielding feature is enabled for the resource.

Possible values:
- **true** - Origin shielding is enabled.
- **false** - Origin shielding is disabled.
- `suspend_date` (String) Date when the CDN resource was suspended automatically if there is no traffic on it for 90 days.

Not specified if the resource was not stopped due to lack of traffic.
- `suspended` (Boolean) Defines whether the CDN resource has been automatically suspended because there was no traffic on it for 90 days.

Possible values:
- **true** - CDN resource is currently automatically suspended.
- **false** - CDN resource is not automatically suspended.

You can enable CDN resource using the `active` field. If there is no traffic on the CDN resource within seven days following activation, it will be suspended again.

To avoid CDN resource suspension due to no traffic, contact technical support.
- `vp_enabled` (Boolean) Defines whether the CDN resource is integrated with the Streaming Platform.

Possible values:
- **true** - CDN resource is configured for Streaming Platform. Changing resource settings can affect its operation.
- **false** - CDN resource is not configured for Streaming Platform.
- `waap_domain_id` (String) The ID of the associated WAAP domain.

<a id="nestedatt--options"></a>
### Nested Schema for `options`

Optional:

- `allowed_http_methods` (Attributes) HTTP methods allowed for content requests from the CDN. (see [below for nested schema](#nestedatt--options--allowed_http_methods))
- `bot_protection` (Attributes) Allows to prevent online services from overloading and ensure your business workflow running smoothly. (see [below for nested schema](#nestedatt--options--bot_protection))
- `brotli_compression` (Attributes) Compresses content with Brotli on the CDN side. CDN servers will request only uncompressed content from the origin.

Notes:

1. CDN only supports "Brotli compression" when the "origin shielding" feature is activated.
2. If a precache server is not active for a CDN resource, no compression occurs, even if the option is enabled.
3. `brotli_compression` is not supported with `fetch_compressed` or `slice` options enabled.
4. `fetch_compressed` option in CDN resource settings overrides `brotli_compression` in rules. If you enabled `fetch_compressed` in CDN resource and want to enable `brotli_compression` in a rule, you must specify `fetch_compressed:false` in the rule. (see [below for nested schema](#nestedatt--options--brotli_compression))
- `browser_cache_settings` (Attributes) Cache expiration time for users browsers in seconds.

Cache expiration time is applied to the following response codes: 200, 201, 204, 206, 301, 302, 303, 304, 307, 308.

Responses with other codes will not be cached. (see [below for nested schema](#nestedatt--options--browser_cache_settings))
- `cors` (Attributes) Enables or disables CORS (Cross-Origin Resource Sharing) header support.

CORS header support allows the CDN to add the Access-Control-Allow-Origin header to a response to a browser. (see [below for nested schema](#nestedatt--options--cors))
- `country_acl` (Attributes) Enables control access to content for specified countries. (see [below for nested schema](#nestedatt--options--country_acl))
- `disable_proxy_force_ranges` (Attributes) Allows 206 responses regardless of the settings of an origin source. (see [below for nested schema](#nestedatt--options--disable_proxy_force_ranges))
- `edge_cache_settings` (Attributes) Cache expiration time for CDN servers.

`value` and `default` fields cannot be used simultaneously. (see [below for nested schema](#nestedatt--options--edge_cache_settings))
- `fastedge` (Attributes) Allows to configure FastEdge app to be called on different request/response phases.

Note: At least one of `on_request_headers`, `on_request_body`, `on_response_headers`, or `on_response_body` must be specified. (see [below for nested schema](#nestedatt--options--fastedge))
- `fetch_compressed` (Attributes) Makes the CDN request compressed content from the origin.

The origin server should support compression. CDN servers will not decompress your content even if a user browser does not accept compression.

Notes:

1. `fetch_compressed` is not supported with `gzipON` or `brotli_compression` or `slice` options enabled.
2. `fetch_compressed` overrides `gzipON` and `brotli_compression` in rule. If you enable it in CDN resource and want to use `gzipON` and `brotli_compression` in a rule, you have to specify `"fetch_compressed": false` in the rule. (see [below for nested schema](#nestedatt--options--fetch_compressed))
- `follow_origin_redirect` (Attributes) Enables redirection from origin.
If the origin server returns a redirect, the option allows the CDN to pull the requested content from the origin server that was returned in the redirect. (see [below for nested schema](#nestedatt--options--follow_origin_redirect))
- `force_return` (Attributes) Applies custom HTTP response codes for CDN content.

The following codes are reserved by our system and cannot be specified in this option: 408, 444, 477, 494, 495, 496, 497, 499. (see [below for nested schema](#nestedatt--options--force_return))
- `forward_host_header` (Attributes) Forwards the Host header from a end-user request to an origin server.

`hostHeader` and `forward_host_header` options cannot be enabled simultaneously. (see [below for nested schema](#nestedatt--options--forward_host_header))
- `grpc_passthrough` (Attributes) Enables gRPC pass-through for the CDN resource. (see [below for nested schema](#nestedatt--options--grpc_passthrough))
- `gzip_on` (Attributes) Compresses content with gzip on the CDN end. CDN servers will request only uncompressed content from the origin.

Notes:

1. Compression with gzip is not supported with `fetch_compressed` or `slice` options enabled.
2. `fetch_compressed` option in CDN resource settings overrides `gzipON` in rules. If you enable `fetch_compressed` in CDN resource and want to enable `gzipON` in rules, you need to specify `"fetch_compressed":false` for rules. (see [below for nested schema](#nestedatt--options--gzip_on))
- `host_header` (Attributes) Sets the Host header that CDN servers use when request content from an origin server.
Your server must be able to process requests with the chosen header.

If the option is `null`, the Host Header value is equal to first CNAME.

`hostHeader` and `forward_host_header` options cannot be enabled simultaneously. (see [below for nested schema](#nestedatt--options--host_header))
- `http3_enabled` (Attributes) Enables HTTP/3 protocol for content delivery.

`http3_enabled` option works only with `"sslEnabled": true`. (see [below for nested schema](#nestedatt--options--http3_enabled))
- `ignore_cookie` (Attributes) Defines whether the files with the Set-Cookies header are cached as one file or as different ones. (see [below for nested schema](#nestedatt--options--ignore_cookie))
- `ignore_query_string` (Attributes) How a file with different query strings is cached: either as one object (option is enabled) or as different objects (option is disabled.)

`ignoreQueryString`, `query_params_whitelist` and `query_params_blacklist` options cannot be enabled simultaneously. (see [below for nested schema](#nestedatt--options--ignore_query_string))
- `image_stack` (Attributes) Transforms JPG and PNG images (for example, resize or crop) and automatically converts them to WebP or AVIF format. (see [below for nested schema](#nestedatt--options--image_stack))
- `ip_address_acl` (Attributes) Controls access to the CDN resource content for specific IP addresses.

If you want to use IPs from our CDN servers IP list for IP ACL configuration, you have to independently monitor their relevance.

We recommend you use a script for automatically update IP ACL. [Read more.](/docs/api-reference/cdn/ip-addresses-list/get-cdn-servers-ip-addresses) (see [below for nested schema](#nestedatt--options--ip_address_acl))
- `limit_bandwidth` (Attributes) Allows to control the download speed per connection. (see [below for nested schema](#nestedatt--options--limit_bandwidth))
- `proxy_cache_key` (Attributes) Allows you to modify your cache key. If omitted, the default value is `$request_uri`.

Combine the specified variables to create a key for caching.
- **$`request_uri`**
- **$scheme**
- **$uri**

**Warning**: Enabling and changing this option can invalidate your current cache and affect the cache hit ratio. Furthermore, the "Purge by pattern" option will not work. (see [below for nested schema](#nestedatt--options--proxy_cache_key))
- `proxy_cache_methods_set` (Attributes) Caching for POST requests along with default GET and HEAD. (see [below for nested schema](#nestedatt--options--proxy_cache_methods_set))
- `proxy_connect_timeout` (Attributes) The time limit for establishing a connection with the origin. (see [below for nested schema](#nestedatt--options--proxy_connect_timeout))
- `proxy_read_timeout` (Attributes) The time limit for receiving a partial response from the origin.
If no response is received within this time, the connection will be closed.

**Note:**
When used with a WebSocket connection, this option supports values only in the range 1–20 seconds (instead of the usual 1–30 seconds). (see [below for nested schema](#nestedatt--options--proxy_read_timeout))
- `query_params_blacklist` (Attributes) Files with the specified query parameters are cached as one object, files with other parameters are cached as different objects.

`ignoreQueryString`, `query_params_whitelist` and `query_params_blacklist` options cannot be enabled simultaneously. (see [below for nested schema](#nestedatt--options--query_params_blacklist))
- `query_params_whitelist` (Attributes) Files with the specified query parameters are cached as different objects, files with other parameters are cached as one object.

`ignoreQueryString`, `query_params_whitelist` and `query_params_blacklist` options cannot be enabled simultaneously. (see [below for nested schema](#nestedatt--options--query_params_whitelist))
- `query_string_forwarding` (Attributes) The Query String Forwarding feature allows for the seamless transfer of parameters embedded in playlist files to the corresponding media chunk files.
This functionality ensures that specific attributes, such as authentication tokens or tracking information, are consistently passed along from the playlist manifest to the individual media segments.
This is particularly useful for maintaining continuity in security, analytics, and any other parameter-based operations across the entire media delivery workflow. (see [below for nested schema](#nestedatt--options--query_string_forwarding))
- `redirect_http_to_https` (Attributes) Enables redirect from HTTP to HTTPS.

`redirect_http_to_https` and `redirect_https_to_http` options cannot be enabled simultaneously. (see [below for nested schema](#nestedatt--options--redirect_http_to_https))
- `redirect_https_to_http` (Attributes) Enables redirect from HTTPS to HTTP.

`redirect_http_to_https` and `redirect_https_to_http` options cannot be enabled simultaneously. (see [below for nested schema](#nestedatt--options--redirect_https_to_http))
- `referrer_acl` (Attributes) Controls access to the CDN resource content for specified domain names. (see [below for nested schema](#nestedatt--options--referrer_acl))
- `request_limiter` (Attributes) Option allows to limit the amount of HTTP requests. (see [below for nested schema](#nestedatt--options--request_limiter))
- `response_headers_hiding_policy` (Attributes) Hides HTTP headers from an origin server in the CDN response. (see [below for nested schema](#nestedatt--options--response_headers_hiding_policy))
- `rewrite` (Attributes) Changes and redirects requests from the CDN to the origin. It operates according to the [Nginx](https://nginx.org/en/docs/http/ngx_http_rewrite_module.html#rewrite) configuration. (see [below for nested schema](#nestedatt--options--rewrite))
- `secure_key` (Attributes) Configures access with tokenized URLs. This makes impossible to access content without a valid (unexpired) token. (see [below for nested schema](#nestedatt--options--secure_key))
- `slice` (Attributes) Requests and caches files larger than 10 MB in parts (no larger than 10 MB per part.) This reduces time to first byte.

The option is based on the [Slice](https://nginx.org/en/docs/http/ngx_http_slice_module.html) module.

Notes:

1. Origin must support HTTP Range requests.
2. Not supported with `gzipON`, `brotli_compression` or `fetch_compressed` options enabled. (see [below for nested schema](#nestedatt--options--slice))
- `sni` (Attributes) The hostname that is added to SNI requests from CDN servers to the origin server via HTTPS.

SNI is generally only required if your origin uses shared hosting or does not have a dedicated IP address.
If the origin server presents multiple certificates, SNI allows the origin server to know which certificate to use for the connection.

The option works only if `originProtocol` parameter is `HTTPS` or `MATCH`. (see [below for nested schema](#nestedatt--options--sni))
- `stale` (Attributes) Serves stale cached content in case of origin unavailability. (see [below for nested schema](#nestedatt--options--stale))
- `static_request_headers` (Attributes) Custom HTTP Headers for a CDN server to add to request. Up to fifty custom HTTP Headers can be specified. (see [below for nested schema](#nestedatt--options--static_request_headers))
- `static_response_headers` (Attributes) Custom HTTP Headers that a CDN server adds to a response. (see [below for nested schema](#nestedatt--options--static_response_headers))
- `tls_versions` (Attributes) List of SSL/TLS protocol versions allowed for HTTPS connections from end users to the domain.

When the option is disabled, all protocols versions are allowed. (see [below for nested schema](#nestedatt--options--tls_versions))
- `use_default_le_chain` (Attributes) Let's Encrypt certificate chain.

The specified chain will be used during the next Let's Encrypt certificate issue or renewal. (see [below for nested schema](#nestedatt--options--use_default_le_chain))
- `use_dns01_le_challenge` (Attributes) DNS-01 challenge to issue a Let's Encrypt certificate for the resource.

DNS service should be activated to enable this option. (see [below for nested schema](#nestedatt--options--use_dns01_le_challenge))
- `use_rsa_le_cert` (Attributes) RSA Let's Encrypt certificate type for the CDN resource.

The specified value will be used during the next Let's Encrypt certificate issue or renewal. (see [below for nested schema](#nestedatt--options--use_rsa_le_cert))
- `user_agent_acl` (Attributes) Controls access to the content for specified User-Agents. (see [below for nested schema](#nestedatt--options--user_agent_acl))
- `waap` (Attributes) Allows to enable WAAP (Web Application and API Protection). (see [below for nested schema](#nestedatt--options--waap))
- `websockets` (Attributes) Enables or disables WebSockets connections to an origin server. (see [below for nested schema](#nestedatt--options--websockets))

<a id="nestedatt--options--allowed_http_methods"></a>
### Nested Schema for `options.allowed_http_methods`

Required:

- `enabled` (Boolean) Controls the option state.

Possible values:
- **true** - Option is enabled.
- **false** - Option is disabled.
- `value` (Set of String)


<a id="nestedatt--options--bot_protection"></a>
### Nested Schema for `options.bot_protection`

Required:

- `bot_challenge` (Attributes) Controls the bot challenge module state. (see [below for nested schema](#nestedatt--options--bot_protection--bot_challenge))
- `enabled` (Boolean) Controls the option state.

Possible values:
- **true** - Option is enabled.
- **false** - Option is disabled.

<a id="nestedatt--options--bot_protection--bot_challenge"></a>
### Nested Schema for `options.bot_protection.bot_challenge`

Optional:

- `enabled` (Boolean) Possible values:
- **true** - Bot challenge is enabled.
- **false** - Bot challenge is disabled.



<a id="nestedatt--options--brotli_compression"></a>
### Nested Schema for `options.brotli_compression`

Required:

- `enabled` (Boolean) Controls the option state.

Possible values:
- **true** - Option is enabled.
- **false** - Option is disabled.
- `value` (Set of String) Allows to select the content types you want to compress.

`text/html` is a mandatory content type.


<a id="nestedatt--options--browser_cache_settings"></a>
### Nested Schema for `options.browser_cache_settings`

Required:

- `enabled` (Boolean) Controls the option state.

Possible values:
- **true** - Option is enabled.
- **false** - Option is disabled.
- `value` (String) Set the cache expiration time to '0s' to disable caching.

The maximum duration is any equivalent to `1y`.


<a id="nestedatt--options--cors"></a>
### Nested Schema for `options.cors`

Required:

- `enabled` (Boolean) Controls the option state.

Possible values:
- **true** - Option is enabled.
- **false** - Option is disabled.
- `value` (Set of String) Value of the Access-Control-Allow-Origin header.

Possible values:
- **Adds * as the Access-Control-Allow-Origin header value** - Content will be uploaded for requests from any domain.
`"value": ["*"]`
- **Adds "$http_origin" as the Access-Control-Allow-Origin header value if the origin matches one of the listed domains** - Content will be uploaded only for requests from the domains specified in the field.
`"value": ["domain.com", "second.dom.com"]`
- **Adds "$http_origin" as the Access-Control-Allow-Origin header value** - Content will be uploaded for requests from any domain, and the domain from which the request was sent will be added to the "Access-Control-Allow-Origin" header in the response.
`"value": ["$http_origin"]`

Optional:

- `always` (Boolean) Defines whether the Access-Control-Allow-Origin header should be added to a response from CDN regardless of response code.

Possible values:
- **true** - Header will be added to a response regardless of response code.
- **false** - Header will only be added to responses with codes: 200, 201, 204, 206, 301, 302, 303, 304, 307, 308.


<a id="nestedatt--options--country_acl"></a>
### Nested Schema for `options.country_acl`

Required:

- `enabled` (Boolean) Controls the option state.

Possible values:
- **true** - Option is enabled.
- **false** - Option is disabled.
- `excepted_values` (Set of String) List of countries according to ISO-3166-1.

The meaning of the parameter depends on `policy_type` value:
- **allow** - List of countries for which access is prohibited.
- **deny** - List of countries for which access is allowed.
- `policy_type` (String) Defines the type of CDN resource access policy.

Possible values:
- **allow** - Access is allowed for all the countries except for those specified in `excepted_values` field.
- **deny** - Access is denied for all the countries except for those specified in `excepted_values` field.
Available values: "allow", "deny".


<a id="nestedatt--options--disable_proxy_force_ranges"></a>
### Nested Schema for `options.disable_proxy_force_ranges`

Required:

- `enabled` (Boolean) Controls the option state.

Possible values:
- **true** - Option is enabled.
- **false** - Option is disabled.
- `value` (Boolean) Possible values:
- **true** - Option is enabled.
- **false** - Option is disabled.


<a id="nestedatt--options--edge_cache_settings"></a>
### Nested Schema for `options.edge_cache_settings`

Required:

- `enabled` (Boolean) Controls the option state.

Possible values:
- **true** - Option is enabled.
- **false** - Option is disabled.

Optional:

- `custom_values` (Map of String) A MAP object representing the caching time in seconds for a response with a specific response code.

These settings have a higher priority than the `value` field.

- Use `any` key to specify caching time for all response codes.
- Use `0s` value to disable caching for a specific response code.
- `default` (String) Enables content caching according to the origin cache settings.

The value is applied to the following response codes 200, 201, 204, 206, 301, 302, 303, 304, 307, 308, if an origin server does not have caching HTTP headers.

Responses with other codes will not be cached.

The maximum duration is any equivalent to `1y`.
- `value` (String) Caching time.

The value is applied to the following response codes: 200, 206, 301, 302.
Responses with codes 4xx, 5xx will not be cached.

Use `0s` to disable caching.

The maximum duration is any equivalent to `1y`.


<a id="nestedatt--options--fastedge"></a>
### Nested Schema for `options.fastedge`

Required:

- `enabled` (Boolean) Controls the option state.

Possible values:
- **true** - Option is enabled.
- **false** - Option is disabled.

Optional:

- `on_request_body` (Attributes) Allows to configure FastEdge application that will be called to handle request body as soon as CDN receives incoming HTTP request. (see [below for nested schema](#nestedatt--options--fastedge--on_request_body))
- `on_request_headers` (Attributes) Allows to configure FastEdge application that will be called to handle request headers as soon as CDN receives incoming HTTP request, **before cache**. (see [below for nested schema](#nestedatt--options--fastedge--on_request_headers))
- `on_response_body` (Attributes) Allows to configure FastEdge application that will be called to handle response body before CDN sends the HTTP response. (see [below for nested schema](#nestedatt--options--fastedge--on_response_body))
- `on_response_headers` (Attributes) Allows to configure FastEdge application that will be called to handle response headers before CDN sends the HTTP response. (see [below for nested schema](#nestedatt--options--fastedge--on_response_headers))

<a id="nestedatt--options--fastedge--on_request_body"></a>
### Nested Schema for `options.fastedge.on_request_body`

Required:

- `app_id` (String) The ID of the application in FastEdge.

Optional:

- `enabled` (Boolean) Determines if the FastEdge application should be called whenever HTTP request headers are received.
- `execute_on_edge` (Boolean) Determines if the request should be executed at the edge nodes.
- `execute_on_shield` (Boolean) Determines if the request should be executed at the shield nodes.
- `interrupt_on_error` (Boolean) Determines if the request execution should be interrupted when an error occurs.


<a id="nestedatt--options--fastedge--on_request_headers"></a>
### Nested Schema for `options.fastedge.on_request_headers`

Required:

- `app_id` (String) The ID of the application in FastEdge.

Optional:

- `enabled` (Boolean) Determines if the FastEdge application should be called whenever HTTP request headers are received.
- `execute_on_edge` (Boolean) Determines if the request should be executed at the edge nodes.
- `execute_on_shield` (Boolean) Determines if the request should be executed at the shield nodes.
- `interrupt_on_error` (Boolean) Determines if the request execution should be interrupted when an error occurs.


<a id="nestedatt--options--fastedge--on_response_body"></a>
### Nested Schema for `options.fastedge.on_response_body`

Required:

- `app_id` (String) The ID of the application in FastEdge.

Optional:

- `enabled` (Boolean) Determines if the FastEdge application should be called whenever HTTP request headers are received.
- `execute_on_edge` (Boolean) Determines if the request should be executed at the edge nodes.
- `execute_on_shield` (Boolean) Determines if the request should be executed at the shield nodes.
- `interrupt_on_error` (Boolean) Determines if the request execution should be interrupted when an error occurs.


<a id="nestedatt--options--fastedge--on_response_headers"></a>
### Nested Schema for `options.fastedge.on_response_headers`

Required:

- `app_id` (String) The ID of the application in FastEdge.

Optional:

- `enabled` (Boolean) Determines if the FastEdge application should be called whenever HTTP request headers are received.
- `execute_on_edge` (Boolean) Determines if the request should be executed at the edge nodes.
- `execute_on_shield` (Boolean) Determines if the request should be executed at the shield nodes.
- `interrupt_on_error` (Boolean) Determines if the request execution should be interrupted when an error occurs.



<a id="nestedatt--options--fetch_compressed"></a>
### Nested Schema for `options.fetch_compressed`

Required:

- `enabled` (Boolean) Controls the option state.

Possible values:
- **true** - Option is enabled.
- **false** - Option is disabled.
- `value` (Boolean) Possible values:
- **true** - Option is enabled.
- **false** - Option is disabled.


<a id="nestedatt--options--follow_origin_redirect"></a>
### Nested Schema for `options.follow_origin_redirect`

Required:

- `codes` (Set of Number) Redirect status code that the origin server returns.

To serve up to date content to end users, you will need to purge the cache after managing the option.
- `enabled` (Boolean) Controls the option state.

Possible values:
- **true** - Option is enabled.
- **false** - Option is disabled.


<a id="nestedatt--options--force_return"></a>
### Nested Schema for `options.force_return`

Required:

- `body` (String) URL for redirection or text.
- `code` (Number) Status code value.
- `enabled` (Boolean) Controls the option state.

Possible values:
- **true** - Option is enabled.
- **false** - Option is disabled.

Optional:

- `time_interval` (Attributes) Controls the time at which a custom HTTP response code should be applied. By default, a custom HTTP response code is applied at any time. (see [below for nested schema](#nestedatt--options--force_return--time_interval))

<a id="nestedatt--options--force_return--time_interval"></a>
### Nested Schema for `options.force_return.time_interval`

Required:

- `end_time` (String) Time until which a custom HTTP response code should be applied. Indicated in 24-hour format.
- `start_time` (String) Time from which a custom HTTP response code should be applied. Indicated in 24-hour format.

Optional:

- `time_zone` (String) Time zone used to calculate time.



<a id="nestedatt--options--forward_host_header"></a>
### Nested Schema for `options.forward_host_header`

Required:

- `enabled` (Boolean) Controls the option state.

Possible values:
- **true** - Option is enabled.
- **false** - Option is disabled.
- `value` (Boolean) Possible values:
- **true** - Option is enabled.
- **false** - Option is disabled.


<a id="nestedatt--options--grpc_passthrough"></a>
### Nested Schema for `options.grpc_passthrough`

Required:

- `enabled` (Boolean) Controls the option state.

Possible values:
- **true** - Option is enabled.
- **false** - Option is disabled.
- `value` (Boolean) Possible values:
- **true** - Option is enabled.
- **false** - Option is disabled.


<a id="nestedatt--options--gzip_on"></a>
### Nested Schema for `options.gzip_on`

Required:

- `enabled` (Boolean) Controls the option state.

Possible values:
- **true** - Option is enabled.
- **false** - Option is disabled.
- `value` (Boolean) Possible values:
- **true** - Option is enabled.
- **false** - Option is disabled.


<a id="nestedatt--options--host_header"></a>
### Nested Schema for `options.host_header`

Required:

- `enabled` (Boolean) Controls the option state.

Possible values:
- **true** - Option is enabled.
- **false** - Option is disabled.
- `value` (String) Host Header value.


<a id="nestedatt--options--http3_enabled"></a>
### Nested Schema for `options.http3_enabled`

Required:

- `enabled` (Boolean) Controls the option state.

Possible values:
- **true** - Option is enabled.
- **false** - Option is disabled.
- `value` (Boolean) Possible values:
- **true** - Option is enabled.
- **false** - Option is disabled.


<a id="nestedatt--options--ignore_cookie"></a>
### Nested Schema for `options.ignore_cookie`

Required:

- `enabled` (Boolean) Controls the option state.

Possible values:
- **true** - Option is enabled.
- **false** - Option is disabled.
- `value` (Boolean) Possible values:
- **true** - Option is enabled, files with cookies are cached as one file.
- **false** - Option is disabled, files with cookies are cached as different files.


<a id="nestedatt--options--ignore_query_string"></a>
### Nested Schema for `options.ignore_query_string`

Required:

- `enabled` (Boolean) Controls the option state.

Possible values:
- **true** - Option is enabled.
- **false** - Option is disabled.
- `value` (Boolean) Possible values:
- **true** - Option is enabled.
- **false** - Option is disabled.


<a id="nestedatt--options--image_stack"></a>
### Nested Schema for `options.image_stack`

Required:

- `enabled` (Boolean) Controls the option state.

Possible values:
- **true** - Option is enabled.
- **false** - Option is disabled.

Optional:

- `avif_enabled` (Boolean) Enables or disables automatic conversion of JPEG and PNG images to AVI format.
- `png_lossless` (Boolean) Enables or disables compression without quality loss for PNG format.
- `quality` (Number) Defines quality settings for JPG and PNG images. The higher the value, the better the image quality, and the larger the file size after conversion.
- `webp_enabled` (Boolean) Enables or disables automatic conversion of JPEG and PNG images to WebP format.


<a id="nestedatt--options--ip_address_acl"></a>
### Nested Schema for `options.ip_address_acl`

Required:

- `enabled` (Boolean) Controls the option state.

Possible values:
- **true** - Option is enabled.
- **false** - Option is disabled.
- `excepted_values` (Set of String) List of IP addresses with a subnet mask.

The meaning of the parameter depends on `policy_type` value:
- **allow** - List of IP addresses for which access is prohibited.
- **deny** - List of IP addresses for which access is allowed.

Examples:
- `192.168.3.2/32`
- `2a03:d000:2980:7::8/128`
- `policy_type` (String) IP access policy type.

Possible values:
- **allow** - Allow access to all IPs except IPs specified in "excepted_values" field.
- **deny** - Deny access to all IPs except IPs specified in "excepted_values" field.
Available values: "allow", "deny".


<a id="nestedatt--options--limit_bandwidth"></a>
### Nested Schema for `options.limit_bandwidth`

Required:

- `enabled` (Boolean) Controls the option state.

Possible values:
- **true** - Option is enabled.
- **false** - Option is disabled.
- `limit_type` (String) Method of controlling the download speed per connection.

Possible values:
- **static** - Use speed and buffer fields to set the download speed limit.
- **dynamic** - Use query strings **speed** and **buffer** to set the download speed limit.

For example, when requesting content at the link

```
http://cdn.example.com/video.mp4?speed=50k&buffer=500k
```

the download speed will be limited to 50kB/s after 500 kB.
Available values: "static", "dynamic".

Optional:

- `buffer` (Number) Amount of downloaded data after which the user will be rate limited.
- `speed` (Number) Maximum download speed per connection.


<a id="nestedatt--options--proxy_cache_key"></a>
### Nested Schema for `options.proxy_cache_key`

Required:

- `enabled` (Boolean) Controls the option state.

Possible values:
- **true** - Option is enabled.
- **false** - Option is disabled.
- `value` (String) Key for caching.


<a id="nestedatt--options--proxy_cache_methods_set"></a>
### Nested Schema for `options.proxy_cache_methods_set`

Required:

- `enabled` (Boolean) Controls the option state.

Possible values:
- **true** - Option is enabled.
- **false** - Option is disabled.
- `value` (Boolean) Possible values:
- **true** - Option is enabled.
- **false** - Option is disabled.


<a id="nestedatt--options--proxy_connect_timeout"></a>
### Nested Schema for `options.proxy_connect_timeout`

Required:

- `enabled` (Boolean) Controls the option state.

Possible values:
- **true** - Option is enabled.
- **false** - Option is disabled.
- `value` (String) Timeout value in seconds.

Supported range: **1s - 5s**.


<a id="nestedatt--options--proxy_read_timeout"></a>
### Nested Schema for `options.proxy_read_timeout`

Required:

- `enabled` (Boolean) Controls the option state.

Possible values:
- **true** - Option is enabled.
- **false** - Option is disabled.
- `value` (String) Timeout value in seconds.

Supported range: **1s - 30s**.


<a id="nestedatt--options--query_params_blacklist"></a>
### Nested Schema for `options.query_params_blacklist`

Required:

- `enabled` (Boolean) Controls the option state.

Possible values:
- **true** - Option is enabled.
- **false** - Option is disabled.
- `value` (Set of String) List of query parameters.


<a id="nestedatt--options--query_params_whitelist"></a>
### Nested Schema for `options.query_params_whitelist`

Required:

- `enabled` (Boolean) Controls the option state.

Possible values:
- **true** - Option is enabled.
- **false** - Option is disabled.
- `value` (Set of String) List of query parameters.


<a id="nestedatt--options--query_string_forwarding"></a>
### Nested Schema for `options.query_string_forwarding`

Required:

- `enabled` (Boolean) Controls the option state.

Possible values:
- **true** - Option is enabled.
- **false** - Option is disabled.
- `forward_from_file_types` (Set of String) The `forward_from_files_types` field specifies the types of playlist files from which parameters will be extracted and forwarded.
This typically includes formats that list multiple media chunk references, such as HLS and DASH playlists.
Parameters associated with these playlist files (like query strings or headers) will be propagated to the chunks they reference.
- `forward_to_file_types` (Set of String) The field specifies the types of media chunk files to which parameters, extracted from playlist files, will be forwarded.
These refer to the actual segments of media content that are delivered to viewers.
Ensuring the correct parameters are forwarded to these files is crucial for maintaining the integrity of the streaming session.

Optional:

- `forward_except_keys` (Set of String) The `forward_except_keys` field provides a mechanism to exclude specific parameters from being forwarded from playlist files to media chunk files.
By listing certain keys in this field, you can ensure that these parameters are omitted during the forwarding process.
This is particularly useful for preventing sensitive or irrelevant information from being included in requests for media chunks, thereby enhancing security and optimizing performance.
- `forward_only_keys` (Set of String) The `forward_only_keys` field allows for granular control over which specific parameters are forwarded from playlist files to media chunk files.
By specifying certain keys, only those parameters will be propagated, ensuring that only relevant information is passed along.
This is particularly useful for security and performance optimization, as it prevents unnecessary or sensitive data from being included in requests for media chunks.


<a id="nestedatt--options--redirect_http_to_https"></a>
### Nested Schema for `options.redirect_http_to_https`

Required:

- `enabled` (Boolean) Controls the option state.

Possible values:
- **true** - Option is enabled.
- **false** - Option is disabled.
- `value` (Boolean) Possible values:
- **true** - Option is enabled.
- **false** - Option is disabled.


<a id="nestedatt--options--redirect_https_to_http"></a>
### Nested Schema for `options.redirect_https_to_http`

Required:

- `enabled` (Boolean) Controls the option state.

Possible values:
- **true** - Option is enabled.
- **false** - Option is disabled.
- `value` (Boolean) Possible values:
- **true** - Option is enabled.
- **false** - Option is disabled.


<a id="nestedatt--options--referrer_acl"></a>
### Nested Schema for `options.referrer_acl`

Required:

- `enabled` (Boolean) Controls the option state.

Possible values:
- **true** - Option is enabled.
- **false** - Option is disabled.
- `excepted_values` (Set of String) List of domain names or wildcard domains (without protocol: `http://` or `https://`.)

The meaning of the parameter depends on `policy_type` value:
- **allow** - List of domain names for which access is prohibited.
- **deny** - List of IP domain names for which access is allowed.

Examples:
- `example.com`
- `*.example.com`
- `policy_type` (String) Policy type.

Possible values:
- **allow** - Allow access to all domain names except the domain names specified in `excepted_values` field.
- **deny** - Deny access to all domain names except the domain names specified in `excepted_values` field.
Available values: "allow", "deny".


<a id="nestedatt--options--request_limiter"></a>
### Nested Schema for `options.request_limiter`

Required:

- `enabled` (Boolean) Controls the option state.

Possible values:
- **true** - Option is enabled.
- **false** - Option is disabled.
- `rate` (Number) Maximum request rate.

Optional:

- `rate_unit` (String) Units of measurement for the `rate` field.

Possible values:
- **r/s** - Requests per second.
- **r/m** - Requests per minute.

If the rate is less than one request per second, it is specified in request per minute (r/m.)
Available values: "r/s", "r/m".

Read-Only:

- `burst` (Number)
- `delay` (Number)


<a id="nestedatt--options--response_headers_hiding_policy"></a>
### Nested Schema for `options.response_headers_hiding_policy`

Required:

- `enabled` (Boolean) Controls the option state.

Possible values:
- **true** - Option is enabled.
- **false** - Option is disabled.
- `mode` (String) How HTTP headers are hidden from the response.

Possible values:
- **show** - Hide only HTTP headers listed in the `excepted` field.
- **hide** - Hide all HTTP headers except headers listed in the "excepted" field.
Available values: "hide", "show".

Optional:

- `excepted` (Set of String) List of HTTP headers.

Parameter meaning depends on the value of the `mode` field:
- **show** - List of HTTP headers to hide from response.
- **hide** - List of HTTP headers to include in response. Other HTTP headers will be hidden.

The following headers are required and cannot be hidden from response:
- `Connection`
- `Content-Length`
- `Content-Type`
- `Date`
- `Server`


<a id="nestedatt--options--rewrite"></a>
### Nested Schema for `options.rewrite`

Required:

- `body` (String) Path for the Rewrite option.

Example:
- `/(.*) /media/$1`
- `enabled` (Boolean) Controls the option state.

Possible values:
- **true** - Option is enabled.
- **false** - Option is disabled.

Optional:

- `flag` (String) Flag for the Rewrite option.

Possible values:
- **last** - Stop processing the current set of `ngx_http_rewrite_module` directives and start a search for a new location matching changed URI.
- **break** - Stop processing the current set of the Rewrite option.
- **redirect** - Return a temporary redirect with the 302 code; used when a replacement string does not start with `http://`, `https://`, or `$scheme`.
- **permanent** - Return a permanent redirect with the 301 code.
Available values: "break", "last", "redirect", "permanent".


<a id="nestedatt--options--secure_key"></a>
### Nested Schema for `options.secure_key`

Required:

- `enabled` (Boolean) Controls the option state.

Possible values:
- **true** - Option is enabled.
- **false** - Option is disabled.
- `key` (String) Key generated on your side that will be used for URL signing.

Optional:

- `type` (Number) Type of URL signing.

Possible types:
- **Type 0** - Includes end user IP to secure token generation.
- **Type 2** - Excludes end user IP from secure token generation.
Available values: 0, 2.


<a id="nestedatt--options--slice"></a>
### Nested Schema for `options.slice`

Required:

- `enabled` (Boolean) Controls the option state.

Possible values:
- **true** - Option is enabled.
- **false** - Option is disabled.
- `value` (Boolean) Possible values:
- **true** - Option is enabled.
- **false** - Option is disabled.


<a id="nestedatt--options--sni"></a>
### Nested Schema for `options.sni`

Required:

- `custom_hostname` (String) Custom SNI hostname.

It is required if `sni_type` is set to custom.
- `enabled` (Boolean) Controls the option state.

Possible values:
- **true** - Option is enabled.
- **false** - Option is disabled.

Optional:

- `sni_type` (String) SNI (Server Name Indication) type.

Possible values:
- **dynamic** - SNI hostname depends on `hostHeader` and `forward_host_header` options.
It has several possible combinations:
- If the `hostHeader` option is enabled and specified, SNI hostname matches the Host header.
- If the `forward_host_header` option is enabled and has true value, SNI hostname matches the Host header used in the request made to a CDN.
- If the `hostHeader` and `forward_host_header` options are disabled, SNI hostname matches the primary CNAME.
- **custom** - custom SNI hostname is in use.
Available values: "dynamic", "custom".


<a id="nestedatt--options--stale"></a>
### Nested Schema for `options.stale`

Required:

- `enabled` (Boolean) Controls the option state.

Possible values:
- **true** - Option is enabled.
- **false** - Option is disabled.
- `value` (Set of String) Defines list of errors for which "Always online" option is applied.


<a id="nestedatt--options--static_request_headers"></a>
### Nested Schema for `options.static_request_headers`

Required:

- `enabled` (Boolean) Controls the option state.

Possible values:
- **true** - Option is enabled.
- **false** - Option is disabled.
- `value` (Map of String) A MAP for static headers in a format of `header_name: header_value`.

Restrictions:
- **Header name** - Maximum 255 symbols, may contain Latin letters (A-Z, a-z), numbers (0-9), dashes, and underscores.
- **Header value** - Maximum 512 symbols, may contain letters (a-z), numbers (0-9), spaces, and symbols (`~!@#%%^&*()-_=+ /|\";:?.,><{}[]). Must start with a letter, number, asterisk or {.


<a id="nestedatt--options--static_response_headers"></a>
### Nested Schema for `options.static_response_headers`

Required:

- `enabled` (Boolean) Controls the option state.

Possible values:
- **true** - Option is enabled.
- **false** - Option is disabled.
- `value` (Attributes List) (see [below for nested schema](#nestedatt--options--static_response_headers--value))

<a id="nestedatt--options--static_response_headers--value"></a>
### Nested Schema for `options.static_response_headers.value`

Required:

- `name` (String) HTTP Header name.

Restrictions:
- Maximum 128 symbols.
- Latin letters (A-Z, a-z,) numbers (0-9,) dashes, and underscores only.
- `value` (Set of String) Header value.

Restrictions:
- Maximum 512 symbols.
- Letters (a-z), numbers (0-9), spaces, and symbols (`~!@#%%^&*()-_=+ /|\";:?.,><{}[]).
- Must start with a letter, number, asterisk or {.
- Multiple values can be added.

Optional:

- `always` (Boolean) Defines whether the header will be added to a response from CDN regardless of response code.

Possible values:
- **true** - Header will be added to a response from CDN regardless of response code.
- **false** - Header will be added only to the following response codes: 200, 201, 204, 206, 301, 302, 303, 304, 307, 308.



<a id="nestedatt--options--tls_versions"></a>
### Nested Schema for `options.tls_versions`

Required:

- `enabled` (Boolean) Controls the option state.

Possible values:
- **true** - Option is enabled.
- **false** - Option is disabled.
- `value` (Set of String) List of SSL/TLS protocol versions (case sensitive).


<a id="nestedatt--options--use_default_le_chain"></a>
### Nested Schema for `options.use_default_le_chain`

Required:

- `enabled` (Boolean) Controls the option state.

Possible values:
- **true** - Option is enabled.
- **false** - Option is disabled.
- `value` (Boolean) Possible values:
- **true** - Default Let's Encrypt certificate chain. This is a deprecated version, use it only for compatibilities with Android devices 7.1.1 or lower.
- **false** - Alternative Let's Encrypt certificate chain.


<a id="nestedatt--options--use_dns01_le_challenge"></a>
### Nested Schema for `options.use_dns01_le_challenge`

Required:

- `enabled` (Boolean) Controls the option state.

Possible values:
- **true** - Option is enabled.
- **false** - Option is disabled.
- `value` (Boolean) Possible values:
- **true** - DNS-01 challenge is used to issue Let's Encrypt certificate.
- **false** - HTTP-01 challenge is used to issue Let's Encrypt certificate.


<a id="nestedatt--options--use_rsa_le_cert"></a>
### Nested Schema for `options.use_rsa_le_cert`

Required:

- `enabled` (Boolean) Controls the option state.

Possible values:
- **true** - Option is enabled.
- **false** - Option is disabled.
- `value` (Boolean) Possible values:
- **true** - RSA Let's Encrypt certificate.
- **false** - ECDSA Let's Encrypt certificate.


<a id="nestedatt--options--user_agent_acl"></a>
### Nested Schema for `options.user_agent_acl`

Required:

- `enabled` (Boolean) Controls the option state.

Possible values:
- **true** - Option is enabled.
- **false** - Option is disabled.
- `excepted_values` (Set of String) List of User-Agents that will be allowed/denied.

The meaning of the parameter depends on `policy_type`:
- **allow** - List of User-Agents for which access is prohibited.
- **deny** - List of User-Agents for which access is allowed.

You can provide exact User-Agent strings or regular expressions. Regular expressions must start
with `~` (case-sensitive) or `~*` (case-insensitive).

Use an empty string `""` to allow/deny access when the User-Agent header is empty.
- `policy_type` (String) User-Agents policy type.

Possible values:
- **allow** - Allow access for all User-Agents except specified in `excepted_values` field.
- **deny** - Deny access for all User-Agents except specified in `excepted_values` field.
Available values: "allow", "deny".


<a id="nestedatt--options--waap"></a>
### Nested Schema for `options.waap`

Required:

- `enabled` (Boolean) Controls the option state.

Possible values:
- **true** - Option is enabled.
- **false** - Option is disabled.
- `value` (Boolean) Possible values:
- **true** - Option is enabled.
- **false** - Option is disabled.


<a id="nestedatt--options--websockets"></a>
### Nested Schema for `options.websockets`

Required:

- `enabled` (Boolean) Controls the option state.

Possible values:
- **true** - Option is enabled.
- **false** - Option is disabled.
- `value` (Boolean) Possible values:
- **true** - Option is enabled.
- **false** - Option is disabled.


## Import

Import is supported using the following syntax:

```shell
$ terraform import gcore_cdn_resource.example '<resource_id>'
```

