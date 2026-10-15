# Create a TLS secret for use with load balancers
resource "gcore_cloud_secret" "example" {
  project_id = 1
  region_id  = 1

  name = "my-tls-certificate"
  payload = {
    certificate_wo       = file("${path.module}/cert.pem")
    certificate_chain_wo = file("${path.module}/chain.pem")
    private_key_wo       = file("${path.module}/key.pem")
  }
  payload_wo_version = 1
  expiration         = "2030-12-31T23:59:59Z"
}
