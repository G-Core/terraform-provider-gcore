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
