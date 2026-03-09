# Reserve a private IP in any available subnet of the network
resource "gcore_cloud_reserved_fixed_ip" "in_any_subnet" {
  project_id = 1
  region_id  = 1

  type       = "any_subnet"
  network_id = gcore_cloud_network.private_network.id

  is_vip = false
}
