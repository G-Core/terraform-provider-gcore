provider gcore {
  permanent_api_token = "251$d3361.............1b35f26d8"
}

resource "gcore_cdn_origingroup" "origin_group_1" {
  name     = "origin_group_1"
  use_next = true
  origin {
    source  = "example.com"
    enabled = true
  }
}

resource "gcore_cdn_resource" "cdn_example_com" {
  cname               = "cdn.example.com"
  origin_group        = gcore_cdn_origingroup.origin_group_1.id
  origin_protocol     = "MATCH"
}

resource "gcore_cdn_rule" "cdn_example_com_rule_1" {
  resource_id = gcore_cdn_resource.cdn_example_com.id
  name        = "Rule with all options"
  rule        = "*.png"
  rule_type   = 0
  weight      = 0

  options {
    allowed_http_methods {
      value = [
        "GET", 
        "POST",
      ]
    }
    brotli_compression {
      value = [
        "text/html", 
        "text/plain",
      ]
    }
    browser_cache_settings {
      value = "3600s"
    }
    cache_http_headers {
      enabled = false
      value = [
        "accept-ranges",
        "content-type",
        "content-encoding",
      ]
    }
    cors {
      value = [
        "*",
      ]
    }
    country_acl {
      policy_type = "allow"
      excepted_values = [
        "GB",
        "DE",
      ]
    }
    disable_cache {
      value = false
    }
    disable_proxy_force_ranges {
      value = true
    }
    edge_cache_settings {
      value = "43200s"
      custom_values = {
        "100" = "400s"
        "101" = "400s"
      }
    }
    fetch_compressed {
      value = false
    }
    follow_origin_redirect {
      codes = [
        301,
        302,
      ]
    }
    force_return {
      enabled = false
      code = 301
      body = "http://example.com/redirect_address"
    }
    forward_host_header {
      value = false
    }
    gzip_on {
      value = true
    }
    host_header {
      value = "host.com"
    }
    ignore_cookie {
      value = true
    }
    ignore_query_string {
      enabled = false
      value = false
    }
    image_stack {
      quality = 80
      avif_enabled = true
      webp_enabled = false
      png_lossless = true
    }
    ip_address_acl {
      policy_type = "deny"
      excepted_values = [
        "192.168.1.100/32"
      ]
    }
    limit_bandwidth {
      limit_type = "static"
      speed = 100
      buffer = 200
    }
    proxy_cache_methods_set {
      value = false
    }
    query_params_blacklist {
      enabled = false
      value = [
        "some",
        "blacklist",
      ]
    }
    query_params_whitelist {
      value = [
        "other",
        "whitelist",
      ]
    }
    redirect_https_to_http {
      enabled = false
      value = false
    }
    redirect_http_to_https {
      enabled = false
      value = true
    }
    referrer_acl {
      policy_type = "deny"
      excepted_values = [
        "*.google.com"
      ]
    }
    request_limiter {
      rate = 5
      burst = 1
      rate_unit = "r/s"
      delay = 0
    }
    response_headers_hiding_policy {
      mode = "hide"
      excepted = [
        "my-header"
      ]
    }
    rewrite {
      body = "/(.*) /additional_path/$1"
      flag = "break"
    }
    secure_key {
      key = "secret"
      type = 2
    }
    slice {
      enabled = false
      value = true
    }
    sni {
      sni_type = "custom"
      custom_hostname = "custom.example.com"
    }
    stale {
      value = [
        "http_404",
        "http_500",
      ]
    }
    static_headers {
      enabled = false
      value = {
        "X-Custom" = "test"
      }
    }
    static_request_headers {
      value = {
        "X-Custom" = "X-Request"
      }
    }
    static_response_headers {
      value {
        name = "X-Custom1"
        value = [
          "Value1",
          "Value2",
        ]
      }
      value {
        name = "X-Custom2"
        value = [
          "CDN"
        ]
        always = true
      }
    }
    user_agent_acl {
      policy_type = "allow"
      excepted_values = [
        "UserAgent"
      ]
    }
    webp {
      enabled = false
      jpg_quality = 55
      png_quality = 66
    }
    websockets {
      value = true
    }
  }
}

