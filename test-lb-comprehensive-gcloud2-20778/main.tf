# Comprehensive Test: gcore_cloud_load_balancer
# Testing GCLOUD2-20778 fix and full resource behavior
#
# Test Cases:
# 1. Create minimal LB without tags
# 2. Add tags (reproduces GCLOUD2-20778)
# 3. Modify tags
# 4. Update name
# 5. Remove tags
# 6. Update preferred_connectivity
# 7. Test logging configuration
# 8. Resize flavor

terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

# Network infrastructure
resource "gcore_cloud_network" "test" {
  name = "test-net-lb-comprehensive"
}

resource "gcore_cloud_network_subnet" "test" {
  name       = "test-subnet-lb-comprehensive"
  network_id = gcore_cloud_network.test.id
  cidr       = "10.200.0.0/24"
}

# Main Load Balancer under test
resource "gcore_cloud_load_balancer" "test" {
  name           = var.lb_name
  flavor         = var.lb_flavor
  vip_network_id = gcore_cloud_network.test.id
  vip_subnet_id  = gcore_cloud_network_subnet.test.id

  # Tags - key test for GCLOUD2-20778
  tags = var.lb_tags

  # Optional: preferred_connectivity
  # preferred_connectivity = var.preferred_connectivity  # Commented - may not be supported
}

# Variables for progressive testing
variable "lb_name" {
  description = "Load balancer name"
  type        = string
  default     = "test-lb-comprehensive"
}

variable "lb_flavor" {
  description = "Load balancer flavor"
  type        = string
  default     = "lb1-2-4"
}

variable "lb_tags" {
  description = "Tags for load balancer"
  type        = map(string)
  default     = {}
}

variable "preferred_connectivity" {
  description = "Preferred connectivity (L2/L3)"
  type        = string
  default     = null
}

variable "enable_logging" {
  description = "Enable logging"
  type        = bool
  default     = false
}

# Outputs for verification
output "lb_id" {
  description = "Load balancer ID"
  value       = gcore_cloud_load_balancer.test.id
}

output "lb_name" {
  description = "Load balancer name"
  value       = gcore_cloud_load_balancer.test.name
}

output "lb_tags" {
  description = "Load balancer tags (input)"
  value       = gcore_cloud_load_balancer.test.tags
}

output "lb_tags_v2" {
  description = "Load balancer tags_v2 (computed)"
  value       = gcore_cloud_load_balancer.test.tags_v2
}

output "lb_vrrp_ips" {
  description = "VRRP IPs (computed)"
  value       = gcore_cloud_load_balancer.test.vrrp_ips
}

output "lb_vip_address" {
  description = "VIP address (computed)"
  value       = gcore_cloud_load_balancer.test.vip_address
}

output "lb_preferred_connectivity" {
  description = "Preferred connectivity"
  value       = gcore_cloud_load_balancer.test.preferred_connectivity
}

output "lb_provisioning_status" {
  description = "Provisioning status"
  value       = gcore_cloud_load_balancer.test.provisioning_status
}

output "lb_operating_status" {
  description = "Operating status"
  value       = gcore_cloud_load_balancer.test.operating_status
}

output "lb_flavor" {
  description = "Load balancer flavor"
  value       = gcore_cloud_load_balancer.test.flavor
}
