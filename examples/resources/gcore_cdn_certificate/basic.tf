# Upload an SSL certificate for CDN
resource "gcore_cdn_certificate" "website_cert" {
  name             = "website-certificate"
  ssl_certificate  = file("cert.pem")
  ssl_private_key  = file("key.pem")
  validate_root_ca = true
}
