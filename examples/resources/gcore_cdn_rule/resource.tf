provider gcore {
  permanent_api_token = "251$d3361.............1b35f26d8"
}

resource "gcore_cdn_rule" "cdn_example_com_rule_1" {
  resource_id = gcore_cdn_resource.cdn_example_com.id
  name        = "All PNG images"
  rule        = "/folder/images/*.png"
  rule_type   = 0
  weight      = 0

  options {
    edge_cache_settings {
      default = "14d"
    }
    browser_cache_settings {
      value = "14d"
    }
    redirect_http_to_https {
      value = true
    }
    request_limiter {
      rate_unit = "r/s"
      rate = 5
      burst = 1
    }
    gzip_on {
      value = true
    }
    cors {
      value = [
        "*"
      ]
    }
    rewrite {
      body = "/(.*) /$1"
    }
    ignore_query_string {
      value = true
    }
  }
}

resource "gcore_cdn_rule" "cdn_example_com_rule_2" {
  resource_id     = gcore_cdn_resource.cdn_example_com.id
  name            = "All JS scripts"
  rule            = "/folder/images/*.js"
  rule_type       = 0
  weight          = 0
  origin_protocol = "HTTP"

  options {
    redirect_http_to_https {
      enabled = false
      value   = true
    }
    gzip_on {
      enabled = false
      value   = true
    }
    query_params_whitelist {
      value = [
        "abc",
      ]
    }
  }
}

resource "gcore_cdn_rule" "cdn_example_com_rule_3" {
  resource_id     = gcore_cdn_resource.cdn_example_com.id
  name            = "Block all png images"
  rule            = "*.png"
  rule_type       = 0
  weight          = 0
  origin_protocol = "HTTP"

  options {
    fastedge {
      on_request_headers {
        enabled = true
        app_id = "1001"
        interrupt_on_error = true
        execute_on_edge    = true
        execute_on_shield  = false
      }
      on_response_headers {
        enabled = true
        app_id = "1002"
        interrupt_on_error = true
      }
    }
    force_return {
      code = 404
      body = "Not found."
    }
  }
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
  secondary_hostnames = ["cdn2.example.com"]
}
