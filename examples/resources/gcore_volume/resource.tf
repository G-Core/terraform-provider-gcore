provider gcore {
  permanent_api_token = "251$d3361.............1b35f26d8"
}

resource "gcore_volume" "volume" {
  name       = "volume_example"
  type_name  = "standard"
  size       = 1
  region_id  = 1
  project_id = 1
}
