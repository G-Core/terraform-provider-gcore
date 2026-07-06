data "gcore_cloud_load_balancer_listeners" "example_cloud_load_balancer_listeners" {
  project_id = 1
  region_id = 1
  load_balancer_id = "00000000-0000-4000-8000-000000000000"
  name = "listener-name"
}
