terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

variable "api_key" {
  type      = string
  sensitive = true
}

provider "gcore" {
  api_key = var.api_key
}

locals {
  project_id = 379987
  region_id  = 76
}

# Load Balancer
resource "gcore_cloud_load_balancer" "test" {
  project_id = local.project_id
  region_id  = local.region_id
  name       = "tf-test-lb-simple-verify"
  flavor     = "lb1-1-2"
}

# Listener
resource "gcore_cloud_load_balancer_listener" "test" {
  project_id       = local.project_id
  region_id        = local.region_id
  load_balancer_id = gcore_cloud_load_balancer.test.id
  name             = "tf-test-listener-simple"
  protocol         = "HTTP"
  protocol_port    = 80
}

# Pool with 2 members
resource "gcore_cloud_load_balancer_pool" "test" {
  project_id       = local.project_id
  region_id        = local.region_id
  load_balancer_id = gcore_cloud_load_balancer.test.id
  listener_id      = gcore_cloud_load_balancer_listener.test.id
  name             = "tf-test-pool-simple"
  protocol         = "HTTP"
  lb_algorithm     = "ROUND_ROBIN"

  # members removed to test GCLOUD2-20778 fix
}

output "pool_id" {
  value = gcore_cloud_load_balancer_pool.test.id
}

output "pool_members" {
  value = gcore_cloud_load_balancer_pool.test.members
}
