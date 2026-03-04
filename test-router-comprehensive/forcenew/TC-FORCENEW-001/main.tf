terraform {
  required_providers {
    gcore = { source = "gcore/gcore" }
  }
}
provider "gcore" {}

variable "project_id" {
  type = number
  default = 379987
}

resource "gcore_cloud_network" "network" {
  project_id = var.project_id
  region_id  = 76
  name       = "test-router-forcenew-network"
  type       = "vxlan"
}

resource "gcore_cloud_network_subnet" "subnet" {
  project_id  = var.project_id
  region_id   = 76
  name        = "test-router-forcenew-subnet"
  cidr        = "192.168.60.0/24"
  network_id  = gcore_cloud_network.network.id
}

resource "gcore_cloud_network_router" "router" {
  project_id = var.project_id
  region_id  = 76
  name       = "test-router-forcenew"
  external_gateway_info = { enable_snat = true }
  interfaces = [{ subnet_id = gcore_cloud_network_subnet.subnet.id, type = "subnet" }]
}

output "router_id" { value = gcore_cloud_network_router.router.id }
