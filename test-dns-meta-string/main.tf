terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

# Step 1: Create a DNS zone with meta (plain string, no jsonencode needed)
resource "gcore_dns_zone" "test" {
  name = "tf-meta-test-20260218.com"

  meta = {
    webhook        = "https://updated-example.com/webhook-v3"
    webhook_method = "PATCH"
  }
}

# Step 2: Update rrset - change TTL and meta
resource "gcore_dns_zone_rrset" "test_a" {
  zone_name  = gcore_dns_zone.test.name
  rrset_name = "www.tf-meta-test-20260218.com."
  rrset_type = "A"
  ttl        = 600

  resource_records = [
    {
      content = ["\"192.168.1.1\""]
      enabled = true
      meta = {
        latlong = "[51.5,31.5]"
        asn     = "[12345]"
        ip      = "[\"192.168.1.0/24\"]"
      }
    },
    {
      content = ["\"10.0.0.1\""]
      enabled = true
      meta = {
        latlong = "[40.0,20.0]"
        asn     = "[67890]"
      }
    },
  ]
}
