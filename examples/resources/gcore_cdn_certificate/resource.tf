# Example: Own SSL certificate
resource "gcore_cdn_certificate" "own_cert" {
  name = "My SSL Certificate"

  # Write-only — values are sent to the API but never stored in state
  ssl_certificate_wo = file("cert.pem")
  ssl_private_key_wo = file("key.pem")

  # Increment to force re-send of cert/key
  ssl_certificate_wo_version = 1

  validate_root_ca = true
}

# Example: Automated Let's Encrypt certificate
resource "gcore_cdn_certificate" "automated_cert" {
  name      = "Auto LE Certificate"
  automated = true
}
