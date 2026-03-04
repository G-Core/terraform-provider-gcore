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

# Test Load Balancer
resource "gcore_cloud_load_balancer" "test" {
  project_id = local.project_id
  region_id  = local.region_id
  name       = "tf-test-member-removal"
  flavor     = "lb1-1-2"
}

# Test Listener
resource "gcore_cloud_load_balancer_listener" "test" {
  project_id       = local.project_id
  region_id        = local.region_id
  load_balancer_id = gcore_cloud_load_balancer.test.id
  name             = "tf-test-listener-member"
  protocol         = "HTTP"
  protocol_port    = 80
}

# Test Pool with 2 members - we'll remove one to test the bug
resource "gcore_cloud_load_balancer_pool" "test" {
  project_id       = local.project_id
  region_id        = local.region_id
  load_balancer_id = gcore_cloud_load_balancer.test.id
  listener_id      = gcore_cloud_load_balancer_listener.test.id
  name             = "tf-test-pool-members"
  protocol         = "HTTP"
  lb_algorithm     = "ROUND_ROBIN"

  # Removing member 10.0.0.2
  members = [
    {
      address       = "10.0.0.1"
      protocol_port = 8080
      weight        = 1
    }
  ]
}

output "pool_id" {
  value = gcore_cloud_load_balancer_pool.test.id
}

output "pool_members" {
  value = gcore_cloud_load_balancer_pool.test.members
}
