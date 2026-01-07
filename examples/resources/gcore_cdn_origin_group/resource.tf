resource "gcore_cdn_origin_group" "example_cdn_origin_group" {
  name = "YourOriginGroup"
  sources = [{
    backup = false
    enabled = true
    source = "yourwebsite.com"
  }, {
    backup = true
    enabled = true
    source = "1.2.3.4:5500"
  }]
  auth_type = "none"
  proxy_next_upstream = ["error", "timeout", "invalid_header", "http_500", "http_502", "http_503", "http_504"]
  use_next = true
}
