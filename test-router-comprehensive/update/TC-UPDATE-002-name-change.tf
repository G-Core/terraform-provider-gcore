terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

# TC-UPDATE-002: Name change test
# Purpose: Verify name update uses PATCH, not replacement

resource "gcore_cloud_network" "network" {
  project_id = 379987
  region_id  = 76
  name       = "test-router-update-002-network"
  type       = "vxlan"
}

resource "gcore_cloud_network_subnet" "subnet" {
  project_id  = 379987
  region_id   = 76
  name        = "test-router-update-002-subnet"
  cidr        = "192.168.40.0/24"
  network_id  = gcore_cloud_network.network.id
}

variable "router_name" {
  type    = string
  default = "test-router-update-002"
}

resource "gcore_cloud_network_router" "router" {
  project_id = 379987
  region_id  = 76
  name       = var.router_name

  external_gateway_info = {
    enable_snat = true
  }

  interfaces = [{
    subnet_id = gcore_cloud_network_subnet.subnet.id
    type      = "subnet"
  }]

  routes = [{
    destination = "10.0.4.0/24"
    nexthop     = "192.168.40.1"
  }]
}

output "router_id" {
  value = gcore_cloud_network_router.router.id
}

output "router_name" {
  value = gcore_cloud_network_router.router.name
}
