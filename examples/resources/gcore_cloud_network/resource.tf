resource "gcore_cloud_network" "example_cloud_network" {
  project_id = 1
  region_id = 1
  name = "my network"
  create_router = true
  tags = {
    my-tag = "my-tag-value"
  }
  type = "vxlan"
}
