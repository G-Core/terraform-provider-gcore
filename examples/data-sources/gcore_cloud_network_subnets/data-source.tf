data "gcore_cloud_network_subnets" "example_cloud_network_subnets" {
  project_id = 1
  region_id = 1
  network_id = "b30d0de7-bca2-4c83-9c57-9e645bd2cc92"
  tag_key = ["key1", "key2"]
  tag_key_value = "tag_key_value"
}
