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
  name       = "mitm-test-network"
}

resource "gcore_cloud_network_subnet" "subnet1" {
  project_id  = 379987
  region_id   = 76
  name        = "mitm-subnet-1"
  network_id  = gcore_cloud_network.test.id
  cidr        = "192.168.10.0/24"
  enable_dhcp = true
}

# Sequential creation - subnet2 waits for subnet1

resource "gcore_cloud_network_subnet" "subnet2" {
  depends_on  = [gcore_cloud_network_subnet.subnet1]
  project_id  = 379987
  region_id   = 76
  name        = "mitm-subnet-2"
  network_id  = gcore_cloud_network.test.id
  cidr        = "192.168.20.0/24"
  enable_dhcp = true
}

resource "gcore_cloud_network_router" "router" {
  project_id = 379987
  region_id  = 76
  name       = "mitm-test-router"

  external_gateway_info = {
    enable_snat = true
    type        = "default"
  }

  interfaces = [
    {
      subnet_id = gcore_cloud_network_subnet.subnet1.id
      type      = "subnet"
    }
  ]
}

output "router_id" {
  value = gcore_cloud_network_router.router.id
}
