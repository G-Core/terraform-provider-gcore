resource "gcore_router" "with_routes" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  name = "router-with-static-routes"

  external_gateway_info {
    type        = "default"
    enable_snat = true
  }

  # Static route to reach 10.0.0.0/8 via a specific next hop
  routes {
    destination = "10.0.0.0/8"
    nexthop     = "192.168.1.254"
  }

  # Static route to reach another network segment
  routes {
    destination = "172.16.0.0/16"
    nexthop     = "192.168.1.253"
  }
}
