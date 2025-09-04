data "gcore_cloud_load_balancers" "example_cloud_load_balancers" {
  project_id = 0
  region_id = 0
  assigned_floating = true
  logging_enabled = true
  name = "name"
  order_by = "order_by"
  show_stats = true
  tag_key = ["string"]
  tag_key_value = "tag_key_value"
  with_ddos = true
}
