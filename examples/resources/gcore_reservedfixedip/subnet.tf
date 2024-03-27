resource "gcore_reservedfixedip" "fixed_ip_subnet" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  type       = "subnet"
  network_id = gcore_network.private_network.id
  subnet_id = gcore_subnet.private_subnet[0].id

  is_vip     = false
}
