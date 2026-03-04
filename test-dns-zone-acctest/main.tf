terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

# Test 1: Basic zone creation
resource "gcore_dns_zone" "basic" {
  name = "tf-test-basic-20260218.com"
}

# Test 2: Zone with meta (MetaStringType)
resource "gcore_dns_zone" "with_meta" {
  name = "tf-test-meta-20260218.com"
  meta = {
    webhook        = "https://example.com/webhook-v2"
    webhook_method = "PATCH"
  }
}

# Test 3: Zone with DNSSEC
resource "gcore_dns_zone" "with_dnssec" {
  name           = "tf-test-dnssec-20260218.com"
  dnssec_enabled = false
}
