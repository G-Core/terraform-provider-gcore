data "gcore_cloud_load_balancer_listener" "example_cloud_load_balancer_listener" {
  project_id = 1
  region_id = 1
  listener_id = "00000000-0000-4000-8000-000000000000"
  show_stats = true
}
