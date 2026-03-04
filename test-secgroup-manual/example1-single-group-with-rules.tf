terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {
}

# ============================================================
# Security Group Resource
# Backend default rules are automatically deleted during creation
# ============================================================

resource "gcore_cloud_security_group" "test" {
  project_id = 379987
  region_id  = 76

  security_group = {
    name        = "manual-test-secgroup-COMPLETE2"
    description = "Manual testing: default rules auto-deleted, no drift"
  }
}

# ============================================================
# Individual Rule Examples
# ============================================================

resource "gcore_cloud_security_group_rule" "ssh" {
  group_id   = gcore_cloud_security_group.test.id
  project_id = 379987
  region_id  = 76

  direction        = "ingress"
  ethertype        = "IPv4"
  protocol         = "tcp"
  port_range_min   = 22
  port_range_max   = 22
  remote_ip_prefix = "0.0.0.0/0"
  description      = "SSH access"
}

resource "gcore_cloud_security_group_rule" "https" {
  group_id   = gcore_cloud_security_group.test.id
  project_id = 379987
  region_id  = 76

  direction        = "ingress"
  ethertype        = "IPv4"
  protocol         = "tcp"
  port_range_min   = 443
  port_range_max   = 443
  remote_ip_prefix = "0.0.0.0/0"
  description      = "HTTPS traffic"
}

# ============================================================
# For-Each Loop Example: Create Multiple Rules from Map
# ============================================================

variable "ports" {
  description = "Map of port configurations for security group rules"
  type = map(object({
    port        = number
    protocol    = string
    description = string
  }))
  default = {
    http = {
      port        = 80
      protocol    = "tcp"
      description = "HTTP web traffic - SHOULD TRIGGER REPLACE 3"
    }
    mysql = {
      port        = 3306
      protocol    = "tcp"
      description = "MySQL database access"
    }
  }
}

resource "gcore_cloud_security_group_rule" "ports" {
  for_each = var.ports

  group_id   = gcore_cloud_security_group.test.id
  project_id = 379987
  region_id  = 76

  direction        = "ingress"
  ethertype        = "IPv4"
  protocol         = each.value.protocol
  port_range_min   = each.value.port
  port_range_max   = each.value.port
  remote_ip_prefix = "0.0.0.0/0"
  description      = each.value.description
}

# ============================================================
# Outputs
# ============================================================

output "security_group_id" {
  value       = gcore_cloud_security_group.test.id
  description = "ID of the security group"
}

output "security_group_name" {
  value       = gcore_cloud_security_group.test.security_group.name
  description = "Name of the security group"
}

output "security_group_rules" {
  value       = gcore_cloud_security_group.test.security_group_rules
  description = "Should be empty list [] - rules managed separately"
}

output "individual_rule_ids" {
  value = {
    ssh   = gcore_cloud_security_group_rule.ssh.id
    https = gcore_cloud_security_group_rule.https.id
  }
  description = "IDs of individually defined rules"
}

output "loop_rule_ids" {
  value       = { for k, r in gcore_cloud_security_group_rule.ports : k => r.id }
  description = "IDs of rules created via for_each loop"
}

output "total_rules_count" {
  value       = 2 + length(var.ports)
  description = "Total expected rules: 2 individual + 2 from loop = 4 total"
}
