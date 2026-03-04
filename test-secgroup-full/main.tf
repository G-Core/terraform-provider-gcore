terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

# Test 1: Security group with no default rules
resource "gcore_cloud_security_group" "test" {
  name        = "tf-test-secgroup-full"
  description = "Full test of SG + rule resources"
}

# Test 2: Egress TCP rule (port range)
resource "gcore_cloud_security_group_rule" "egress_tcp" {
  group_id         = gcore_cloud_security_group.test.id
  direction        = "egress"
  protocol         = "tcp"
  ethertype        = "IPv4"
  port_range_min   = 443
  port_range_max   = 443
  remote_ip_prefix = "0.0.0.0/0"
  description      = "Allow HTTPS outbound"
}

# Test 3: Ingress ICMP rule (no ports)
resource "gcore_cloud_security_group_rule" "ingress_icmp" {
  group_id         = gcore_cloud_security_group.test.id
  direction        = "ingress"
  protocol         = "icmp"
  ethertype        = "IPv4"
  remote_ip_prefix = "10.0.0.0/8"
  description      = "Allow ICMP from private"
}

output "sg_id" {
  value = gcore_cloud_security_group.test.id
}

output "egress_rule_id" {
  value = gcore_cloud_security_group_rule.egress_tcp.id
}

output "ingress_rule_id" {
  value = gcore_cloud_security_group_rule.ingress_icmp.id
}
