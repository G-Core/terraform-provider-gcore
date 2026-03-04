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
  name       = "qa-terr-router-network"
}

resource "gcore_cloud_network_subnet" "sb" {
  project_id  = 379987
  region_id   = 76
  name        = "qa-terr-router-subnet"
  cidr        = "192.168.0.0/24"
  network_id  = gcore_cloud_network.network.id
}

resource "gcore_cloud_network_router" "router" {
  project_id = 379987
  region_id  = 76
  name       = "qa-terr-route-bug"
  external_gateway_info = {
    enable_snat = true
  }
  interfaces = [{
    subnet_id = gcore_cloud_network_subnet.sb.id
    type      = "subnet"
  }]
  # routes removed - should delete the route
}
