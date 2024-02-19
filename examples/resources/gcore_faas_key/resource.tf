provider gcore {
  permanent_api_token = "251$d3361.............1b35f26d8"
}

resource "gcore_faas_key" "secret" {
  name = "key-name"
  description = "Keys description"
  project_id = 1
  region_id = 1
}

# To get sensitive value use `terraform output secret`
output "secret" {
  value = resource.gcore_faas_key.secret
  sensitive = true
}
