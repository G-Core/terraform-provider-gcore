# Test Case: Update existing rule - change port range
resource "gcore_cloud_security_group_rule" "tcp_range_updated" {
  group_id   = gcore_cloud_security_group.test.id
  project_id = 379987
  region_id  = 76

  direction      = "ingress"
  ethertype      = "IPv4"
  protocol       = "tcp"
  port_range_min = 9000  # Changed from 8000
  port_range_max = 9100  # Changed from 8100
  remote_ip_prefix = "192.168.1.0/24"
  description    = "TCP port range 9000-9100 (updated)"
}

output "updated_rule_id" {
  value = gcore_cloud_security_group_rule.tcp_range_updated.id
}
