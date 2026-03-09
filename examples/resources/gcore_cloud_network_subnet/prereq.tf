# Create a network to hold the subnets
resource "gcore_cloud_network" "network" {
  project_id = 1
  region_id  = 1

  name = "network-example"
  type = "vxlan"
}
