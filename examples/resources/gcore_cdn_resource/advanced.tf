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
