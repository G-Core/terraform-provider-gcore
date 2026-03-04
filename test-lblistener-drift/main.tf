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
resource "gcore_cloud_load_balancer" "lb" {
  name       = "qa-lb-renamed"
  flavor     = "lb1-1-2"
  project_id = 379987
  region_id  = 76
}

# Test lblistener resource to reproduce drift issue
resource "gcore_cloud_load_balancer_listener" "ls" {
  name             = "qa-ls"
  load_balancer_id = gcore_cloud_load_balancer.lb.id
  protocol         = "HTTP"
  protocol_port    = 81
  project_id       = 379987
  region_id        = 76

  depends_on = [gcore_cloud_load_balancer.lb]
}

output "load_balancer_id" {
  value = gcore_cloud_load_balancer.lb.id
}

output "listener_id" {
  value = gcore_cloud_load_balancer_listener.ls.id
}
