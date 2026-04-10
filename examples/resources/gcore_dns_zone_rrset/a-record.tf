# Create an A record set
resource "gcore_dns_zone_rrset" "a_record" {
  zone_name  = "example.com"
  rrset_name = "example.com"
  rrset_type = "A"
  ttl        = 120

  resource_records = [
    {
      content = ["127.0.0.100"]
    },
    {
      content = ["127.0.0.200"]
    },
  ]
}
