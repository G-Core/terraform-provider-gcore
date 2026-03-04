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
  name       = "test-lb-drift-02"
  flavor     = "lb1-1-2"
  project_id = var.project_id
  region_id  = var.region_id

  tags = {
    environment = "test"
    purpose     = "drift-detection"
  }
}

resource "gcore_cloud_load_balancer_listener" "test" {
  name              = "test-listener-drift-02"
  load_balancer_id  = gcore_cloud_load_balancer.test.id
  protocol          = "HTTP"
  protocol_port     = 80
  project_id        = var.project_id
  region_id         = var.region_id

  allowed_cidrs          = ["0.0.0.0/0"]
  connection_limit       = 100000
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
