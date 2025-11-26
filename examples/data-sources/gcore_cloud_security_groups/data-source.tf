data "gcore_cloud_security_groups" "example_cloud_security_groups" {
  project_id = 1
  region_id = 1
  tag_key = ["my-tag"]
  tag_key_value = "tag_key_value"
}
