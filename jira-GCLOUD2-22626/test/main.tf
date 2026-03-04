terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

variable "cert_name" {
  default = "tf-test-cdn-cert-gcloud2-22626"
}

# Test 1: Create automated LE cert
resource "gcore_cdn_certificate" "test" {
  name      = var.cert_name
  automated = true
}

# Test 7: Data source read (enabled after initial create)
data "gcore_cdn_certificate" "test" {
  ssl_id = gcore_cdn_certificate.test.ssl_id
}

# Outputs to verify computed fields
output "resource_id" {
  value = gcore_cdn_certificate.test.id
}
output "resource_ssl_id" {
  value = gcore_cdn_certificate.test.ssl_id
}
output "resource_name" {
  value = gcore_cdn_certificate.test.name
}
output "resource_automated" {
  value = gcore_cdn_certificate.test.automated
}
output "resource_cert_issuer" {
  value = gcore_cdn_certificate.test.cert_issuer
}
output "resource_cert_subject_cn" {
  value = gcore_cdn_certificate.test.cert_subject_cn
}
output "resource_cert_subject_alt" {
  value = gcore_cdn_certificate.test.cert_subject_alt
}
output "resource_validity_not_after" {
  value = gcore_cdn_certificate.test.validity_not_after
}
output "resource_validity_not_before" {
  value = gcore_cdn_certificate.test.validity_not_before
}
output "resource_has_related_resources" {
  value = gcore_cdn_certificate.test.has_related_resources
}
output "resource_validate_root_ca" {
  value = gcore_cdn_certificate.test.validate_root_ca
}
output "resource_deleted" {
  value = gcore_cdn_certificate.test.deleted
}

# Data source outputs
output "ds_name" {
  value = data.gcore_cdn_certificate.test.name
}
output "ds_automated" {
  value = data.gcore_cdn_certificate.test.automated
}
output "ds_cert_issuer" {
  value = data.gcore_cdn_certificate.test.cert_issuer
}
output "ds_cert_subject_cn" {
  value = data.gcore_cdn_certificate.test.cert_subject_cn
}
output "ds_id" {
  value = data.gcore_cdn_certificate.test.id
}
