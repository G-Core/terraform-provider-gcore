terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

# TC-UPDATE-001: Route removal test (main bug scenario)
# Purpose: Verify route deletion when routes block is removed from config

resource "gcore_cloud_network" "network" {
  project_id = 379987
  region_id  = 76
  name       = "test-router-update-001-network"
  type       = "vxlan"
}

resource "gcore_cloud_network_subnet" "subnet" {
  project_id  = 379987
  region_id   = 76
  name        = "test-router-update-001-subnet"
  cidr        = "192.168.30.0/24"
  network_id  = gcore_cloud_network.network.id
}

variable "include_routes" {
  type    = bool
  default = true
}

resource "gcore_cloud_network_router" "router" {
  project_id = 379987
  region_id  = 76
  name       = "test-router-update-001"

  external_gateway_info = {
    enable_snat = true
  }

  interfaces = [{
    subnet_id = gcore_cloud_network_subnet.subnet.id
    type      = "subnet"
  }]

  # Conditionally include routes
  routes = var.include_routes ? [
    {
      destination = "10.0.3.0/24"
      nexthop     = "192.168.30.1"
    }
  ] : []
}

output "router_id" {
  value = gcore_cloud_network_router.router.id
}

output "router_routes" {
  value = gcore_cloud_network_router.router.routes
}
