# Test auto-inherit of project_id and region_id
terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {
  # Uses environment variables
}

# Security group with explicit project_id/region_id
resource "gcore_cloud_security_group" "test_inherit" {
  project_id = 379987
  region_id  = 76

  security_group = {
    name        = "test-auto-inherit"
    description = "Test auto-inherit of project_id/region_id"
  }
}

# Rule WITHOUT explicit project_id/region_id - should auto-inherit
resource "gcore_cloud_security_group_rule" "auto_inherit" {
  group_id = gcore_cloud_security_group.test_inherit.id

  direction        = "ingress"
  ethertype        = "IPv4"
  protocol         = "tcp"
  port_range_min   = 8080
  port_range_max   = 8080
  remote_ip_prefix = "0.0.0.0/0"
  description      = "Test auto-inherit - no explicit project/region"
}

output "security_group_id" {
  value = gcore_cloud_security_group.test_inherit.id
}

output "rule_id" {
  value = gcore_cloud_security_group_rule.auto_inherit.id
}

output "rule_project_id" {
  value = gcore_cloud_security_group_rule.auto_inherit.project_id
  description = "Should be 379987 (auto-inherited)"
}

output "rule_region_id" {
  value = gcore_cloud_security_group_rule.auto_inherit.region_id
  description = "Should be 76 (auto-inherited)"
}
