# Create a container registry credential for private image pulls
resource "gcore_cloud_inference_registry_credential" "example" {
  project_id         = 1
  name               = "docker-io"
  username           = "my-username"
  password_wo        = "my-password"
  password_wo_version = 1
  registry_url       = "docker.io"
}
