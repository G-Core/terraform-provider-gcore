# Create an IPv4 subnet with custom DNS and host routes
resource "gcore_cloud_network_subnet" "subnet_ipv4" {
  project_id = 1
  region_id  = 1

  name            = "subnet-ipv4"
  cidr            = "192.168.10.0/24"
  network_id      = gcore_cloud_network.network.id
  dns_nameservers = ["8.8.4.4", "1.1.1.1"]
  enable_dhcp     = true
  gateway_ip      = "192.168.10.1"
  ip_version      = 4

  host_routes = [
    {
      destination = "10.0.3.0/24"
      nexthop     = "10.0.0.13"
    },
    {
      destination = "10.0.4.0/24"
      nexthop     = "10.0.0.14"
    },
  ]

  tags = {
    environment = "production"
  }
}
