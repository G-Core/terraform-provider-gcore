# Create a Windows baremetal server with a public interface
resource "gcore_cloud_baremetal_server" "windows_server" {
  project_id          = 1
  region_id           = 1
  flavor              = "bm1-infrastructure-small"
  name                = "my-windows-bare-metal"
  image_id            = "408a0e4d-6a28-4bae-93fa-f738d964f555"
  password_wo         = "my-s3cR3tP@ssw0rd"
  password_wo_version = 1

  interfaces = [{
    type      = "external"
    ip_family = "ipv4"
  }]
}
