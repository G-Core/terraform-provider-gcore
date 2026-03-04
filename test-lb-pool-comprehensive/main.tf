terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

variable "api_key" {
  type        = string
  description = "Gcore API key"
  sensitive   = true
}

provider "gcore" {
  api_key = var.api_key
}

locals {
  project_id = 379987
  region_id  = 76
}

# Load Balancer (dependency)
resource "gcore_cloud_load_balancer" "test" {
  project_id = local.project_id
  region_id  = local.region_id
  name       = "tf-test-lb-pool-comprehensive"
  flavor     = "lb1-1-2"
}

# Listener (dependency)
resource "gcore_cloud_load_balancer_listener" "test" {
  project_id       = local.project_id
  region_id        = local.region_id
  load_balancer_id = gcore_cloud_load_balancer.test.id
  name             = "tf-test-listener-comprehensive"
  protocol         = "HTTP"
  protocol_port    = 80
}

# Main Pool under test
resource "gcore_cloud_load_balancer_pool" "test" {
  project_id       = local.project_id
  region_id        = local.region_id
  load_balancer_id = gcore_cloud_load_balancer.test.id
  listener_id      = gcore_cloud_load_balancer_listener.test.id
  name             = "tf-test-pool-hm-fix"
  protocol         = "HTTP"
  lb_algorithm     = "ROUND_ROBIN"

  # Test 5: Add healthmonitor
  members = [
    {
      address       = "10.0.0.1"
      protocol_port = 8080
      weight        = 1
    },
    {
      address       = "10.0.0.2"
      protocol_port = 8080
      weight        = 1
    }
  ]

  # Case 1: Delete entire attribute (omit from config)
}

# Outputs for verification
output "pool_id" {
  value = gcore_cloud_load_balancer_pool.test.id
}

output "pool_members" {
  value = gcore_cloud_load_balancer_pool.test.members
}

output "pool_healthmonitor" {
  value = gcore_cloud_load_balancer_pool.test.healthmonitor
}

output "pool_operating_status" {
  value = gcore_cloud_load_balancer_pool.test.operating_status
}
