provider gcore {
  permanent_api_token = "251$d3361.............1b35f26d8"
}

data "gcore_cdn_preset" "live_streaming" {
  id = 52
  name = "LIVE STREAMING"
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
}

resource "gcore_cdn_applied_preset" "demo_preset" {
  preset_id = data.gcore_cdn_preset.live_streaming.id
  object_id = gcore_cdn_resource.demo_resource.id
}
