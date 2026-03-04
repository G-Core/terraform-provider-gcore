# Comprehensive DNS Zone Test - GCLOUD2-22173
# Testing regenerated code with id_property fix

terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

variable "zone_name" {
  default = "tf-test-dns-1765268016.com"
}

variable "enabled" {
  default = true
}

variable "nx_ttl" {
  default = 300
}

variable "retry" {
  default = 3600
}

variable "expiry" {
  default = 1209600
}

# Test 1: Basic zone creation
resource "gcore_dns_zone" "test" {
  name    = var.zone_name
  nx_ttl  = var.nx_ttl
  retry   = var.retry
  expiry  = var.expiry
  enabled = var.enabled
}

# Outputs for verification
output "zone_id" {
  value = gcore_dns_zone.test.id
}

output "zone_name" {
  value = gcore_dns_zone.test.name
}

output "zone_enabled" {
  value = gcore_dns_zone.test.enabled
}

output "zone_warnings" {
  value = gcore_dns_zone.test.warnings
}

output "zone_details" {
  value = gcore_dns_zone.test.zone
}
