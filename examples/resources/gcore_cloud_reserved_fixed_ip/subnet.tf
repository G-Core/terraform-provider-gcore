# Reserve a private IP in a specific subnet
resource "gcore_cloud_reserved_fixed_ip" "in_subnet" {
  project_id = 1
  region_id  = 1

  type       = "subnet"
  network_id = gcore_cloud_network.private_network.id
  subnet_id  = gcore_cloud_network_subnet.private_subnet_0.id

  is_vip = false
}
