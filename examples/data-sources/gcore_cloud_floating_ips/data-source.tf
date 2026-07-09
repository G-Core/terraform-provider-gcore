data "gcore_cloud_floating_ips" "example_cloud_floating_ips" {
  project_id = 1
  region_id = 1
  port_ids = ["ee2402d0-f0cd-4503-9b75-69be1d11c5f1"]
  status = "ACTIVE"
  tag_key = ["key1", "key2"]
  tag_key_value = "tag_key_value"
}
