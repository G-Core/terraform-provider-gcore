terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

# TC-DRIFT-002: Router with routes specified
# Purpose: Verify no drift when routes are explicitly configured

resource "gcore_cloud_network" "network" {
  project_id = 379987
  region_id  = 76
  name       = "test-router-drift-002-network"
  type       = "vxlan"
}

resource "gcore_cloud_network_subnet" "subnet" {
  project_id  = 379987
  region_id   = 76
  name        = "test-router-drift-002-subnet"
  cidr        = "192.168.20.0/24"
  network_id  = gcore_cloud_network.network.id
}

resource "gcore_cloud_network_router" "router" {
  project_id = 379987
  region_id  = 76
  name       = "test-router-drift-002"

  external_gateway_info = {
    enable_snat = true
  }

  interfaces = [{
    subnet_id = gcore_cloud_network_subnet.subnet.id
    type      = "subnet"
  }]

  # Routes explicitly specified
  routes = [
    {
      destination = "10.0.1.0/24"
      nexthop     = "192.168.20.1"
    },
    {
      destination = "10.0.2.0/24"
      nexthop     = "192.168.20.1"
    }
  ]
}

output "router_id" {
  value = gcore_cloud_network_router.router.id
}

output "router_routes" {
  value = gcore_cloud_network_router.router.routes
}
