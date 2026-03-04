terraform {
  required_providers {
    gcore = { source = "gcore/gcore" }
  }
}
provider "gcore" {}

resource "gcore_cloud_network" "network" {
  project_id = 379987
  region_id  = 76
  name       = "test-router-edge-network"
  type       = "vxlan"
}

resource "gcore_cloud_network_subnet" "subnet" {
  project_id  = 379987
  region_id   = 76
  name        = "test-router-edge-subnet"
  cidr        = "192.168.80.0/24"
  network_id  = gcore_cloud_network.network.id
}

variable "routes_config" {
  type = string
  default = "empty"
}

resource "gcore_cloud_network_router" "router" {
  project_id = 379987
  region_id  = 76
  name       = "test-router-edge"
  external_gateway_info = { enable_snat = true }
  interfaces = [{ subnet_id = gcore_cloud_network_subnet.subnet.id, type = "subnet" }]
  routes = var.routes_config == "empty" ? [] : (
    var.routes_config == "single" ? [{ destination = "10.0.8.0/24", nexthop = "192.168.80.1" }] : [
      { destination = "10.0.8.0/24", nexthop = "192.168.80.1" },
      { destination = "10.0.9.0/24", nexthop = "192.168.80.1" },
      { destination = "10.0.10.0/24", nexthop = "192.168.80.1" }
    ]
  )
}
