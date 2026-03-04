terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

resource "gcore_cloud_network" "net" {
  project_id = 379987
  region_id  = 76
  name       = "test-router-fix-network"
}

resource "gcore_cloud_subnet" "subnet1" {
  project_id  = 379987
  region_id   = 76
  name        = "test-subnet-1"
  network_id  = gcore_cloud_network.net.id
  cidr        = "192.168.1.0/24"
  enable_dhcp = true
}

resource "gcore_cloud_router" "router" {
  project_id = 379987
  region_id  = 76
  name       = "test-router-fix"

  external_gateway_info = {
    enable_snat = true
    type        = "default"
  }
}

output "router_id" {
  value = gcore_cloud_router.router.id
}
