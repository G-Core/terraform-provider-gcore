data "gcore_cloud_load_balancer" "example_cloud_load_balancer" {
  project_id = 0
  region_id = 0
  load_balancer_id = "load_balancer_id"
  show_stats = true
  with_ddos = true
}
