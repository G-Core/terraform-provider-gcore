resource "gcore_cloud_network" "private_network_dualstack" {
  project_id = 1
  region_id  = 1

  name = "my-private-network-dualstack"
}

resource "gcore_cloud_network_subnet" "private_subnet_ipv4" {
  project_id = 1
  region_id  = 1

  cidr       = "10.0.0.0/24"
  name       = "my-private-network-subnet-ipv4"
  network_id = gcore_cloud_network.private_network_dualstack.id
}

resource "gcore_cloud_network_subnet" "private_subnet_ipv6" {
  project_id = 1
  region_id  = 1

  cidr       = "fd00::/120"
  name       = "my-private-network-subnet-ipv6"
  network_id = gcore_cloud_network.private_network_dualstack.id
}

# A dual-stack Reserved Fixed IP takes one IPv4 and one IPv6 address from the
# network's subnets; the load balancer uses its port as the VIP.
resource "gcore_cloud_reserved_fixed_ip" "private_lb_dualstack_vip" {
  project_id = 1
  region_id  = 1

  type       = "any_subnet"
  network_id = gcore_cloud_network.private_network_dualstack.id
  ip_family  = "dual"

  is_vip = false

  depends_on = [
    gcore_cloud_network_subnet.private_subnet_ipv4,
    gcore_cloud_network_subnet.private_subnet_ipv6,
  ]
}

resource "gcore_cloud_load_balancer" "private_lb_dualstack" {
  project_id = 1
  region_id  = 1

  name        = "My first private dual stack load balancer"
  flavor      = "lb1-1-2"
  vip_port_id = gcore_cloud_reserved_fixed_ip.private_lb_dualstack_vip.port_id
}
