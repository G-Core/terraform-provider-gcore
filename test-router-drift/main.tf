terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {
  # Credentials from environment variables
}

# Create network for testing
resource "gcore_cloud_network" "test_network" {
  project_id = 379987
  region_id  = 76
  name       = "qa-router-drift-network"
}

# Create subnet 1
resource "gcore_cloud_network_subnet" "subnet1" {
  project_id = 379987
  region_id  = 76
  name       = "qa-router-drift-subnet-1"
  network_id = gcore_cloud_network.test_network.id
  cidr       = "192.168.100.0/24"
  enable_dhcp = true
}

# Create subnet 2
resource "gcore_cloud_network_subnet" "subnet2" {
  project_id = 379987
  region_id  = 76
  name       = "qa-router-drift-subnet-2"
  network_id = gcore_cloud_network.test_network.id
  cidr       = "192.168.200.0/24"
  enable_dhcp = true
}

# Test 1: Basic router with external gateway (default type)
resource "gcore_cloud_network_router" "router_basic" {
  project_id = 379987
  region_id  = 76
  name       = "qa-router-drift-basic"

  external_gateway_info = {
    enable_snat = true
    type        = "default"
  }
}

# Test 2: Router with interfaces
resource "gcore_cloud_network_router" "router_interfaces" {
  project_id = 379987
  region_id  = 76
  name       = "qa-router-drift-interfaces"

  external_gateway_info = {
    enable_snat = true
    type        = "default"
  }

  interfaces = [
    {
      subnet_id = gcore_cloud_network_subnet.subnet1.id
      type      = "subnet"
    },
    {
      subnet_id = gcore_cloud_network_subnet.subnet2.id
      type      = "subnet"
    }
  ]
}

# Test 3: Router with custom routes
resource "gcore_cloud_network_router" "router_routes" {
  project_id = 379987
  region_id  = 76
  name       = "qa-router-drift-routes"

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

  routes = [
    {
      destination = "10.0.0.0/24"
      nexthop     = "192.168.100.10"
    },
    {
      destination = "10.1.0.0/24"
      nexthop     = "192.168.100.20"
    }
  ]
}

# Test 4: Router with everything
resource "gcore_cloud_network_router" "router_complete" {
  project_id = 379987
  region_id  = 76
  name       = "qa-router-drift-complete"

  external_gateway_info = {
    enable_snat = true
    type        = "default"
  }

  interfaces = [
    {
      subnet_id = gcore_cloud_network_subnet.subnet1.id
      type      = "subnet"
    },
    {
      subnet_id = gcore_cloud_network_subnet.subnet2.id
      type      = "subnet"
    }
  ]

  routes = [
    {
      destination = "172.16.0.0/16"
      nexthop     = "192.168.100.1"
    }
  ]
}

# Outputs for verification
output "router_basic_id" {
  value = gcore_cloud_network_router.router_basic.id
}

output "router_basic_external_gateway" {
  value = gcore_cloud_network_router.router_basic.external_gateway_info
}

output "router_interfaces_id" {
  value = gcore_cloud_network_router.router_interfaces.id
}

output "router_interfaces_list" {
  value = gcore_cloud_network_router.router_interfaces.interfaces
}

output "router_routes_id" {
  value = gcore_cloud_network_router.router_routes.id
}

output "router_routes_list" {
  value = gcore_cloud_network_router.router_routes.routes
}

output "router_complete_id" {
  value = gcore_cloud_network_router.router_complete.id
}

output "router_complete_state" {
  value = {
    interfaces = gcore_cloud_network_router.router_complete.interfaces
    routes     = gcore_cloud_network_router.router_complete.routes
    external   = gcore_cloud_network_router.router_complete.external_gateway_info
  }
}
