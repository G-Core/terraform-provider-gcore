terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

# Test 7: validate_root_ca defaults to false when not specified
resource "gcore_cdn_certificate" "test" {
  name            = "tf-test-default-validate"
  ssl_certificate = "-----BEGIN CERTIFICATE-----\nfake\n-----END CERTIFICATE-----"
  ssl_private_key = "-----BEGIN PRIVATE KEY-----\nfake\n-----END PRIVATE KEY-----"
}

output "validate_root_ca" {
  value = gcore_cdn_certificate.test.validate_root_ca
}
