terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

# Test Load Balancer - testing flavor import fix
resource "gcore_cloud_load_balancer" "lb" {
  project_id = 379987
  region_id  = 76
  name       = "test-import-drift-lb"
  flavor     = "lb1-1-2"
}

# Test Listener
resource "gcore_cloud_load_balancer_listener" "ls" {
  project_id       = 379987
  region_id        = 76
  load_balancer_id = gcore_cloud_load_balancer.lb.id
  name             = "test-listener"
  protocol         = "HTTP"
  protocol_port    = 80
}

# Test Pool - using listener_id which causes the import drift issue
resource "gcore_cloud_load_balancer_pool" "pool" {
  project_id   = 379987
  region_id    = 76
  listener_id  = gcore_cloud_load_balancer_listener.ls.id
  lb_algorithm = "ROUND_ROBIN"
  name         = "test-pool-import"
  protocol     = "HTTP"

  # Healthmonitor removed - should trigger deletion via DELETE endpoint
}

output "pool_id" {
  value = gcore_cloud_load_balancer_pool.pool.id
}

output "listener_id" {
  value = gcore_cloud_load_balancer_listener.ls.id
}

output "lb_id" {
  value = gcore_cloud_load_balancer.lb.id
}
