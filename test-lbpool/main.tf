terraform {
  required_providers {
    gcore = {
      source  = "stainless-sdks/gcore"
      version = "~> 1.0"
    }
  }
}

provider "gcore" {
  # Uses environment variables from .env file
}

# First create a simple load balancer
resource "gcore_cloud_load_balancer" "test" {
  name       = "test-lb-simple-${formatdate("YYYY-MM-DD-hhmm", timestamp())}"
  flavor     = "lb1-1-2"
  project_id = 379987
  region_id  = 76
}

# Create a pool associated with the load balancer
resource "gcore_cloud_lbpool" "test" {
  name            = "test-pool-${formatdate("YYYY-MM-DD-hhmm", timestamp())}"
  lb_algorithm    = "ROUND_ROBIN"
  protocol        = "HTTP"
  loadbalancer_id = gcore_cloud_load_balancer.test.id
  project_id      = 379987
  region_id       = 76

  # Simple health monitor
  healthmonitor = {
    delay       = 10
    max_retries = 3
    timeout     = 5
    type        = "TCP"
  }

  depends_on = [gcore_cloud_load_balancer.test]
}

output "load_balancer_id" {
  value = gcore_cloud_load_balancer.test.id
}

output "pool_id" {
  value = gcore_cloud_lbpool.test.id
}

output "pool_name" {
  value = gcore_cloud_lbpool.test.name
}

output "pool_status" {
  value = gcore_cloud_lbpool.test.operating_status
}