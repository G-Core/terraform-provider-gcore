data "gcore_cloud_registry_artifacts" "example_cloud_registry_artifacts" {
  project_id = 1
  region_id = 1
  registry_id = 1
  repository_name = "nginx"
}
