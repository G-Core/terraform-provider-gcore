resource "gcore_subnet" "subnet_ipv6" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  name            = "subnet_ipv6"
  cidr            = "fd00::/8"
  network_id      = gcore_network.network.id
}
