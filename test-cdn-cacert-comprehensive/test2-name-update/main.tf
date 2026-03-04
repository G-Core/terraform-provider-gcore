terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

variable "cert_name" {
  type    = string
  default = "tf-test-cacert-update-v1"
}

resource "gcore_cdn_trusted_ca_certificate" "test" {
  name            = var.cert_name
  ssl_certificate = <<-EOT
-----BEGIN CERTIFICATE-----
MIIDTTCCAjWgAwIBAgIUSxFCl2F6hwTdsUxcSyXfA8UNBd4wDQYJKoZIhvcNAQEL
BQAwNjEfMB0GA1UEAwwWdGYtdGVzdC1jYS5leGFtcGxlLmNvbTETMBEGA1UECgwK
VEYgVGVzdCBDQTAeFw0yNjAyMTkwOTUzNDBaFw0zNjAyMTcwOTUzNDBaMDYxHzAd
BgNVBAMMFnRmLXRlc3QtY2EuZXhhbXBsZS5jb20xEzARBgNVBAoMClRGIFRlc3Qg
Q0EwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDjQdPc0vEL19sGV3AI
yEe5GODxQRjRxFnfkWX2KvMYQipA5JZN0WgHFe3EMHNRs+9yULyz5VfX/LDMEaQh
gt0lQ0ZDatLSSKvGLwA9Ben74+MIFIqSycXyo11k5zESz2vYd0z78esVYxo+SBO7
OmkZdpuoKoYnBxsNCk5eN/T58+BxfsDce2dDG64zjW1eWnFAcJ/WQuxybJQTBWQr
9PNqGKilPJLvLe9eUhLCbjDojvNxcfyqEv5w+Ln5Fse//sypljheJRlLyHwGWLZn
3nYMml+S6Z7hpE6kS3YwZfdjm8+2unHMoUmNNxtYHJVFkUX/h8vuU7lZJScIWSgx
4zZbAgMBAAGjUzBRMB0GA1UdDgQWBBR+1ZMyCKYP3vMkVtEmHMSsQIWI4jAfBgNV
HSMEGDAWgBR+1ZMyCKYP3vMkVtEmHMSsQIWI4jAPBgNVHRMBAf8EBTADAQH/MA0G
CSqGSIb3DQEBCwUAA4IBAQBzkOmUVy/8YOHJNLcvluTbhH+0q3DgZfmCmVG+iKk4
dXyzgwmfXtS6Y8LUzT+hub2l+Rlt78qczLoON6V/tM0BnTEuUATSc3DmMx2DRsIZ
FnFM6mwk2WmUs9j3E9EsM9JoFrKl5RQz9H76/q0si5Z+g5lqvKswcLtjGiBZQLt0
fIiMi8TO1oPhIjq7LsK8u1eUQcVetuGOQVRiQ1FK1t/WXCnTNRjWYXydalR2GU2b
iPRh/ZWXfWJLjYbe021r8uuNeYfja+Vlqdan3drTkpeC+57HILU7EpITGbZ8WXVS
OK79sVMVjRGHAEhkkCz2uULKjDpjfb12QW45cKEo/Hce
-----END CERTIFICATE-----
EOT
}

output "cert_id" {
  value = gcore_cdn_trusted_ca_certificate.test.id
}
output "cert_name" {
  value = gcore_cdn_trusted_ca_certificate.test.name
}
