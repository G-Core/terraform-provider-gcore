terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

data "gcore_cloud_region" "rg" {
  name = "Luxembourg-2"
}

resource "gcore_cloud_network" "net" {
  project_id = 379987
  region_id  = data.gcore_cloud_region.rg.id
  name       = "test-router-comprehensive-net"
}

resource "gcore_cloud_network_subnet" "subnet1" {
  project_id  = 379987
  region_id   = data.gcore_cloud_region.rg.id
  name        = "test-subnet-1"
  network_id  = gcore_cloud_network.net.id
  cidr        = "192.168.1.0/24"
  enable_dhcp = true
}

resource "gcore_cloud_network_subnet" "subnet2" {
  project_id  = 379987
  region_id   = data.gcore_cloud_region.rg.id
  name        = "test-subnet-2"
  network_id  = gcore_cloud_network.net.id
  cidr        = "192.168.2.0/24"
  enable_dhcp = true
}

resource "gcore_cloud_network_router" "router" {
  project_id = 379987
  region_id  = data.gcore_cloud_region.rg.id
  name       = "test-router-comprehensive"

  external_gateway_info = {
    enable_snat = true
    type        = "default"
  }
}

output "router_id" {
  value = gcore_cloud_network_router.router.id
}
