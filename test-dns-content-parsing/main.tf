terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

# Test: A Record
resource "gcore_dns_zone_rrset" "test_a" {
  zone_name  = "maxima.lt"
  rrset_name = "tf-content-test-a.maxima.lt"
  rrset_type = "A"
  ttl        = var.ttl_a

  resource_records = [
    {
      content = ["\"192.168.1.100\""]
      enabled = true
    }
  ]
}

variable "ttl_a" {
  default = 300
}

output "a_record" {
  value = {
    name    = gcore_dns_zone_rrset.test_a.name
    content = gcore_dns_zone_rrset.test_a.resource_records[0].content
    ttl     = gcore_dns_zone_rrset.test_a.ttl
  }
}
