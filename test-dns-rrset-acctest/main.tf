terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

# Parent zone
resource "gcore_dns_zone" "test" {
  name = "tf-test-rrset-20260218.com"
}

# Test 1: Basic A record — updated TTL + added second record
resource "gcore_dns_zone_rrset" "basic" {
  zone_name  = gcore_dns_zone.test.name
  rrset_name = "www.tf-test-rrset-20260218.com."
  rrset_type = "A"
  ttl        = 600

  resource_records = [
    {
      content = ["\"192.168.1.1\""]
      enabled = true
    },
    {
      content = ["\"192.168.1.2\""]
      enabled = true
    },
  ]
}

# Test 2: Multiple records — unchanged
resource "gcore_dns_zone_rrset" "multi" {
  zone_name  = gcore_dns_zone.test.name
  rrset_name = "multi.tf-test-rrset-20260218.com."
  rrset_type = "A"
  ttl        = 300

  resource_records = [
    {
      content = ["\"10.0.0.1\""]
      enabled = true
    },
    {
      content = ["\"10.0.0.2\""]
      enabled = true
    },
    {
      content = ["\"10.0.0.3\""]
      enabled = true
    },
  ]
}

# Test 3: Records with meta — updated latlong values
resource "gcore_dns_zone_rrset" "with_meta" {
  zone_name  = gcore_dns_zone.test.name
  rrset_name = "geo.tf-test-rrset-20260218.com."
  rrset_type = "A"
  ttl        = 300

  resource_records = [
    {
      content = ["\"192.168.1.1\""]
      enabled = true
      meta = {
        latlong = "[52.0,32.0]"
        asn     = "[12345]"
        ip      = "[\"192.168.1.0/24\"]"
      }
    },
    {
      content = ["\"10.0.0.1\""]
      enabled = true
      meta = {
        latlong = "[41.0,21.0]"
        asn     = "[67890]"
      }
    },
  ]

  pickers = [
    {
      type = "geodns"
    },
    {
      type = "default"
    },
  ]
}
