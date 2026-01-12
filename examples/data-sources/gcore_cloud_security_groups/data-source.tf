data "gcore_cloud_security_groups" "example_cloud_security_groups" {
  project_id = 1
  region_id = 1
  name = "my_security_group"
  tag_key = ["key1", "key2"]
  tag_key_value = "tag_key_value"
}
