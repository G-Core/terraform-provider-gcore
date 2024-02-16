resource "gcore_network" "private_network" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  name = "my-private-network"
}

resource "gcore_subnet" "private_subnet" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  cidr       = "10.0.0.0/24"
  name       = "my-private-network-subnet"
  network_id = gcore_network.private_network.id
}

resource "gcore_loadbalancerv2" "private_lb" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  name       = "My first private load balancer"
  flavor     = "lb1-1-2"
  vip_network_id = gcore_network.private_network.id
  vip_subnet_id = gcore_subnet.private_subnet.id
}

output "private_lb_ip" {
  value = gcore_loadbalancerv2.private_lb.vip_address
}
