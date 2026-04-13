resource "gcore_cdn_origin_group" "example_cdn_origin_group" {
  name = "YourOriginGroup"
  sources = [{
    source  = "yourwebsite.com"
    backup  = false
    enabled = true
    tag     = "default"
  }]
  proxy_next_upstream = ["error", "timeout", "invalid_header", "http_500", "http_502", "http_503", "http_504"]
  use_next            = true
}
