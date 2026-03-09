# Create a CDN resource rule
resource "gcore_cdn_cdn_resource_rule" "example" {
  resource_id = 12345
  name        = "my-cdn-rule"
  rule        = "/assets/*.png"
  rule_type   = 0
  active      = true

  override_origin_protocol = "HTTPS"
}
