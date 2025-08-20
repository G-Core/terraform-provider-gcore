resource "gcore_cloud_registry_user" "example_cloud_registry_user" {
  project_id = 0
  region_id = 0
  registry_id = 0
  duration = 14
  name = "user1"
  read_only = false
  secret = "secret"
}
