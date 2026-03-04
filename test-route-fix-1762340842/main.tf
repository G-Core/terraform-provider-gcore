terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

resource "gcore_cloud_network" "network" {
  project_id = 379987
  region_id  = 76
  name       = "test-route-fix-network"
  type       = "vxlan"
}

resource "gcore_cloud_network_subnet" "subnet" {
  project_id  = 379987
  region_id   = 76
  name        = "test-route-fix-subnet"
  cidr        = "192.168.50.0/24"
  network_id  = gcore_cloud_network.network.id
}

variable "include_routes" {
  type    = bool
  default = true
}

resource "gcore_cloud_network_router" "router" {
  project_id = 379987
  region_id  = 76
  name       = "test-route-fix-router"

  external_gateway_info = {
    enable_snat = true
  }

  interfaces = [{
    subnet_id = gcore_cloud_network_subnet.subnet.id
    type      = "subnet"
  }]

  routes = var.include_routes ? [
    {
      destination = "10.0.5.0/24"
      nexthop     = "192.168.50.1"
    }
  ] : []
}

output "router_id" {
  value = gcore_cloud_network_router.router.id
}

output "router_routes" {
  value = gcore_cloud_network_router.router.routes
}
