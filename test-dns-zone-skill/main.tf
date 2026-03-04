terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

# Main test zone
resource "gcore_dns_zone" "test" {
  name = var.zone_name

  # SOA fields (optional)
  contact        = var.contact
  refresh        = var.refresh
  retry          = var.retry
  expiry         = var.expiry
  nx_ttl         = var.nx_ttl
  primary_server = var.primary_server

  # Meta
  meta = var.meta

  # Enabled
  enabled = var.enabled

  # DNSSEC
  dnssec_enabled = var.dnssec_enabled
}

variable "zone_name" {
  default = "test-tf-dns-skill-1770884578.com"
}

variable "contact" {
  type    = string
  default = null
}

variable "refresh" {
  type    = number
  default = null
}

variable "retry" {
  type    = number
  default = null
}

variable "expiry" {
  type    = number
  default = null
}

variable "nx_ttl" {
  type    = number
  default = null
}

variable "primary_server" {
  type    = string
  default = null
}

variable "meta" {
  type    = map(string)
  default = null
}

variable "enabled" {
  type    = bool
  default = true
}

variable "dnssec_enabled" {
  type    = bool
  default = null
}

# Full zone with all optional fields
variable "create_full_zone" {
  type    = bool
  default = false
}

resource "gcore_dns_zone" "full" {
  count = var.create_full_zone ? 1 : 0

  name           = "test-tf-dns-full-1770884578.com"
  contact        = "admin@example.com"
  refresh        = 7200
  retry          = 1800
  expiry         = 604800
  nx_ttl         = 600
  enabled        = true
  dnssec_enabled = false
  meta = {
    "webhook" = "https://example.com/webhook"
  }
}

# Data source test
variable "test_data_source" {
  type    = bool
  default = false
}

data "gcore_dns_zone" "test" {
  count = var.test_data_source ? 1 : 0
  name  = gcore_dns_zone.test.name
}

output "zone_id" {
  value = gcore_dns_zone.test.id
}

output "zone_serial" {
  value = gcore_dns_zone.test.serial
}

output "zone_dnssec" {
  value = gcore_dns_zone.test.dnssec_enabled
}

output "zone_status" {
  value = gcore_dns_zone.test.status
}

output "zone_contact" {
  value = gcore_dns_zone.test.contact
}

output "zone_records" {
  value = gcore_dns_zone.test.records
}

output "zone_rrsets_amount" {
  value = gcore_dns_zone.test.rrsets_amount
}

output "data_source_zone" {
  value = var.test_data_source ? data.gcore_dns_zone.test[0].zone : null
}
