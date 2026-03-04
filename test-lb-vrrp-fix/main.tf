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

# Network infrastructure
resource "gcore_cloud_network" "test" {
  project_id = local.project_id[0]
  region_id  = data.gcore_cloud_region.rg.id
  name       = "test-net-lb-comprehensive"
}

resource "gcore_cloud_network_subnet" "test" {
  project_id = local.project_id[0]
  region_id  = data.gcore_cloud_region.rg.id
  network_id = gcore_cloud_network.test.id
  name       = "test-subnet-lb-comprehensive"
  cidr       = "10.0.20.0/24"
}

# ===== LOAD BALANCER =====
# Test: vrrp_ips preservation during Update (PATCH + GET pattern)
resource "gcore_cloud_load_balancer" "test" {
  project_id     = local.project_id[0]
  region_id      = data.gcore_cloud_region.rg.id
  flavor         = "lb1-2-4"
  name           = var.lb_name
  vip_network_id = gcore_cloud_network.test.id
  vip_subnet_id  = gcore_cloud_network_subnet.test.id
}

# ===== LOAD BALANCER LISTENER =====
# Test: UpdateAndPoll pattern
resource "gcore_cloud_load_balancer_listener" "test" {
  project_id       = local.project_id[0]
  region_id        = data.gcore_cloud_region.rg.id
  load_balancer_id = gcore_cloud_load_balancer.test.id
  name             = var.listener_name
  protocol         = "HTTP"
  protocol_port    = 80

  # Test computed_optional fields
  timeout_client_data    = var.timeout_client_data
  timeout_member_connect = var.timeout_member_connect
  timeout_member_data    = var.timeout_member_data
}

# ===== LOAD BALANCER POOL =====
# Test: UpdateAndPoll pattern with health_monitor
resource "gcore_cloud_load_balancer_pool" "test" {
  project_id   = local.project_id[0]
  region_id    = data.gcore_cloud_region.rg.id
  listener_id  = gcore_cloud_load_balancer_listener.test.id
  name         = var.pool_name
  lb_algorithm = "ROUND_ROBIN"
  protocol     = "HTTP"

  # Health monitor configuration
  healthmonitor = var.with_health_monitor ? {
    delay       = 10
    max_retries = 3
    timeout     = 5
    type        = "HTTP"
    url_path    = var.health_check_path
    # Let API compute: http_method, max_retries_down
  } : null
}

# ===== LOAD BALANCER POOL MEMBER =====
# Test: AddAndPoll/UpdateAndPoll pattern
resource "gcore_cloud_load_balancer_pool_member" "test" {
  count = var.with_member ? 1 : 0

  project_id    = local.project_id[0]
  region_id     = data.gcore_cloud_region.rg.id
  pool_id       = gcore_cloud_load_balancer_pool.test.id
  address       = var.member_address
  protocol_port = var.member_port
  weight        = var.member_weight
}

# Variables for testing different scenarios
variable "lb_name" {
  default = "test-lb-comprehensive"
}

variable "listener_name" {
  default = "test-listener"
}

variable "pool_name" {
  default = "test-pool"
}

variable "timeout_client_data" {
  default = 50000
}

variable "timeout_member_connect" {
  default = 5000
}

variable "timeout_member_data" {
  default = 50000
}

variable "with_health_monitor" {
  default = true
}

variable "health_check_path" {
  default = "/health"
}

variable "with_member" {
  default = false
}

variable "member_address" {
  default = "10.0.20.100"
}

variable "member_port" {
  default = 80
}

variable "member_weight" {
  default = 1
}

# Outputs
output "lb_id" {
  value = gcore_cloud_load_balancer.test.id
}

output "vrrp_ips" {
  value = gcore_cloud_load_balancer.test.vrrp_ips
}

output "vrrp_ips_count" {
  value = length(gcore_cloud_load_balancer.test.vrrp_ips)
}

output "listener_id" {
  value = gcore_cloud_load_balancer_listener.test.id
}

output "pool_id" {
  value = gcore_cloud_load_balancer_pool.test.id
}

output "member_id" {
  value = var.with_member ? gcore_cloud_load_balancer_pool_member.test[0].id : null
}
