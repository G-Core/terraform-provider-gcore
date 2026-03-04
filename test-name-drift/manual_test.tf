terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

# =============================================================================
# Manual Test: Security Group Name Drift & Rules Nulling
# =============================================================================
#
# This test verifies:
# 1. Top-level 'name' field doesn't cause drift (Tomer's finding)
# 2. Flat 'security_group_rules' stays null for nested syntax (original issue)
#
# HOW TO TEST:
# 1. Set environment variables:
#    export TF_CLI_CONFIG_FILE="/Users/user/repos/gcore-terraform/.terraformrc"
#    source /Users/user/repos/gcore-terraform/.env
#
# 2. Run terraform apply 3 times:
#    terraform apply -auto-approve
#    terraform apply -auto-approve
#    terraform apply -auto-approve
#
# 3. Check for drift:
#    terraform plan
#
# EXPECTED RESULTS:
# - All applies should show "No changes" after first create
# - Final plan should show "No changes"
# - State should NOT contain top-level 'name' or flat 'security_group_rules'
#
# =============================================================================

resource "gcore_cloud_security_group" "test" {
  project_id = 379987
  region_id  = 76

  security_group = {
    name        = "manual-test-name-drift"
    description = "Manual test for name drift and rules nulling"

    security_group_rules = [
      {
        direction        = "ingress"
        ethertype        = "IPv4"
        protocol         = "tcp"
        port_range_min   = 22
        port_range_max   = 22
        remote_ip_prefix = "0.0.0.0/0"
        description      = "SSH access"
      },
      {
        direction        = "ingress"
        ethertype        = "IPv4"
        protocol         = "tcp"
        port_range_min   = 443
        port_range_max   = 443
        remote_ip_prefix = "0.0.0.0/0"
        description      = "HTTPS access"
      }
    ]
  }
}

# =============================================================================
# Outputs for verification
# =============================================================================

output "id" {
  value = gcore_cloud_security_group.test.id
}

output "top_level_name_value" {
  value       = gcore_cloud_security_group.test.name
  description = "Should be null - user doesn't set this, only security_group.name"
}

output "flat_rules_is_null" {
  value = gcore_cloud_security_group.test.security_group_rules == null
  description = "Should be true - flat field should be null for nested syntax"
}

output "nested_rules_count" {
  value       = length(gcore_cloud_security_group.test.security_group.security_group_rules)
  description = "Should be 2 - user's rules"
}

# =============================================================================
# After testing, clean up with:
#   terraform destroy -auto-approve
# =============================================================================
