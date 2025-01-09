provider gcore {
  permanent_api_token = "251$d3361.............1b35f26d8"
}

variable "cert" {
  type      = string
  sensitive = true
}

resource "gcore_cdn_cacert" "cdnopt_cert" {
  name        = "Test CA cert"
  cert        = var.cert
}
