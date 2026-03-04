terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

# Test: Create A record on existing zone
resource "gcore_dns_zone_rrset" "test" {
  zone_name  = "test-minimal-1765280657.dev"
  rrset_name = "www.test-minimal-1765280657.dev"
  rrset_type = "A"
  ttl        = 600  # Changed from 300 to 600
  
  resource_records = [
    {
      content = ["\"1.2.3.4\""]
      enabled = true
    }
  ]
}

output "rrset_info" {
  value = {
    zone_name  = gcore_dns_zone_rrset.test.zone_name
    rrset_name = gcore_dns_zone_rrset.test.rrset_name
    rrset_type = gcore_dns_zone_rrset.test.rrset_type
    name       = gcore_dns_zone_rrset.test.name
    type       = gcore_dns_zone_rrset.test.type
    ttl        = gcore_dns_zone_rrset.test.ttl
  }
}
