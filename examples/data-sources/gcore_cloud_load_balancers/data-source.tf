data "gcore_cloud_load_balancers" "example_cloud_load_balancers" {
  project_id = 1
  region_id = 7
  assigned_floating = true
  logging_enabled = true
  name = "lb_name"
  operating_status = "DEGRADED"
  provisioning_status = "ACTIVE"
  tag_key = ["key1", "key2"]
  tag_key_value = "tag_key_value"
}
