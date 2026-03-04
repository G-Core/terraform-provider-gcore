terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

# Get project data
data "gcore_cloud_projects" "my_projects" {
  name = "default"
}

locals {
  project_id = [for p in data.gcore_cloud_projects.my_projects.items : p.id]
}

data "gcore_cloud_region" "rg" {
  region_id = 76
}

# ===== LOAD BALANCER =====
resource "gcore_cloud_load_balancer" "lb" {
  project_id = local.project_id[0]
  region_id  = data.gcore_cloud_region.rg.id
  flavor     = var.lb_flavor
  name       = var.lb_name
}

# ===== LOAD BALANCER LISTENER =====
resource "gcore_cloud_load_balancer_listener" "ls" {
  project_id       = local.project_id[0]
  region_id        = data.gcore_cloud_region.rg.id
  load_balancer_id = gcore_cloud_load_balancer.lb.id
  name             = var.listener_name
  protocol         = var.listener_protocol
  protocol_port    = var.listener_port

  timeout_client_data    = var.timeout_client_data
  timeout_member_connect = var.timeout_member_connect
  timeout_member_data    = var.timeout_member_data
}

# ===== LOAD BALANCER POOL =====
resource "gcore_cloud_load_balancer_pool" "pool" {
  count = var.create_pool ? 1 : 0

  project_id  = local.project_id[0]
  region_id   = data.gcore_cloud_region.rg.id
  listener_id = gcore_cloud_load_balancer_listener.ls.id

  name         = var.pool_name
  lb_algorithm = var.pool_algorithm
  protocol     = var.pool_protocol


}

output "lb_id" {
  value = gcore_cloud_load_balancer.lb.id
}

output "listener_id" {
  value = gcore_cloud_load_balancer_listener.ls.id
}

output "pool_id" {
  value = var.create_pool ? gcore_cloud_load_balancer_pool.pool[0].id : null
}
