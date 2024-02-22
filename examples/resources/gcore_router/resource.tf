provider gcore {
  permanent_api_token = "251$d3361.............1b35f26d8"
}

resource "gcore_router" "router" {
  name = "router_example"

  dynamic external_gateway_info {
    iterator = egi
    for_each = var.external_gateway_info
    content {
      type        = egi.value.type
      enable_snat = egi.value.enable_snat
      network_id  = egi.value.network_id
    }
  }

  dynamic interfaces {
    iterator = ifaces
    for_each = var.interfaces
    content {
      type      = ifaces.value.type
      subnet_id = ifaces.value.subnet_id
    }
  }

  dynamic routes {
    iterator = rs
    for_each = var.routes
    content {
      destination = rs.value.destination
      nexthop     = rs.value.nexthop
    }
  }

  region_id  = 1
  project_id = 1
}
