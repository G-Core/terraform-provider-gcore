# Test for GCLOUD2-22173: dns_zone issues
# Issue 2: Data source zone = null
# Issue 3: Drift after import

terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

variable "zone_name" {
  default = "tf-test-gcloud2-22173.com"
}

# Resource to create the zone
resource "gcore_dns_zone" "test" {
  name    = var.zone_name
  enabled = true
  nx_ttl  = 300
  retry   = 3600
  expiry  = 1209600
}

# Data source to test zone attribute
data "gcore_dns_zone" "test" {
  name = gcore_dns_zone.test.name
}

# Outputs to verify
output "resource_id" {
  value = gcore_dns_zone.test.id
}

output "resource_name" {
  value = gcore_dns_zone.test.name
}

output "resource_enabled" {
  value = gcore_dns_zone.test.enabled
}

output "resource_warnings" {
  value = gcore_dns_zone.test.warnings
}

output "resource_zone" {
  value = gcore_dns_zone.test.zone
}

# Data source outputs - Issue 2 test
output "datasource_id" {
  value = data.gcore_dns_zone.test.id
}

output "datasource_name" {
  value = data.gcore_dns_zone.test.name
}

output "datasource_zone" {
  description = "This should NOT be null - Issue 2"
  value       = data.gcore_dns_zone.test.zone
}
