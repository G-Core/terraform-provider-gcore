terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

resource "gcore_cloud_network" "test" {
  project_id = 379987
  region_id  = 76
  name       = "test-add-route-iface-net"
}

resource "gcore_cloud_network_subnet" "subnet1" {
  project_id = 379987
  region_id  = 76
  name       = "subnet1"
  cidr       = "192.168.1.0/24"
  network_id = gcore_cloud_network.test.id
}

resource "gcore_cloud_network_subnet" "subnet2" {
  project_id = 379987
  region_id  = 76
  name       = "subnet2"
  cidr       = "192.168.2.0/24"
  network_id = gcore_cloud_network.test.id
}

resource "gcore_cloud_network_router" "test" {
  project_id = 379987
  region_id  = 76
  name       = "test-add-route-iface"

  external_gateway_info = {
    enable_snat = true
  }

  # Start with ONLY subnet1
  interfaces = [
    {
      type      = "subnet"
      subnet_id = gcore_cloud_network_subnet.subnet1.id
    }
  ]

  # NO routes initially
}

output "router_id" {
  value = gcore_cloud_network_router.test.id
}

output "subnet2_gateway" {
  value = gcore_cloud_network_subnet.subnet2.gateway_ip
}
