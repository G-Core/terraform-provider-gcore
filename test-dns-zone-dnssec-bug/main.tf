terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

resource "gcore_dns_zone" "test" {
  name           = "test-tf-dnssec-bug-1770883208.com"
  dnssec_enabled = var.dnssec_enabled
}

variable "dnssec_enabled" {
  type    = bool
  default = false
}

output "zone_name" {
  value = gcore_dns_zone.test.name
}

output "zone_dnssec" {
  value = gcore_dns_zone.test.dnssec_enabled
}
