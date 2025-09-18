data "gcore_cloud_networks" "example_cloud_networks" {
  project_id = 1
  region_id = 1
  name = "my-network"
  tag_key = ["key1", "key2"]
  tag_key_value = "tag_key_value"
}
