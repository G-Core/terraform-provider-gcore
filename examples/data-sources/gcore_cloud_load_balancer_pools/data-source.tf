data "gcore_cloud_load_balancer_pools" "example_cloud_load_balancer_pools" {
  project_id = 1
  region_id = 1
  listener_id = "00000000-0000-4000-8000-000000000000"
  load_balancer_id = "00000000-0000-4000-8000-000000000000"
  name = "lb-pool-name"
}
