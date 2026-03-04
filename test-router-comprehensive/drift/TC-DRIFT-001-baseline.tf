terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

# TC-DRIFT-001: Baseline router with interface, no routes
# Purpose: Verify no drift with minimal configuration

resource "gcore_cloud_network" "network" {
  project_id = 379987
  region_id  = 76
  name       = "test-router-drift-001-network"
  type       = "vxlan"
}

resource "gcore_cloud_network_subnet" "subnet" {
  project_id  = 379987
  region_id   = 76
  name        = "test-router-drift-001-subnet"
  cidr        = "192.168.10.0/24"
  network_id  = gcore_cloud_network.network.id
}

resource "gcore_cloud_network_router" "router" {
  project_id = 379987
  region_id  = 76
  name       = "test-router-drift-001"

  external_gateway_info = {
    enable_snat = true
  }

  interfaces = [{
    subnet_id = gcore_cloud_network_subnet.subnet.id
    type      = "subnet"
  }]

  # No routes specified - should remain empty with no drift
}

output "router_id" {
  value = gcore_cloud_network_router.router.id
}

output "router_routes" {
  value = gcore_cloud_network_router.router.routes
}
