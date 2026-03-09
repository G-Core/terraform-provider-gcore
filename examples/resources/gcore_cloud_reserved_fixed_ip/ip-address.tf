# Reserve a specific IP address in a subnet
resource "gcore_cloud_reserved_fixed_ip" "specific_ip" {
  project_id = 1
  region_id  = 1

  type       = "ip_address"
  network_id = gcore_cloud_network.private_network.id
  subnet_id  = gcore_cloud_network_subnet.private_subnet_0.id
  ip_address = "172.16.0.254"

  is_vip = false
}
