data "gcore_cloud_floating_ips" "example_cloud_floating_ips" {
  project_id = 1
  region_id = 1
  tag_key = ["key1", "key2"]
  tag_key_value = "tag_key_value"
}
