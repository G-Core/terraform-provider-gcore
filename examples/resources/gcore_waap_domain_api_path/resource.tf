resource "gcore_waap_domain_api_path" "example_waap_domain_api_path" {
  domain_id = 1
  http_scheme = "HTTP"
  method = "GET"
  path = "/api/v1/paths/{path_id}"
  api_groups = ["accounts", "internal"]
  api_version = "v1"
  tags = ["sensitivedataurl", "highriskurl"]
}
