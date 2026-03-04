# DNS Zone RRSet Comprehensive Test
# Testing gcore_dns_zone_rrset resource after fixes

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

# Test: A record
resource "gcore_dns_zone_rrset" "test_a" {
  zone_name  = "maxima.lt"
  rrset_name = "tf-test-fixed.maxima.lt"
  rrset_type = "A"
  ttl        = var.ttl

  resource_records = var.resource_records
}

variable "ttl" {
  default = 300
}

variable "resource_records" {
  default = [
    {
      content = ["\"192.168.1.100\""]
      enabled = true
    }
  ]
}

# Outputs for verification
output "rrset_name" {
  value = gcore_dns_zone_rrset.test_a.name
}

output "rrset_type" {
  value = gcore_dns_zone_rrset.test_a.type
}

output "rrset_ttl" {
  value = gcore_dns_zone_rrset.test_a.ttl
}

output "resource_record_ids" {
  value = gcore_dns_zone_rrset.test_a.resource_records != null ? [for r in gcore_dns_zone_rrset.test_a.resource_records : r.id] : []
}
