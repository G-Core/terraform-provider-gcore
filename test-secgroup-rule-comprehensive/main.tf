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

# Test security group for comprehensive rule testing
resource "gcore_cloud_security_group" "test" {
  project_id = 379987
  region_id  = 76

  security_group = {
    name        = "test-secgroup-comprehensive"
    description = "Comprehensive testing of security group rules"
  }
}

# Test Case 1: Basic TCP rule with port range
resource "gcore_cloud_security_group_rule" "tcp_range" {
  group_id   = gcore_cloud_security_group.test.id
  project_id = 379987
  region_id  = 76

  direction      = "ingress"
  ethertype      = "IPv4"
  protocol       = "tcp"
  port_range_min = 8000
  port_range_max = 8100
  remote_ip_prefix = "192.168.1.0/24"
  description    = "TCP port range 8000-8100"
}

# Test Case 2: UDP single port
resource "gcore_cloud_security_group_rule" "udp_single" {
  group_id   = gcore_cloud_security_group.test.id
  project_id = 379987
  region_id  = 76

  direction      = "ingress"
  ethertype      = "IPv4"
  protocol       = "udp"
  port_range_min = 53
  port_range_max = 53
  remote_ip_prefix = "0.0.0.0/0"
  description    = "DNS UDP"
}

# Test Case 3: ICMP (no ports)
resource "gcore_cloud_security_group_rule" "icmp" {
  group_id   = gcore_cloud_security_group.test.id
  project_id = 379987
  region_id  = 76

  direction        = "ingress"
  ethertype        = "IPv4"
  protocol         = "icmp"
  remote_ip_prefix = "0.0.0.0/0"
  description      = "ICMP ping"
}

# Test Case 4: IPv6 rule
resource "gcore_cloud_security_group_rule" "ipv6" {
  group_id   = gcore_cloud_security_group.test.id
  project_id = 379987
  region_id  = 76

  direction        = "ingress"
  ethertype        = "IPv6"
  protocol         = "tcp"
  port_range_min   = 443
  port_range_max   = 443
  remote_ip_prefix = "::/0"
  description      = "HTTPS IPv6"
}

# Test Case 5: Egress rule
resource "gcore_cloud_security_group_rule" "egress" {
  group_id   = gcore_cloud_security_group.test.id
  project_id = 379987
  region_id  = 76

  direction        = "egress"
  ethertype        = "IPv4"
  protocol         = "tcp"
  port_range_min   = 443
  port_range_max   = 443
  remote_ip_prefix = "0.0.0.0/0"
  description      = "Outbound HTTPS"
}

# Test Case 6: Rule without description (optional field)
resource "gcore_cloud_security_group_rule" "no_desc" {
  group_id   = gcore_cloud_security_group.test.id
  project_id = 379987
  region_id  = 76

  direction        = "ingress"
  ethertype        = "IPv4"
  protocol         = "tcp"
  port_range_min   = 22
  port_range_max   = 22
  remote_ip_prefix = "10.0.0.0/8"
}

# Test Case 7: Protocol "any" (all protocols)
resource "gcore_cloud_security_group_rule" "any_protocol" {
  group_id   = gcore_cloud_security_group.test.id
  project_id = 379987
  region_id  = 76

  direction        = "egress"
  ethertype        = "IPv4"
  protocol         = "any"
  remote_ip_prefix = "0.0.0.0/0"
  description      = "Allow all outbound"
}

# Outputs for verification
output "security_group_id" {
  value = gcore_cloud_security_group.test.id
}

output "rule_ids" {
  value = {
    tcp_range    = gcore_cloud_security_group_rule.tcp_range.id
    udp_single   = gcore_cloud_security_group_rule.udp_single.id
    icmp         = gcore_cloud_security_group_rule.icmp.id
    ipv6         = gcore_cloud_security_group_rule.ipv6.id
    egress       = gcore_cloud_security_group_rule.egress.id
    no_desc      = gcore_cloud_security_group_rule.no_desc.id
    any_protocol = gcore_cloud_security_group_rule.any_protocol.id
  }
}

output "rule_count" {
  value = 7
}
