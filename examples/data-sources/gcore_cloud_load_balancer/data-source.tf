data "gcore_cloud_load_balancer" "example_cloud_load_balancer" {
  project_id = 1
  region_id = 7
  load_balancer_id = "ac307687-31a4-4a11-a949-6bea1b2878f5"
  show_stats = true
  with_ddos = true
}
