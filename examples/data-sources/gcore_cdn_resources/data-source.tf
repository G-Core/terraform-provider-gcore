data "gcore_cdn_resources" "example_cdn_resources" {
  active = true
  cname = "cname"
  deleted = true
  enabled = true
  is_primary = true
  max_created = "max_created"
  max_updated = "max_updated"
  min_created = "min_created"
  min_updated = "min_updated"
  name = "name"
  origin_group = 0
  origin_protocol = "HTTP"
  rules = "rules"
  secondary_hostnames = "secondaryHostnames"
  shield_dc = "shield_dc"
  shielded = true
  ssl_data = 0
  ssl_data_in = 0
  ssl_enabled = true
  status = "active"
  suspend = true
  suspended = true
  vp_enabled = true
}
