# Create a TLS secret for use with load balancers
resource "gcore_cloud_secret" "example" {
  project_id = 1
  region_id  = 1

  name = "my-tls-certificate"
  payload = {
    certificate       = file("${path.module}/cert.pem")
    certificate_chain = file("${path.module}/chain.pem")
    private_key       = file("${path.module}/key.pem")
  }
  expiration = "2025-12-28T19:14:44.000Z"
}
