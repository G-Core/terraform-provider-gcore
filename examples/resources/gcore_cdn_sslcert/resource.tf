provider gcore {
  permanent_api_token = "251$d3361.............1b35f26d8"
}

variable "cert" {
  type      = string
  sensitive = true
}

variable "private_key" {
  type      = string
  sensitive = true
}

resource "gcore_cdn_sslcert" "cdnopt_cert" {
  name        = "Test cert for cdnopt_bookatest_by"
  cert        = var.cert
  private_key = var.private_key
}

resource "gcore_cdn_sslcert" "lets_encrypt_cert" {
  name        = "Test Let's Encrypt certificate"
  automated   = true
}
