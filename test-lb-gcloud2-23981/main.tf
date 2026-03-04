terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

# =============================================================================
# Test: GCLOUD2-23981 - Remove listeners from LB, deprecated timeout fields
# =============================================================================
# Resources under test:
#   1. gcore_cloud_load_balancer - listeners block removed
#   2. gcore_cloud_load_balancer_listener - timeout_member_connect/data removed
#   3. gcore_cloud_load_balancer_pool - timeout_client_data removed
#
# Regression checks:
#   - No perpetual drift on operating_status, provisioning_status, updated_at
#   - No perpetual drift on vrrp_ips, tags_v2
#   - No "Provider produced inconsistent result" errors
# =============================================================================

# --- Load Balancer (NO listeners block) ---
resource "gcore_cloud_load_balancer" "test" {
  project_id = 379987
  region_id  = 76
  name       = var.lb_name
  flavor     = "lb1-1-2"

  tags = {
    "tf-test" = "gcloud2-23981"
    "purpose" = "lb-regression-test"
  }
}

# --- Listener (NO timeout_member_connect / timeout_member_data) ---
resource "gcore_cloud_load_balancer_listener" "test" {
  project_id       = 379987
  region_id        = 76
  load_balancer_id = gcore_cloud_load_balancer.test.id
  name             = var.listener_name
  protocol         = "HTTP"
  protocol_port    = 80
  admin_state_up   = true

  # This field should still work (not deprecated)
  timeout_client_data = 50000
}

# --- Pool (NO timeout_client_data) ---
resource "gcore_cloud_load_balancer_pool" "test" {
  project_id   = 379987
  region_id    = 76
  listener_id  = gcore_cloud_load_balancer_listener.test.id
  name         = var.pool_name
  protocol     = "HTTP"
  lb_algorithm = "ROUND_ROBIN"

  # timeout_member_connect and timeout_member_data should still work on pool
  timeout_member_connect = 5000
  timeout_member_data    = 50000

  healthmonitor = {
    type        = "HTTP"
    delay       = 10
    max_retries = 3
    timeout     = 5
    url_path    = "/health"
  }

  members = [
    {
      address       = "10.0.0.10"
      protocol_port = 8080
      weight        = 1
    }
  ]
}

# --- Variables for update testing ---
variable "lb_name" {
  default = "tf-test-lb-gcloud2-23981"
}

variable "listener_name" {
  default = "tf-test-listener-23981"
}

variable "pool_name" {
  default = "tf-test-pool-23981"
}

# --- Outputs for verification ---
output "lb_id" {
  value = gcore_cloud_load_balancer.test.id
}

output "lb_operating_status" {
  value = gcore_cloud_load_balancer.test.operating_status
}

output "lb_provisioning_status" {
  value = gcore_cloud_load_balancer.test.provisioning_status
}

output "lb_vrrp_ips" {
  value = gcore_cloud_load_balancer.test.vrrp_ips
}

output "listener_id" {
  value = gcore_cloud_load_balancer_listener.test.id
}

output "pool_id" {
  value = gcore_cloud_load_balancer_pool.test.id
}

output "pool_operating_status" {
  value = gcore_cloud_load_balancer_pool.test.operating_status
}
