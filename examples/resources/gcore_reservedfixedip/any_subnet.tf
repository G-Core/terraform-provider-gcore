resource "gcore_reservedfixedip" "fixed_ip_in_any_subnet" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  type       = "any_subnet"
  network_id = gcore_network.private_network.id

  is_vip     = false
}
