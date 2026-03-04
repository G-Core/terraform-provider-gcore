terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {
  # Credentials from .env file
}

# Get project and region data
data "gcore_cloud_project" "project" {
  id = 1
}

data "gcore_cloud_region" "region" {
  id = 1
}

# Create a router for testing
resource "gcore_cloud_network_router" "test_router" {
  project_id = data.gcore_cloud_project.project.id
  region_id  = data.gcore_cloud_region.region.id
  name       = "terraform-test-router-drift"
  
  external_gateway_info = {
    enable_snat = true
    type        = "default"
  }
}

output "router_id" {
  value = gcore_cloud_network_router.test_router.id
}

output "router_name" {
  value = gcore_cloud_network_router.test_router.name
}

output "router_status" {
  value = gcore_cloud_network_router.test_router.status
}
