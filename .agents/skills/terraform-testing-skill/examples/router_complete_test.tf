# Router Testing Configuration - Complete Example
# This tests interfaces and routes with corner cases

terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
      version = "~> 1.0"
    }
  }
}

provider "gcore" {
  # Credentials from environment variables
}

# Variables for dynamic testing
variable "test_name_suffix" {
  default = "test"
}

variable "enable_routes" {
  default = true
  type    = bool
}

# Create networks for testing
resource "gcore_cloud_network" "test_net1" {
  name = "test-network-1-${var.test_name_suffix}"
}

resource "gcore_cloud_subnet" "test_subnet1" {
  name       = "test-subnet-1-${var.test_name_suffix}"
  network_id = gcore_cloud_network.test_net1.id
  cidr       = "10.1.0.0/24"
}

resource "gcore_cloud_network" "test_net2" {
  name = "test-network-2-${var.test_name_suffix}"
}

resource "gcore_cloud_subnet" "test_subnet2" {
  name       = "test-subnet-2-${var.test_name_suffix}"
  network_id = gcore_cloud_network.test_net2.id
  cidr       = "10.2.0.0/24"
}

# Create instance for interface attachment testing
resource "gcore_cloud_instance" "test" {
  name   = "test-instance-${var.test_name_suffix}"
  flavor = "g1-standard-1-2"
  
  boot_volume {
    size   = 10
    source = "image"
    image_id = data.gcore_cloud_image.ubuntu.id
  }
  
  network_interface {
    network_id = gcore_cloud_network.test_net1.id
    subnet_id  = gcore_cloud_subnet.test_subnet1.id
  }
}

data "gcore_cloud_image" "ubuntu" {
  name = "ubuntu-22.04"
  os_distro = "ubuntu"
}

# Main router resource to test
resource "gcore_cloud_router" "test" {
  name = "test-router-${var.test_name_suffix}"
  
  # Test interfaces - attach/detach operations
  interfaces = [
    {
      network_id = gcore_cloud_network.test_net1.id
      subnet_id  = gcore_cloud_subnet.test_subnet1.id
    },
    {
      network_id = gcore_cloud_network.test_net2.id
      subnet_id  = gcore_cloud_subnet.test_subnet2.id
    }
  ]
  
  # Test routes - should trigger PATCH with routes=[] when cleared
  routes = var.enable_routes ? [
    {
      destination = "192.168.1.0/24"
      nexthop     = "10.1.0.1"
    },
    {
      destination = "192.168.2.0/24"
      nexthop     = "10.2.0.1"
    },
    {
      destination = "172.16.0.0/16"
      nexthop     = "10.1.0.254"
    }
  ] : []
  
  # Test that external gateway info is handled correctly
  external_gateway_info {
    network_id = data.gcore_cloud_network.external.id
    enable_snat = true
  }
}

# Data source for external network
data "gcore_cloud_network" "external" {
  name = "external-network"
  external = true
}

# Outputs for verification
output "router_id" {
  value = gcore_cloud_router.test.id
  description = "Router ID - should remain stable during updates"
}

output "router_interfaces" {
  value = gcore_cloud_router.test.interfaces
  description = "Attached interfaces"
}

output "router_routes" {
  value = gcore_cloud_router.test.routes
  description = "Configured routes"
}

# Test scenarios to run:
# 1. Initial apply - verify all resources created
# 2. terraform plan - should show no changes (drift test)
# 3. terraform apply -var="enable_routes=false" - should send routes=[] 
# 4. terraform apply -var="test_name_suffix=updated" - should PATCH name
# 5. Remove interface from list - should call detach_interface
# 6. Add interface to list - should call attach_interface