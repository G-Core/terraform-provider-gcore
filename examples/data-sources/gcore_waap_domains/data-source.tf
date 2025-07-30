data "gcore_waap_domains" "example_waap_domains" {
  ids = [1]
  name = "name"
  ordering = "id"
  status = "active"
}
