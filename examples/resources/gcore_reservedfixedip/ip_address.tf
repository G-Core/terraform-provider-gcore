locals {
  selected_subnet = gcore_subnet.private_subnet[0]
}

resource "gcore_reservedfixedip" "fixed_ip_ip_address" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  type       = "ip_address"
  network_id = gcore_network.private_network.id
  subnet_id = local.selected_subnet.id

  fixed_ip_address = cidrhost(local.selected_subnet.cidr, 254)

  is_vip     = false
}
