resource "gcore_dns_zone" "example_dns_zone" {
  name = "example.com"
  contact = "contact"
  enabled = true
  expiry = 0
  meta = {
    foo = {

    }
  }
  nx_ttl = 0
  primary_server = "primary_server"
  refresh = 0
  retry = 0
  serial = 0
}
