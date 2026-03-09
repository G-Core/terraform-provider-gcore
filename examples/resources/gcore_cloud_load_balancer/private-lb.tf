resource "gcore_cloud_network" "private_network" {
  project_id = 1
  region_id  = 1

  name = "my-private-network"
}

resource "gcore_cloud_network_subnet" "private_subnet" {
  project_id = 1
  region_id  = 1

  cidr       = "10.0.0.0/24"
  name       = "my-private-network-subnet"
  network_id = gcore_cloud_network.private_network.id
}

resource "gcore_cloud_load_balancer" "private_lb" {
  project_id = 1
  region_id  = 1

  name           = "My first private load balancer"
  flavor         = "lb1-1-2"
  vip_network_id = gcore_cloud_network.private_network.id
  vip_subnet_id  = gcore_cloud_network_subnet.private_subnet.id
}

output "private_lb_ip" {
  value = gcore_cloud_load_balancer.private_lb.vip_address
}
