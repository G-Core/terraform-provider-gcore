terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

# Test 1: Create automated LE cert via terraform apply
resource "gcore_cdn_certificate" "test" {
  name      = "tf-test-cdn-cert-create"
  automated = true
}

output "cert_id" {
  value = gcore_cdn_certificate.test.id
}

output "cert_ssl_id" {
  value = gcore_cdn_certificate.test.ssl_id
}

output "cert_name" {
  value = gcore_cdn_certificate.test.name
}
