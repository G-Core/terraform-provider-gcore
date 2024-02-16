provider gcore {
  permanent_api_token = "251$d3361.............1b35f26d8"
}

resource "gcore_reservedfixedip" "fixed_ip" {
  project_id = 1
  region_id  = 1
  type       = "external"
  is_vip     = false
}
