terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {
}

# TC-01: Basic security group creation
# Expected: Default rules are automatically deleted during Create
resource "gcore_cloud_security_group" "test" {
  project_id = 379987
  region_id  = 76

  security_group = {
    name        = "test-secgroup-skill-updated"
    description = "Updated description for testing"
  }
}

# TC-03: Separate rule resources
# Expected: Rules are created independently and preserved
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
  description      = "Allow HTTPS traffic"
}

# SSH rule removed for testing deletion
# resource "gcore_cloud_security_group_rule" "ssh" {
#   group_id   = gcore_cloud_security_group.test.id
#   project_id = 379987
#   region_id  = 76
#
#   direction        = "ingress"
#   ethertype        = "IPv4"
#   protocol         = "tcp"
#   port_range_min   = 22
#   port_range_max   = 22
#   remote_ip_prefix = "0.0.0.0/0"
#   description      = "Allow SSH access"
# }

# TC-04: Rule without description (edge case)
resource "gcore_cloud_security_group_rule" "custom" {
  group_id   = gcore_cloud_security_group.test.id
  project_id = 379987
  region_id  = 76

  direction        = "ingress"
  ethertype        = "IPv4"
  protocol         = "tcp"
  port_range_min   = 8080
  port_range_max   = 8080
  remote_ip_prefix = "0.0.0.0/0"
  # No description - testing edge case
}

output "security_group_id" {
  value = gcore_cloud_security_group.test.id
}

output "security_group_rules" {
  value = gcore_cloud_security_group.test.security_group_rules
  description = "Should be empty list after default rule deletion"
}

output "rule_ids" {
  value = {
    https  = gcore_cloud_security_group_rule.https.id
    # ssh removed for testing
    custom = gcore_cloud_security_group_rule.custom.id
  }
}
