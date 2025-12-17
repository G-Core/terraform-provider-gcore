provider gcore {
  permanent_api_token = "251$d3361.............1b35f26d8"
}

resource "gcore_storage_s3" "example_s3" {
  name     = "example"
  location = "s-region-1"
}
