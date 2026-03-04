terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

# Dependencies
resource "gcore_cloud_network" "test" {
  project_id = var.project_id
  region_id  = var.region_id
  name       = "test-router-net"
}

resource "gcore_cloud_network_subnet" "subnet1" {
  project_id = var.project_id
  region_id  = var.region_id
  name       = "subnet1"
  cidr       = "192.168.1.0/24"
  network_id = gcore_cloud_network.test.id
}

resource "gcore_cloud_network_subnet" "subnet2" {
  project_id = var.project_id
  region_id  = var.region_id
  name       = "subnet2"
  cidr       = "192.168.2.0/24"
  network_id = gcore_cloud_network.test.id
}

# Router under test - using attribute syntax
resource "gcore_cloud_network_router" "test" {
  project_id = var.project_id
  region_id  = var.region_id
  name       = var.router_name

  interfaces = [
    for subnet_id in var.interfaces : {
      type      = "subnet"
      subnet_id = subnet_id
    }
  ]

  routes = var.routes

  external_gateway_info = var.enable_external_gateway ? {
    enable_snat = var.external_gateway_snat
  } : null
}

output "router_id" {
  value = gcore_cloud_network_router.test.id
}

output "subnet1_id" {
  value = gcore_cloud_network_subnet.subnet1.id
}

output "subnet2_id" {
  value = gcore_cloud_network_subnet.subnet2.id
}
