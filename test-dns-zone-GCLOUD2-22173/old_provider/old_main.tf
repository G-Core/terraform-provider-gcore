# Old Provider Test Configuration for GCLOUD2-22173
# Testing dns_zone resource compatibility

terraform {
  required_providers {
    gcore = {
      source  = "G-Core/gcore"
      version = "0.32.2"  # Latest old provider version
    }
  }
}

provider "gcore" {
  permanent_api_token = var.gcore_api_key
}

variable "gcore_api_key" {
  type        = string
  description = "Gcore API Key"
  sensitive   = true
}

variable "zone_name" {
  type        = string
  default     = "tf-test-zone-22173.example.com"
  description = "DNS Zone name for testing"
}

# Test 1: Create DNS Zone with basic configuration
resource "gcore_dns_zone" "test_basic" {
  name = var.zone_name

  # Optional SOA fields - note: contact editing is prohibited by tariff
  # contact   = "admin@example.com"  # API defaults to support@gcore.com
  nx_ttl    = 300
  # refresh defaults to 0 by API
  retry     = 3600
  expiry    = 1209600

  # Zone enabled/disabled
  enabled = true

  # DNSSEC toggle (old provider feature)
  dnssec = false
}

# Test 2: DNS RRSet (A record)
resource "gcore_dns_zone_record" "test_a" {
  zone   = gcore_dns_zone.test_basic.name
  domain = "www.${var.zone_name}"
  type   = "A"
  ttl    = 300

  resource_record {
    content = "192.168.1.100"
    enabled = true
  }

  resource_record {
    content = "192.168.1.101"
    enabled = true
  }
}

# Test 3: DNS RRSet (TXT record)
resource "gcore_dns_zone_record" "test_txt" {
  zone   = gcore_dns_zone.test_basic.name
  domain = "test-txt.${var.zone_name}"
  type   = "TXT"
  ttl    = 600

  resource_record {
    content = "v=spf1 include:_spf.google.com ~all"
    enabled = true
  }
}

# Test 4: DNS RRSet (CNAME record)
resource "gcore_dns_zone_record" "test_cname" {
  zone   = gcore_dns_zone.test_basic.name
  domain = "alias.${var.zone_name}"
  type   = "CNAME"
  ttl    = 300

  resource_record {
    content = "www.${var.zone_name}."
    enabled = true
  }
}

# Test 5: DNS RRSet (MX record)
resource "gcore_dns_zone_record" "test_mx" {
  zone   = gcore_dns_zone.test_basic.name
  domain = var.zone_name
  type   = "MX"
  ttl    = 300

  resource_record {
    content = "10 mail.${var.zone_name}."
    enabled = true
  }
}

# Outputs for state capture
output "zone_id" {
  value = gcore_dns_zone.test_basic.id
}

output "zone_name" {
  value = gcore_dns_zone.test_basic.name
}

output "a_record_id" {
  value = gcore_dns_zone_record.test_a.id
}

output "txt_record_id" {
  value = gcore_dns_zone_record.test_txt.id
}
