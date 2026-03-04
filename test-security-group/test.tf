terraform {
  required_providers {
    gcore = {
      source  = "gcore/gcore"
      version = "~> 1.0"
    }
  }
}

provider "gcore" {
  # Reads from environment variables:
  # - GCORE_API_KEY
  # - GCORE_CLOUD_PROJECT_ID
  # - GCORE_CLOUD_REGION_ID
}

# Minimal example: Security Group with one rule
# This demonstrates the AWS-like pattern where:
# 1. Security group is created without rules
# 2. Backend adds default rules (usually egress allow-all)
# 3. On second apply, Terraform detects drift and suggests removing default rules
# 4. User can either let Terraform remove them or import them

resource "gcore_cloud_security_group" "minimal" {
  project_id = 379987
  region_id  = 76

  security_group = {
    name        = "minimal-test-sg"
    description = "Minimal security group for testing separate rules"
  }
}

# Add a single rule to allow HTTPS
resource "gcore_cloud_security_group_rule" "https_only" {
  group_id   = gcore_cloud_security_group.minimal.id
  project_id = 379987
  region_id  = 76

  direction        = "ingress"
  ethertype        = "IPv4"
  protocol         = "tcp"
  port_range_min   = 443
  port_range_max   = 443
  remote_ip_prefix = "0.0.0.0/0"
  description      = "Allow HTTPS"
}

output "security_group_id" {
  value = gcore_cloud_security_group.minimal.id
}

output "rule_id" {
  value = gcore_cloud_security_group_rule.https_only.id
}
