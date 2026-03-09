# Upload a trusted CA certificate for CDN origin verification
resource "gcore_cdn_trusted_ca_certificate" "example" {
  name            = "My CA Certificate"
  ssl_certificate = file("${path.module}/ca-cert.pem")
}
