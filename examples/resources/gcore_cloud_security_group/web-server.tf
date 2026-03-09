# Create a security group and manage rules as separate resources
resource "gcore_cloud_security_group" "example" {
  project_id  = 1
  region_id   = 1
  name        = "web-security-group"
  description = "Allow HTTP, HTTPS, and outbound traffic"

  tags = {
    environment = "production"
  }
}

resource "gcore_cloud_security_group_rule" "allow_http" {
  project_id     = 1
  region_id      = 1
  group_id       = gcore_cloud_security_group.example.id
  direction      = "ingress"
  ethertype      = "IPv4"
  protocol       = "tcp"
  port_range_min = 80
  port_range_max = 80
  description    = "Allow HTTP"
}

resource "gcore_cloud_security_group_rule" "allow_https" {
  project_id     = 1
  region_id      = 1
  group_id       = gcore_cloud_security_group.example.id
  direction      = "ingress"
  ethertype      = "IPv4"
  protocol       = "tcp"
  port_range_min = 443
  port_range_max = 443
  description    = "Allow HTTPS"
}

resource "gcore_cloud_security_group_rule" "allow_egress_tcp" {
  project_id  = 1
  region_id   = 1
  group_id    = gcore_cloud_security_group.example.id
  direction   = "egress"
  ethertype   = "IPv4"
  protocol    = "tcp"
  description = "Allow all outbound TCP"
}
