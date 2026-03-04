terraform {
  required_providers {
    gcore = {
      source  = "gcore/gcore"
      version = "~> 1.0"
    }
  }
}

provider "gcore" {
  # Uses environment variables from .env file
}

# First create a load balancer
resource "gcore_cloud_load_balancer" "test" {
  name       = "test-lb-listener-${formatdate("YYYY-MM-DD-hhmm", timestamp())}"
  flavor     = "lb1-1-2"
  project_id = 379987
  region_id  = 76

  tags = {
    test = "lblistener"
  }
}

# Test lblistener resource with *AndPoll methods
resource "gcore_cloud_load_balancer_listener" "test" {
  name            = "test-listener-andpoll"
  loadbalancer_id = gcore_cloud_load_balancer.test.id
  protocol        = "HTTP"
  protocol_port   = 80
  project_id      = 379987
  region_id       = 76

  depends_on = [gcore_cloud_load_balancer.test]
}

output "load_balancer_id" {
  value = gcore_cloud_load_balancer.test.id
}

output "listener_id" {
  value = gcore_cloud_load_balancer_listener.test.id
}