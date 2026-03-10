resource "gcore_cdn_resource" "example_cdn_resource" {
  cname = "cdn.site.com"
  active = true
  description = "My resource"
  name = "Resource for images"
  options = {
    allowed_http_methods = {
      enabled = true
      value = ["GET", "POST"]
    }
    bot_protection = {
      bot_challenge = {
        enabled = true
      }
      enabled = true
    }
    brotli_compression = {
      enabled = true
      value = ["text/html", "text/plain"]
    }
    browser_cache_settings = {
      enabled = true
      value = "3600s"
    }
    cache_http_headers = {
      enabled = false
      value = ["vary", "content-length", "last-modified", "connection", "accept-ranges", "content-type", "content-encoding", "etag", "cache-control", "expires", "keep-alive", "server"]
    }
    cors = {
      enabled = true
      value = ["domain.com", "domain2.com"]
      always = true
    }
    country_acl = {
      enabled = true
      excepted_values = ["GB", "DE"]
      policy_type = "allow"
    }
    disable_cache = {
      enabled = true
      value = false
    }
    disable_proxy_force_ranges = {
      enabled = true
      value = true
    }
    edge_cache_settings = {
      enabled = true
      custom_values = {
        "100" = "43200s"
      }
      default = "321669910225"
      value = "43200s"
    }
    fastedge = {
      enabled = true
      on_request_body = {
        app_id = "1001"
        enabled = true
        execute_on_edge = true
        execute_on_shield = false
        interrupt_on_error = true
      }
      on_request_headers = {
        app_id = "1001"
        enabled = true
        execute_on_edge = true
        execute_on_shield = false
        interrupt_on_error = true
      }
      on_response_body = {
        app_id = "1001"
        enabled = true
        execute_on_edge = true
        execute_on_shield = false
        interrupt_on_error = true
      }
      on_response_headers = {
        app_id = "1001"
        enabled = true
        execute_on_edge = true
        execute_on_shield = false
        interrupt_on_error = true
      }
    }
    fetch_compressed = {
      enabled = true
      value = false
    }
    follow_origin_redirect = {
      codes = [302, 308]
      enabled = true
    }
    force_return = {
      body = "http://example.com/redirect_address"
      code = 301
      enabled = true
      time_interval = {
        end_time = "20:00"
        start_time = "09:00"
        time_zone = "CET"
      }
    }
    forward_host_header = {
      enabled = false
      value = false
    }
    grpc_passthrough = {
      enabled = true
      value = true
    }
    gzip_on = {
      enabled = true
      value = true
    }
    host_header = {
      enabled = true
      value = "host.com"
    }
    http3_enabled = {
      enabled = true
      value = true
    }
    ignore_cookie = {
      enabled = true
      value = true
    }
    ignore_query_string = {
      enabled = true
      value = false
    }
    image_stack = {
      enabled = true
      avif_enabled = true
      png_lossless = true
      quality = 80
      webp_enabled = false
    }
    ip_address_acl = {
      enabled = true
      excepted_values = ["192.168.1.100/32"]
      policy_type = "deny"
    }
    limit_bandwidth = {
      enabled = true
      limit_type = "static"
      buffer = 200
      speed = 100
    }
    proxy_cache_key = {
      enabled = true
      value = "$scheme$uri"
    }
    proxy_cache_methods_set = {
      enabled = true
      value = false
    }
    proxy_connect_timeout = {
      enabled = true
      value = "4s"
    }
    proxy_read_timeout = {
      enabled = true
      value = "10s"
    }
    query_params_blacklist = {
      enabled = true
      value = ["some", "blacklisted", "query"]
    }
    query_params_whitelist = {
      enabled = true
      value = ["some", "whitelisted", "query"]
    }
    query_string_forwarding = {
      enabled = true
      forward_from_file_types = ["m3u8", "mpd"]
      forward_to_file_types = ["ts", "mp4"]
      forward_except_keys = ["debug_info"]
      forward_only_keys = ["auth_token", "session_id"]
    }
    redirect_http_to_https = {
      enabled = true
      value = true
    }
    redirect_https_to_http = {
      enabled = false
      value = true
    }
    referrer_acl = {
      enabled = true
      excepted_values = ["example.com", "*.example.net"]
      policy_type = "deny"
    }
    request_limiter = {
      enabled = true
      rate = 5
      rate_unit = "r/s"
    }
    response_headers_hiding_policy = {
      enabled = true
      excepted = ["my-header"]
      mode = "hide"
    }
    rewrite = {
      body = "/(.*) /additional_path/$1"
      enabled = true
      flag = "break"
    }
    secure_key = {
      enabled = true
      key = "secretkey"
      type = 2
    }
    slice = {
      enabled = true
      value = true
    }
    sni = {
      custom_hostname = "custom.example.com"
      enabled = true
      sni_type = "custom"
    }
    stale = {
      enabled = true
      value = ["http_404", "http_500"]
    }
    static_response_headers = {
      enabled = true
      value = [{
        name = "X-Example"
        value = ["Value_1"]
        always = true
      }, {
        name = "X-Example-Multiple"
        value = ["Value_1", "Value_2", "Value_3"]
        always = false
      }]
    }
    static_headers = {
      enabled = true
      value = {
        X-Example = "Value_1"
        X-Example-Multiple = ["Value_2", "Value_3"]
      }
    }
    static_request_headers = {
      enabled = true
      value = {
        Header-One = "Value 1"
        Header-Two = "Value 2"
      }
    }
    tls_versions = {
      enabled = true
      value = ["SSLv3", "TLSv1.3"]
    }
    use_default_le_chain = {
      enabled = true
      value = true
    }
    use_dns01_le_challenge = {
      enabled = true
      value = true
    }
    use_rsa_le_cert = {
      enabled = true
      value = true
    }
    user_agent_acl = {
      enabled = true
      excepted_values = ["UserAgent Value", "~*.*bot.*", ""]
      policy_type = "allow"
    }
    waap = {
      enabled = true
      value = true
    }
    websockets = {
      enabled = true
      value = true
    }
  }
  origin = "example.com"
  origin_group = 132
  origin_protocol = "HTTPS"
  primary_resource = null
  proxy_ssl_ca = null
  proxy_ssl_data = null
  proxy_ssl_enabled = false
  secondary_hostnames = ["first.example.com", "second.example.com"]
  ssl_data = 192
  ssl_enabled = false
  waap_api_domain_enabled = true
}
