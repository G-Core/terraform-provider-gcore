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
# Security Group: Web Tier
# ============================================================

resource "gcore_cloud_security_group" "web" {
  project_id = 379987
  region_id  = 76

  security_group = {
    name        = "test-web-tier"
    description = "Security group for web servers"
  }
}

# ============================================================
# Security Group: Database Tier
# ============================================================

resource "gcore_cloud_security_group" "database" {
  project_id = 379987
  region_id  = 76

  security_group = {
    name        = "test-database-tier"
    description = "Security group for database servers"
  }
}

# ============================================================
# Variables: Define rules for each security group
# ============================================================

variable "web_rules" {
  description = "Rules for web security group"
  type = map(object({
    port        = number
    protocol    = string
    description = string
  }))
  default = {
    http = {
      port        = 80
      protocol    = "tcp"
      description = "HTTP traffic"
    }
    https = {
      port        = 443
      protocol    = "tcp"
      description = "HTTPS traffic"
    }
    ssh = {
      port        = 22
      protocol    = "tcp"
      description = "SSH access"
    }
  }
}

variable "database_rules" {
  description = "Rules for database security group"
  type = map(object({
    port        = number
    protocol    = string
    description = string
    cidr        = string
  }))
  default = {
    postgresql = {
      port        = 5432
      protocol    = "tcp"
      description = "PostgreSQL access"
      cidr        = "10.0.0.0/8"
    }
    mysql = {
      port        = 3306
      protocol    = "tcp"
      description = "MySQL access"
      cidr        = "10.0.0.0/8"
    }
    redis = {
      port        = 6379
      protocol    = "tcp"
      description = "Redis access"
      cidr        = "10.0.0.0/8"
    }
  }
}

# ============================================================
# Rules for Web Security Group (using for_each loop)
# ============================================================

resource "gcore_cloud_security_group_rule" "web_rules" {
  for_each = var.web_rules

  group_id   = gcore_cloud_security_group.web.id
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
# Rules for Database Security Group (using for_each loop)
# ============================================================

resource "gcore_cloud_security_group_rule" "database_rules" {
  for_each = var.database_rules

  group_id   = gcore_cloud_security_group.database.id
  project_id = 379987
  region_id  = 76

  direction        = "ingress"
  ethertype        = "IPv4"
  protocol         = each.value.protocol
  port_range_min   = each.value.port
  port_range_max   = each.value.port
  remote_ip_prefix = each.value.cidr
  description      = each.value.description
}

# ============================================================
# Outputs
# ============================================================

output "web_security_group" {
  value = {
    id    = gcore_cloud_security_group.web.id
    name  = gcore_cloud_security_group.web.security_group.name
    rules = gcore_cloud_security_group.web.security_group_rules
  }
  description = "Web security group details"
}

output "database_security_group" {
  value = {
    id    = gcore_cloud_security_group.database.id
    name  = gcore_cloud_security_group.database.security_group.name
    rules = gcore_cloud_security_group.database.security_group_rules
  }
  description = "Database security group details"
}

output "web_rule_ids" {
  value       = { for k, r in gcore_cloud_security_group_rule.web_rules : k => r.id }
  description = "IDs of web security group rules"
}

output "database_rule_ids" {
  value       = { for k, r in gcore_cloud_security_group_rule.database_rules : k => r.id }
  description = "IDs of database security group rules"
}

output "summary" {
  value = {
    web_group_id      = gcore_cloud_security_group.web.id
    web_rules_count   = length(var.web_rules)
    db_group_id       = gcore_cloud_security_group.database.id
    db_rules_count    = length(var.database_rules)
    total_groups      = 2
    total_rules       = length(var.web_rules) + length(var.database_rules)
  }
  description = "Summary of all resources"
}
