# New Provider: gcore_dns_zone_rrset
terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {
  # Uses GCORE_API_KEY from environment
}

# Simple A record for testing
resource "gcore_dns_zone_rrset" "test_a" {
  zone_name  = "maxima.lt"
  rrset_name = "tf-new-provider-test.maxima.lt"
  rrset_type = "A"
  ttl        = 300

  resource_records = [
    {
      content = ["\"192.168.100.1\""]
      enabled = true
    }
  ]
}

output "zone_name" {
  value = gcore_dns_zone_rrset.test_a.zone_name
}

output "rrset_name" {
  value = gcore_dns_zone_rrset.test_a.rrset_name
}

output "rrset_type" {
  value = gcore_dns_zone_rrset.test_a.rrset_type
}

output "name" {
  value = gcore_dns_zone_rrset.test_a.name
}

output "type" {
  value = gcore_dns_zone_rrset.test_a.type
}
