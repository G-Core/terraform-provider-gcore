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

resource "gcore_cloud_load_balancer" "test" {
  name       = "test-lb-drift-03"
  flavor     = "lb1-1-2"
  project_id = var.project_id
  region_id  = var.region_id

  tags = {
    environment = "test"
    purpose     = "drift-detection"
  }
}

resource "gcore_cloud_load_balancer_listener" "test" {
  name              = "test-listener-drift-03"
  load_balancer_id  = gcore_cloud_load_balancer.test.id
  protocol          = "HTTP"
  protocol_port     = 80
  project_id        = var.project_id
  region_id         = var.region_id
}

resource "gcore_cloud_load_balancer_pool" "test" {
  name             = "test-pool-drift-03"
  lb_algorithm     = "ROUND_ROBIN"
  protocol         = "HTTP"
  listener_id      = gcore_cloud_load_balancer_listener.test.id
  load_balancer_id = gcore_cloud_load_balancer.test.id
  project_id       = var.project_id
  region_id        = var.region_id

  healthmonitor = {
    delay          = 10
    max_retries    = 3
    timeout        = 5
    type           = "HTTP"
    url_path       = "/health"
    expected_codes = "200"
  }

  timeout_client_data    = 50000
  timeout_member_connect = 5000
  timeout_member_data    = 50000
}

output "load_balancer_id" {
  value = gcore_cloud_load_balancer.test.id
}

output "listener_id" {
  value = gcore_cloud_load_balancer_listener.test.id
}

output "pool_id" {
  value = gcore_cloud_load_balancer_pool.test.id
}
