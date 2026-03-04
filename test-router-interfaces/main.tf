terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {
  # Credentials from environment
}

# Create a network
resource "gcore_cloud_network" "test_network" {
  project_id = 379987
  region_id  = 76
  name       = "test-router-network"
}

# Create a subnet
resource "gcore_cloud_network_subnet" "test_subnet" {
  project_id = 379987
  region_id  = 76
  name       = "test-router-subnet"
  network_id = gcore_cloud_network.test_network.id
  cidr       = "192.168.1.0/24"

  # Enable gateway for router attachment
  enable_dhcp = true
}

# Create a second subnet for testing interface add
resource "gcore_cloud_network_subnet" "test_subnet_2" {
  project_id = 379987
  region_id  = 76
  name       = "test-router-subnet-2"
  network_id = gcore_cloud_network.test_network.id
  cidr       = "192.168.2.0/24"

  enable_dhcp = true
}

# Create a router WITHOUT interfaces initially
resource "gcore_cloud_network_router" "test_router" {
  project_id = 379987
  region_id  = 76
  name       = "test-router-interface-fix"

  external_gateway_info = {
    enable_snat = true
    type        = "default"
  }

  # Start with one interface
  interfaces = [
    {
      subnet_id = gcore_cloud_network_subnet.test_subnet.id
      type      = "subnet"
    }
  ]
}

output "router_id" {
  value = gcore_cloud_network_router.test_router.id
}

output "router_interfaces" {
  value = gcore_cloud_network_router.test_router.interfaces
}

output "subnet_1_id" {
  value = gcore_cloud_network_subnet.test_subnet.id
}

output "subnet_2_id" {
  value = gcore_cloud_network_subnet.test_subnet_2.id
}
