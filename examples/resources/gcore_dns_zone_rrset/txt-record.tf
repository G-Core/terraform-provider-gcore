# Create a TXT record on a subdomain
resource "gcore_dns_zone_rrset" "txt_record" {
  zone_name  = "example.com"
  rrset_name = "subdomain.example.com"
  rrset_type = "TXT"
  ttl        = 120

  resource_records = [{
    content = ["v=spf1 include:_spf.google.com ~all"]
    enabled = true
  }]

  pickers = [{
    type   = "geodns"
    limit  = 1
    strict = true
  }]
}
