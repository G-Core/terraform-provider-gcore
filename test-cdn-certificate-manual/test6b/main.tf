terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

# Test 6b: Only ssl_private_key without ssl_certificate — should fail validation
resource "gcore_cdn_certificate" "test" {
  name            = "tf-test-key-only"
  ssl_private_key = "-----BEGIN PRIVATE KEY-----\nfake\n-----END PRIVATE KEY-----"
}
