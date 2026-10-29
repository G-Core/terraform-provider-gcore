resource "gcore_cloud_security_group" "example_cloud_security_group" {
  project_id = 1
  region_id = 1
  name = "my_security_group"
  description = "My security group description"
  tags = {
    my-tag = "my-tag-value"
  }
}
