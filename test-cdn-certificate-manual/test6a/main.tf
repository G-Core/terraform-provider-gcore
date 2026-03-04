terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

# Test 6a: Only ssl_certificate without ssl_private_key — should fail validation
resource "gcore_cdn_certificate" "test" {
  name            = "tf-test-cert-only"
  ssl_certificate = "-----BEGIN CERTIFICATE-----\nfake\n-----END CERTIFICATE-----"
}
