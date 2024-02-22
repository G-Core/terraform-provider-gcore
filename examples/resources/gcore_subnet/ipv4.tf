resource "gcore_subnet" "subnet_ipv4" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  name            = "subnet_ipv4"
  cidr            = "192.168.10.0/24"
  network_id      = gcore_network.network.id
  dns_nameservers = var.dns_nameservers

  dynamic host_routes {
    iterator = hr
    for_each = var.host_routes
    content {
      destination = hr.value.destination
      nexthop     = hr.value.nexthop
    }
  }

  gateway_ip = "192.168.10.1"
}