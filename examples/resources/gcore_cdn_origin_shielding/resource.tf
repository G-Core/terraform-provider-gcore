resource "gcore_cdn_origin_shielding" "example_cdn_origin_shielding" {
  resource_id   = 0 # ID of the CDN resource to protect
  shielding_pop = 0 # origin shielding location ID; look it up via the gcore_cdn_origin_shielding data source or GET /cdn/shieldingpop_v2
}
