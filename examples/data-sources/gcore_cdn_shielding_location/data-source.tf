provider gcore {
  permanent_api_token = "251$d3361.............1b35f26d8"
}

data "gcore_cdn_shielding_location" "sl" {
  city = "Luxembourg"
}

resource "gcore_cdn_originshielding" "origin_shielding_1" {
  resource_id   = 1
  shielding_pop = data.gcore_cdn_shielding_location.sl.id
}