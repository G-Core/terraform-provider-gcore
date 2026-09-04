# Apply a CDN preset to a CDN resource.
#
# Presets are read-only, account-scoped bundles of CDN settings, so look the one
# you want up instead of hardcoding its ID. Only presets whose object_type is
# "CDNResource" can be applied to a CDN resource; "Rule" presets apply to
# gcore_cdn_resource_rule objects.
data "gcore_cdn_presets" "all" {}

locals {
  live_streaming_preset_id = one([
    for preset in data.gcore_cdn_presets.all.items :
    preset.id if preset.object_type == "CDNResource" && preset.name == "LIVE STREAMING"
  ])
}

resource "gcore_cdn_origin_group" "example" {
  name = "origin_group_1"
  sources = [{
    source = "example.com"
  }]
}

resource "gcore_cdn_resource" "example" {
  cname        = "cdn.example.com"
  origin_group = gcore_cdn_origin_group.example.id
}

resource "gcore_cdn_applied_preset" "example" {
  preset_id = local.live_streaming_preset_id
  object_id = gcore_cdn_resource.example.id
}
