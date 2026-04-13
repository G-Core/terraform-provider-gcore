resource "gcore_dns_zone_rrset" "example_dns_zone_rrset" {
  zone_name  = "zoneName"
  rrset_name = "rrsetName"
  rrset_type = "rrsetType"
  resource_records = [{
    content = ["value"]
    enabled = true
    meta = {
      foo = "bar"
    }
  }]
  meta = {
    key = "value"
  }
  pickers = [{
    type   = "geodns"
    limit  = 0
    strict = true
  }]
  ttl = 0
}
