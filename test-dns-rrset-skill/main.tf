terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

variable "zone_name" {
  default = "test-rrset-skill-20260211.com"
}

variable "ttl" {
  default = 300
}

variable "test_mx" {
  default = false
}

variable "test_pickers" {
  default = false
}

# First create the zone
resource "gcore_dns_zone" "test" {
  name    = var.zone_name
  enabled = true
}

# Test: A record
resource "gcore_dns_zone_rrset" "a_record" {
  zone_name  = gcore_dns_zone.test.name
  rrset_name = "www.${var.zone_name}"
  rrset_type = "A"
  ttl        = var.ttl

  resource_records = [
    {
      content = ["\"192.168.1.1\""]
      enabled = true
    }
  ]
}

# Test: MX record with meta
resource "gcore_dns_zone_rrset" "mx_record" {
  count = var.test_mx ? 1 : 0

  zone_name  = gcore_dns_zone.test.name
  rrset_name = var.zone_name
  rrset_type = "MX"
  ttl        = 3600

  meta = {
    geodns_link = var.zone_name
  }

  resource_records = [
    {
      content = [10, "\"mail1.${var.zone_name}\""]
      enabled = true
      meta = {
        countries = "[\"us\",\"ca\"]"
      }
    },
    {
      content = [20, "\"mail2.${var.zone_name}\""]
      enabled = true
    }
  ]
}

# Test: Record with pickers
resource "gcore_dns_zone_rrset" "with_pickers" {
  count = var.test_pickers ? 1 : 0

  zone_name  = gcore_dns_zone.test.name
  rrset_name = "geo.${var.zone_name}"
  rrset_type = "A"
  ttl        = 300

  resource_records = [
    {
      content = ["\"10.0.0.1\""]
      enabled = true
      meta = {
        countries = "[\"us\"]"
      }
    },
    {
      content = ["\"10.0.0.2\""]
      enabled = true
      meta = {
        countries = "[\"de\"]"
      }
    }
  ]

  pickers = [
    {
      type = "country"
    }
  ]
}

output "a_record_name" {
  value = gcore_dns_zone_rrset.a_record.name
}

output "a_record_type" {
  value = gcore_dns_zone_rrset.a_record.type
}
