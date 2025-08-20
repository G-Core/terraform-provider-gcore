data "gcore_waap_domain_api_paths" "example_waap_domain_api_paths" {
  domain_id = 1
  api_group = "api_group"
  api_version = "api_version"
  http_scheme = "HTTP"
  ids = ["182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e", "182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"]
  method = "GET"
  ordering = "id"
  path = "path"
  source = "API_DESCRIPTION_FILE"
  status = ["CONFIRMED_API", "POTENTIAL_API"]
}
