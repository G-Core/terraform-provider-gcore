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
    cors {
      value = [
        "*",
      ]
      always = true
    }
    country_acl {
      policy_type = "allow"
      excepted_values = [
        "GB",
        "DE",
      ]
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
    fastedge {
      on_request_headers {
        enabled = true
        app_id = "1001"
        interrupt_on_error = true
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
    http3_enabled {
      enabled = false
      value = true
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
    proxy_cache_key {
      value = "$scheme$uri"
    }
    proxy_cache_methods_set {
      value = false
    }
    proxy_connect_timeout {
      enabled = true
      value = "1s"
    }
    proxy_read_timeout {
      enabled = true
      value = "1s"
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
    query_string_forwarding {
      forward_from_file_types = ["m3u8", "m3u"]
      forward_to_file_types = ["m3u8", "m3u", "ts", "m4s"]
      forward_only_keys = ["auth_token", "session_id"]
      forward_except_keys = ["debug_info"]
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
        "connection",
        "content-length",
        "content-type",
        "date",
        "server",
        "my-header",
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
    tls_versions {
      value = [
        "TLSv1.2",
      ]
    }
    use_default_le_chain {
      value = false
    }
    user_agent_acl {
      policy_type = "allow"
      excepted_values = [
        "UserAgent"
      ]
    }
    use_rsa_le_cert {
      value = true
    }
    websockets {
      value = true
    }
  }
}
