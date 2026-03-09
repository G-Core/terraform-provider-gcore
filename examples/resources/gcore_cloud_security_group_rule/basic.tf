# Allow inbound HTTPS traffic from a specific CIDR
resource "gcore_cloud_security_group_rule" "allow_https" {
  project_id       = 1
  region_id        = 1
  group_id         = gcore_cloud_security_group.web.id
  direction        = "ingress"
  ethertype        = "IPv4"
  protocol         = "tcp"
  port_range_min   = 443
  port_range_max   = 443
  remote_ip_prefix = "10.0.0.0/8"
  description      = "Allow HTTPS from internal network"
}
