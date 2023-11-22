provider gcore {
  permanent_api_token = "251$d3361.............1b35f26d8"
}

resource "gcore_network" "network" {
  name       = "network_example"
  type       = "vxlan"
  region_id  = 1
  project_id = 1
}

resource "gcore_subnet" "subnet" {
  name            = "subnet_example"
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
  region_id  = 1
  project_id = 1
}