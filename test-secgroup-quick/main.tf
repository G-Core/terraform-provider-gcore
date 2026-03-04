terraform {
  required_providers {
    gcore = {
      source  = "gcore/gcore"
      version = "~> 1.0"
    }
  }
}

provider "gcore" {
}

# Security group without inline rules
resource "gcore_cloud_security_group" "test" {
  project_id = 379987
  region_id  = 76

  security_group = {
    name        = "test-rule-drift"
    description = "Testing rule drift with separate resources"
  }
}

# Rule WITH description (user-created)
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
  description      = "Allow HTTPS"
}

# Rule WITH description (user-created)
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
  description      = "Allow SSH"
}

# Rule WITHOUT description (user-created) - this is the critical test case
# Previously this would be deleted on Apply 2, now it should be preserved
resource "gcore_cloud_security_group_rule" "no_desc" {
  group_id   = gcore_cloud_security_group.test.id
  project_id = 379987
  region_id  = 76

  direction        = "ingress"
  ethertype        = "IPv4"
  protocol         = "tcp"
  port_range_min   = 8080
  port_range_max   = 8080
  remote_ip_prefix = "0.0.0.0/0"
  # NO description field - testing TC-06 edge case
}

output "security_group_id" {
  value = gcore_cloud_security_group.test.id
}

output "rule_ids" {
  value = {
    https   = gcore_cloud_security_group_rule.https.id
    ssh     = gcore_cloud_security_group_rule.ssh.id
    no_desc = gcore_cloud_security_group_rule.no_desc.id
  }
}
