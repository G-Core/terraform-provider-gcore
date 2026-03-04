terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {
  # These will be read from environment variables:
  # GCORE_API_KEY, GCORE_CLOUD_PROJECT_ID, GCORE_CLOUD_REGION_ID
}

# Test: Create security group with direct attributes (using nested block as required)
resource "gcore_cloud_security_group" "test" {
  project_id = 379987
  region_id  = 76

  security_group = {
    name        = "test-secgroup-rule-mvp"
    description = "Test security group for rules MVP testing"
  }
}

# Test rule - MVP requires explicit project_id/region_id
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

# Test second rule
resource "gcore_cloud_security_group_rule" "http" {
  group_id   = gcore_cloud_security_group.test.id
  project_id = 379987
  region_id  = 76

  direction        = "ingress"
  ethertype        = "IPv4"
  protocol         = "tcp"
  port_range_min   = 80
  port_range_max   = 80
  remote_ip_prefix = "0.0.0.0/0"
  description      = "Allow HTTP"
}

output "security_group_id" {
  value = gcore_cloud_security_group.test.id
}

output "https_rule_id" {
  value = gcore_cloud_security_group_rule.https.id
}

output "http_rule_id" {
  value = gcore_cloud_security_group_rule.http.id
}

output "rule_details" {
  value = {
    https = {
      id                = gcore_cloud_security_group_rule.https.id
      direction         = gcore_cloud_security_group_rule.https.direction
      protocol          = gcore_cloud_security_group_rule.https.protocol
      port_range_min    = gcore_cloud_security_group_rule.https.port_range_min
      port_range_max    = gcore_cloud_security_group_rule.https.port_range_max
      security_group_id = gcore_cloud_security_group_rule.https.security_group_id
    }
    http = {
      id                = gcore_cloud_security_group_rule.http.id
      direction         = gcore_cloud_security_group_rule.http.direction
      protocol          = gcore_cloud_security_group_rule.http.protocol
      port_range_min    = gcore_cloud_security_group_rule.http.port_range_min
      port_range_max    = gcore_cloud_security_group_rule.http.port_range_max
      security_group_id = gcore_cloud_security_group_rule.http.security_group_id
    }
  }
}
