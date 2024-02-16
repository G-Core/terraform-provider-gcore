resource "gcore_network" "private_network_dualstack" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  name = "my-private-network-dualstack"
}

resource "gcore_subnet" "private_subnet_ipv4" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  cidr       = "10.0.0.0/24"
  name       = "my-private-network-subnet-ipv4"
  network_id = gcore_network.private_network_dualstack.id
}

resource "gcore_subnet" "private_subnet_ipv6" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  cidr       = "fd00::/120"
  name       = "my-private-network-subnet-ipv6"
  network_id = gcore_network.private_network_dualstack.id
}

resource "gcore_loadbalancerv2" "private_lb_dualstack" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  name       = "My first private dual stack load balancer"
  flavor     = "lb1-1-2"
  vip_network_id = gcore_network.private_network_dualstack.id
  vip_ip_family = "dual"
  depends_on = [
    gcore_subnet.private_subnet_ipv4,
    gcore_subnet.private_subnet_ipv6,
  ]
}
