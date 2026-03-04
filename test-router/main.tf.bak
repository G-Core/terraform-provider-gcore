terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {
  # Credentials from environment variables:
  # GCORE_API_KEY, GCORE_CLOUD_PROJECT_ID, GCORE_CLOUD_REGION_ID
}

# Create a network for testing
resource "gcore_cloud_network" "test_network" {
  project_id = 379987
  region_id  = 76
  name       = "qa-terr-router-test-network"
}

# Create subnet 1 for router interface
resource "gcore_cloud_network_subnet" "subnet1" {
  project_id = 379987
  region_id  = 76
  name       = "qa-terr-router-subnet-1"
  network_id = gcore_cloud_network.test_network.id
  cidr       = "192.168.10.0/24"

  enable_dhcp = true
}

# Create subnet 2 for testing interface add/remove
resource "gcore_cloud_network_subnet" "subnet2" {
  project_id = 379987
  region_id  = 76
  name       = "qa-terr-router-subnet-2"
  network_id = gcore_cloud_network.test_network.id
  cidr       = "192.168.20.0/24"

  enable_dhcp = true
}

# Create router with ONE interface initially
resource "gcore_cloud_network_router" "router" {
  project_id = 379987
  region_id  = 76
  name       = "qa-terr-router"

  external_gateway_info = {
    enable_snat = true
    type        = "default"
  }

  # Add second interface
  interfaces = [
    {
      subnet_id = gcore_cloud_network_subnet.subnet1.id
      type      = "subnet"
    }
  ]
}

# Outputs for verification
output "router_id" {
  value = gcore_cloud_network_router.router.id
}

output "router_name" {
  value = gcore_cloud_network_router.router.name
}

output "router_status" {
  value = gcore_cloud_network_router.router.status
}

output "router_interfaces" {
  value = gcore_cloud_network_router.router.interfaces
}

output "subnet1_id" {
  value = gcore_cloud_network_subnet.subnet1.id
}

output "subnet2_id" {
  value = gcore_cloud_network_subnet.subnet2.id
}

output "network_id" {
  value = gcore_cloud_network.test_network.id
}
