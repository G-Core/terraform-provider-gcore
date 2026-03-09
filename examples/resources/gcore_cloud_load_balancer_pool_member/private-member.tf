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

resource "gcore_cloud_reserved_fixed_ip" "fixed_ip" {
  project_id = 1
  region_id  = 1

  type             = "ip_address"
  network_id       = gcore_cloud_network.private_network.id
  subnet_id        = gcore_cloud_network_subnet.private_subnet.id
  fixed_ip_address = "10.0.0.10"
  is_vip           = false
}

resource "gcore_cloud_load_balancer_pool_member" "private_member" {
  project_id = 1
  region_id  = 1

  pool_id = gcore_cloud_load_balancer_pool.http.id

  address   = gcore_cloud_reserved_fixed_ip.fixed_ip.fixed_ip_address
  subnet_id = gcore_cloud_reserved_fixed_ip.fixed_ip.subnet_id

  protocol_port = 80
  weight        = 1
}
