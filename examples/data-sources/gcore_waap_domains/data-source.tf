data "gcore_waap_domains" "example_waap_domains" {
  ids = [1]
  name = "*example.com"
  ordering = "id"
  status = "active"
}
