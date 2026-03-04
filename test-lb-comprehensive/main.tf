terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

locals {
  project_id = 379987
  region_id  = 76 # Luxembourg-2
}

# ============================================================================
# TEST GROUP 1-3: Load Balancer
# ============================================================================
resource "gcore_cloud_load_balancer" "test_lb" {
  project_id = local.project_id
  region_id  = local.region_id
  name       = "comprehensive-test-lb"
  flavor     = "lb1-1-2"
}

# ============================================================================
# TEST GROUP 4-5: Listener
# ============================================================================
resource "gcore_cloud_load_balancer_listener" "http_listener" {
  project_id       = local.project_id
  region_id        = local.region_id
  load_balancer_id = gcore_cloud_load_balancer.test_lb.id
  name             = "http-listener"
  protocol         = "HTTP"
  protocol_port    = 80
}

# ============================================================================
# TEST GROUP 6-10: Pool with inline configuration
# ============================================================================
resource "gcore_cloud_load_balancer_pool" "test_pool" {
  project_id   = local.project_id
  region_id    = local.region_id
  listener_id  = gcore_cloud_load_balancer_listener.http_listener.id
  name         = "comprehensive-test-pool"
  lb_algorithm = "ROUND_ROBIN"
  protocol     = "HTTP"

  # Health monitor - can be added/removed for testing
  # healthmonitor = {
  #   type        = "HTTP"
  #   delay       = 10
  #   max_retries = 3
  #   timeout     = 5
  #   url_path    = "/health"
  # }

  # Members - can be added/removed for testing
  # members attribute omitted to test no drift

  # Session persistence - can be added/removed for testing
  # session_persistence = {
  #   type = "SOURCE_IP"
  # }
}

# ============================================================================
# OUTPUTS for import testing
# ============================================================================
output "lb_id" {
  value = gcore_cloud_load_balancer.test_lb.id
}

output "listener_id" {
  value = gcore_cloud_load_balancer_listener.http_listener.id
}

output "pool_id" {
  value = gcore_cloud_load_balancer_pool.test_pool.id
}

output "import_lb_cmd" {
  value = "terraform import gcore_cloud_load_balancer.test_lb ${local.project_id}/${local.region_id}/${gcore_cloud_load_balancer.test_lb.id}"
}

output "import_listener_cmd" {
  value = "terraform import gcore_cloud_load_balancer_listener.http_listener ${local.project_id}/${local.region_id}/${gcore_cloud_load_balancer_listener.http_listener.id}"
}

output "import_pool_cmd" {
  value = "terraform import gcore_cloud_load_balancer_pool.test_pool ${local.project_id}/${local.region_id}/${gcore_cloud_load_balancer_pool.test_pool.id}"
}
