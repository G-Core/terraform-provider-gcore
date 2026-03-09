# Create an IPv6 subnet
resource "gcore_cloud_network_subnet" "subnet_ipv6" {
  project_id = 1
  region_id  = 1

  name       = "subnet-ipv6"
  cidr       = "fd00::/8"
  network_id = gcore_cloud_network.network.id
  ip_version = 6
}
