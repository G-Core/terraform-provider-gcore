# New Provider Test Configuration for GCLOUD2-22173
# Testing dns_zone resource with Stainless-generated provider

terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {
  api_key = var.gcore_api_key
}

variable "gcore_api_key" {
  type        = string
  description = "Gcore API Key"
  sensitive   = true
}

variable "zone_name" {
  type        = string
  default     = "tf-test-zone-22173-new.example.com"
  description = "DNS Zone name for testing"
}

# Test 1: Create DNS Zone with basic configuration
resource "gcore_dns_zone" "test_basic" {
  name = var.zone_name

  # Optional SOA fields (all trigger RequiresReplace in new provider!)
  # Note: contact editing is prohibited by tariff, so we omit it
  nx_ttl         = 300
  retry          = 3600
  expiry         = 1209600

  # Zone enabled/disabled
  enabled = true

  # Note: DNSSEC is NOT supported in new provider!
}

# Test 2: DNS Zone RRSet (A record)
resource "gcore_dns_zone_rrset" "test_a" {
  zone_name  = gcore_dns_zone.test_basic.name
  rrset_name = "www.${var.zone_name}"
  rrset_type = "A"
  ttl        = 300

  resource_records = [
    {
      content = ["\"192.168.1.100\""]
      enabled = true
    },
    {
      content = ["\"192.168.1.101\""]
      enabled = true
    }
  ]
}

# Test 3: DNS Zone RRSet (TXT record)
resource "gcore_dns_zone_rrset" "test_txt" {
  zone_name  = gcore_dns_zone.test_basic.name
  rrset_name = "test-txt.${var.zone_name}"
  rrset_type = "TXT"
  ttl        = 600

  resource_records = [
    {
      content = ["\"v=spf1 include:_spf.google.com ~all\""]
      enabled = true
    }
  ]
}

# Test 4: DNS Zone RRSet (CNAME record)
resource "gcore_dns_zone_rrset" "test_cname" {
  zone_name  = gcore_dns_zone.test_basic.name
  rrset_name = "alias.${var.zone_name}"
  rrset_type = "CNAME"
  ttl        = 300

  resource_records = [
    {
      content = ["\"www.${var.zone_name}.\""]
      enabled = true
    }
  ]
}

# Test 5: DNS Zone RRSet (MX record)
resource "gcore_dns_zone_rrset" "test_mx" {
  zone_name  = gcore_dns_zone.test_basic.name
  rrset_name = var.zone_name
  rrset_type = "MX"
  ttl        = 300

  resource_records = [
    {
      content = [10, "\"mail.${var.zone_name}.\""]
      enabled = true
    }
  ]
}

# Outputs for state capture
output "zone_id" {
  value = gcore_dns_zone.test_basic.id
}

output "zone_name" {
  value = gcore_dns_zone.test_basic.name
}

output "zone_details" {
  value = gcore_dns_zone.test_basic.zone
}

output "a_record_name" {
  value = gcore_dns_zone_rrset.test_a.name
}
