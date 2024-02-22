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

resource "gcore_reservedfixedip" "fixed_ip" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  type             = "ip_address"
  network_id       = gcore_network.private_network.id
  subnet_id        = gcore_subnet.private_subnet.id
  fixed_ip_address = "10.0.0.10"
  is_vip           = false
}

resource "gcore_lbmember" "private_member" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  pool_id = gcore_lbpool.http.id

  address   = gcore_reservedfixedip.fixed_ip.fixed_ip_address
  subnet_id = gcore_reservedfixedip.fixed_ip.subnet_id

  protocol_port = 80
  weight        = 1
}
