terraform {
  required_version = ">= 1.5"

  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

variable "project_id" {
  default = 379987
}

variable "region_id" {
  default = 76
}

variable "pool_name" {
  default = "test-pool-update-01"
}

resource "gcore_cloud_load_balancer" "test" {
  name       = "test-lb-update-01"
  flavor     = "lb1-1-2"
  project_id = var.project_id
  region_id  = var.region_id
}

resource "gcore_cloud_load_balancer_listener" "test" {
  name              = "test-listener-update-01"
  load_balancer_id  = gcore_cloud_load_balancer.test.id
  protocol          = "HTTP"
  protocol_port     = 80
  project_id        = var.project_id
  region_id         = var.region_id
}

resource "gcore_cloud_load_balancer_pool" "test" {
  name             = var.pool_name
  lb_algorithm     = "ROUND_ROBIN"
  protocol         = "HTTP"
  listener_id      = gcore_cloud_load_balancer_listener.test.id
  load_balancer_id = gcore_cloud_load_balancer.test.id
  project_id       = var.project_id
  region_id        = var.region_id

  timeout_client_data    = 50000
  timeout_member_connect = 5000
  timeout_member_data    = 50000
}

output "pool_id" {
  value = gcore_cloud_load_balancer_pool.test.id
}

output "pool_name" {
  value = gcore_cloud_load_balancer_pool.test.name
}
