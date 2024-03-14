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
  secondary_hostnames = ["cdn2.example.com"]

  options {
    edge_cache_settings {
      default = "8d"
    }
    browser_cache_settings {
      value = "1d"
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

    tls_versions {
      enabled = true
      value = [
        "TLSv1.2",
      ]
    }

    force_return {
      code = 200
      body = "OK"
    }
  }
}

data "gcore_cdn_shielding_location" "sl" {
  datacenter = "am3"
}

resource "gcore_cdn_originshielding" "origin_shielding_1" {
  resource_id   = gcore_cdn_resource.cdn_example_com.id
  shielding_pop = data.gcore_cdn_shielding_location.sl.id
}