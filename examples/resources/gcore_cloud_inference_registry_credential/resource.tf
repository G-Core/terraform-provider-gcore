resource "gcore_cloud_inference_registry_credential" "example_cloud_inference_registry_credential" {
  project_id = 1
  name = "docker-io"
  password = "password"
  registry_url = "registry.example.com"
  username = "username"
}
