terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

variable "api_key" {
  type        = string
  description = "Gcore API key"
  sensitive   = true
}

provider "gcore" {
  api_key = var.api_key
}

# Test Load Balancer
resource "gcore_cloud_load_balancer" "test" {
  project_id = 379987
  region_id  = 76
  name       = "tf-test-hm-bug"
  flavor     = "lb1-1-2"
}

# Test Listener
resource "gcore_cloud_load_balancer_listener" "test" {
  project_id       = 379987
  region_id        = 76
  load_balancer_id = gcore_cloud_load_balancer.test.id
  name             = "tf-test-listener-hm"
  protocol         = "HTTP"
  protocol_port    = 80
}

# Test Pool with health monitor - start with delay=10
resource "gcore_cloud_load_balancer_pool" "test" {
  project_id       = 379987
  region_id        = 76
  load_balancer_id = gcore_cloud_load_balancer.test.id
  listener_id      = gcore_cloud_load_balancer_listener.test.id
  name             = "tf-test-pool-hm"
  protocol         = "HTTP"
  lb_algorithm     = "ROUND_ROBIN"

  healthmonitor = {
    type           = "HTTP"
    delay          = 15
    max_retries    = 3
    timeout        = 5
    http_method    = "GET"
    url_path       = "/"
    expected_codes = "200"
  }
}

output "pool_id" {
  value = gcore_cloud_load_balancer_pool.test.id
}

output "pool_healthmonitor_delay" {
  value = gcore_cloud_load_balancer_pool.test.healthmonitor.delay
}
