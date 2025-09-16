data "gcore_cloud_lblistener" "example_cloud_lblistener" {
  project_id = 1
  region_id = 1
  listener_id = "00000000-0000-4000-8000-000000000000"
  show_stats = true
}
