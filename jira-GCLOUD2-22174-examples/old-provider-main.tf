# Old Provider: gcore_dns_zone_record
terraform {
  required_providers {
    gcore = {
      source  = "G-Core/gcore"
      version = "~> 0.10"
    }
  }
}

provider "gcore" {
  permanent_api_token = var.api_token
}

variable "api_token" {
  sensitive = true
}

# Simple A record for testing
resource "gcore_dns_zone_record" "test_a" {
  zone   = "maxima.lt"
  domain = "tf-old-provider-test.maxima.lt"
  type   = "A"
  ttl    = 300

  resource_record {
    content = "192.168.100.1"
    enabled = true
  }
}

output "zone" {
  value = gcore_dns_zone_record.test_a.zone
}

output "domain" {
  value = gcore_dns_zone_record.test_a.domain
}

output "type" {
  value = gcore_dns_zone_record.test_a.type
}
