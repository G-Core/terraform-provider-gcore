terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

# Test 1: Security group with no default rules
resource "gcore_cloud_security_group" "test" {
  name        = var.sg_name
  description = var.sg_description
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

# Test 4: Ingress UDP rule (port range 5000-6000)
resource "gcore_cloud_security_group_rule" "ingress_udp" {
  group_id         = gcore_cloud_security_group.test.id
  direction        = "ingress"
  protocol         = "udp"
  ethertype        = "IPv4"
  port_range_min   = 5000
  port_range_max   = 6000
  remote_ip_prefix = "192.168.0.0/16"
  description      = "Allow UDP 5000-6000 from private"
}

variable "sg_name" {
  default = "tf-test-secgroup-comprehensive"
}

variable "sg_description" {
  default = "Comprehensive test of SG + rule resources"
}

output "sg_id" {
  value = gcore_cloud_security_group.test.id
}

output "egress_tcp_rule_id" {
  value = gcore_cloud_security_group_rule.egress_tcp.id
}

output "ingress_icmp_rule_id" {
  value = gcore_cloud_security_group_rule.ingress_icmp.id
}

output "ingress_udp_rule_id" {
  value = gcore_cloud_security_group_rule.ingress_udp.id
}
