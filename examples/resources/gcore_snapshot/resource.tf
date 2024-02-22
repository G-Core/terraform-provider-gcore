provider gcore {
  permanent_api_token = "251$d3361.............1b35f26d8"
}

resource "gcore_snapshot" "snapshot" {
  project_id  = 1
  region_id   = 1
  name        = "snapshot example"
  volume_id   = "28e9edcb-1593-41fe-971b-da729c6ec301"
  description = "snapshot example description"
  metadata    = {
    env = "test"
  }
}
