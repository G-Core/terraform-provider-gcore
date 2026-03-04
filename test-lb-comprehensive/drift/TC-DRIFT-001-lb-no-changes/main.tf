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
  name       = "test-lb-drift-01"
  flavor     = "lb1-1-2"
  project_id = var.project_id
  region_id  = var.region_id

  tags = {
    environment = "test"
    purpose     = "drift-detection"
  }
}

output "load_balancer_id" {
  value = gcore_cloud_load_balancer.test.id
}
