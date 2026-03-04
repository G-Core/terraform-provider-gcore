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

# Example 1: Web Server Security Group with Separate Rules (Recommended Pattern)
# This is the recommended approach for production use

resource "gcore_cloud_security_group" "web" {
  project_id = 379987
  region_id  = 76

  security_group = {
    name        = "web-server-sg"
    description = "Security group for web servers"
  }
}

# Allow HTTP from anywhere
resource "gcore_cloud_security_group_rule" "http" {
  group_id   = gcore_cloud_security_group.web.id
  project_id = 379987
  region_id  = 76

  direction        = "ingress"
  ethertype        = "IPv4"
  protocol         = "tcp"
  port_range_min   = 80
  port_range_max   = 80
  remote_ip_prefix = "0.0.0.0/0"
  description      = "Allow HTTP from anywhere"
}

# Allow HTTPS from anywhere
resource "gcore_cloud_security_group_rule" "https" {
  group_id   = gcore_cloud_security_group.web.id
  project_id = 379987
  region_id  = 76

  direction        = "ingress"
  ethertype        = "IPv4"
  protocol         = "tcp"
  port_range_min   = 443
  port_range_max   = 443
  remote_ip_prefix = "0.0.0.0/0"
  description      = "Allow HTTPS from anywhere"
}

# Allow SSH from management network only
resource "gcore_cloud_security_group_rule" "ssh" {
  group_id   = gcore_cloud_security_group.web.id
  project_id = 379987
  region_id  = 76

  direction        = "ingress"
  ethertype        = "IPv4"
  protocol         = "tcp"
  port_range_min   = 22
  port_range_max   = 22
  remote_ip_prefix = "10.0.0.0/24"
  description      = "Allow SSH from management network"
}

# Allow all outbound traffic (IPv4)
resource "gcore_cloud_security_group_rule" "egress_ipv4" {
  group_id   = gcore_cloud_security_group.web.id
  project_id = 379987
  region_id  = 76

  direction        = "egress"
  ethertype        = "IPv4"
  protocol         = "any"
  remote_ip_prefix = "0.0.0.0/0"
  description      = "Allow all outbound IPv4 traffic"
}

# Allow all outbound traffic (IPv6)
resource "gcore_cloud_security_group_rule" "egress_ipv6" {
  group_id   = gcore_cloud_security_group.web.id
  project_id = 379987
  region_id  = 76

  direction        = "egress"
  ethertype        = "IPv6"
  protocol         = "any"
  remote_ip_prefix = "::/0"
  description      = "Allow all outbound IPv6 traffic"
}

# Example 2: Database Security Group with Inter-Security-Group Rules

resource "gcore_cloud_security_group" "database" {
  project_id = 379987
  region_id  = 76

  security_group = {
    name        = "database-sg"
    description = "Security group for database servers"
  }
}

# Allow MySQL from web security group
resource "gcore_cloud_security_group_rule" "mysql_from_web" {
  group_id   = gcore_cloud_security_group.database.id
  project_id = 379987
  region_id  = 76

  direction       = "ingress"
  ethertype       = "IPv4"
  protocol        = "tcp"
  port_range_min  = 3306
  port_range_max  = 3306
  remote_group_id = gcore_cloud_security_group.web.id
  description     = "Allow MySQL from web servers"
}

# Example 3: Dynamic Rules using for_each

locals {
  web_ports = {
    http       = { port = 80, description = "HTTP" }
    https      = { port = 443, description = "HTTPS" }
    http_alt   = { port = 8080, description = "Alternative HTTP" }
    https_alt  = { port = 8443, description = "Alternative HTTPS" }
  }
}

resource "gcore_cloud_security_group" "dynamic_web" {
  project_id = 379987
  region_id  = 76

  security_group = {
    name        = "dynamic-web-sg"
    description = "Web server with dynamic rules"
  }
}

resource "gcore_cloud_security_group_rule" "dynamic_web_ports" {
  for_each = local.web_ports

  group_id   = gcore_cloud_security_group.dynamic_web.id
  project_id = 379987
  region_id  = 76

  direction        = "ingress"
  ethertype        = "IPv4"
  protocol         = "tcp"
  port_range_min   = each.value.port
  port_range_max   = each.value.port
  remote_ip_prefix = "0.0.0.0/0"
  description      = "Allow ${each.value.description}"
}

# Outputs
output "web_security_group_id" {
  value       = gcore_cloud_security_group.web.id
  description = "The ID of the web security group"
}

output "database_security_group_id" {
  value       = gcore_cloud_security_group.database.id
  description = "The ID of the database security group"
}

output "dynamic_web_security_group_id" {
  value       = gcore_cloud_security_group.dynamic_web.id
  description = "The ID of the dynamic web security group"
}

output "security_group_rules" {
  value = {
    http  = gcore_cloud_security_group_rule.http.id
    https = gcore_cloud_security_group_rule.https.id
    ssh   = gcore_cloud_security_group_rule.ssh.id
  }
  description = "Map of rule names to rule IDs"
}
