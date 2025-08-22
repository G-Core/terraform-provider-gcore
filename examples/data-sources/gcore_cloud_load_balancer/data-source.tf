data "gcore_cloud_load_balancer" "example_cloud_load_balancer" {
  project_id = 0
  region_id = 0
  loadbalancer_id = "loadbalancer_id"
  show_stats = true
  with_ddos = true
}
