# Create a security group, then add rules as separate resources
resource "gcore_cloud_security_group" "web_server" {
  project_id = 1
  region_id  = 1
  name       = "web-server-only"
}

resource "gcore_cloud_security_group_rule" "egress_low" {
  project_id     = 1
  region_id      = 1
  group_id       = gcore_cloud_security_group.web_server.id
  direction      = "egress"
  ethertype      = "IPv4"
  protocol       = "tcp"
  port_range_min = 1
  port_range_max = 24
  description    = "Allow outgoing TCP except SMTP"
}

resource "gcore_cloud_security_group_rule" "egress_high" {
  project_id     = 1
  region_id      = 1
  group_id       = gcore_cloud_security_group.web_server.id
  direction      = "egress"
  ethertype      = "IPv4"
  protocol       = "tcp"
  port_range_min = 26
  port_range_max = 65535
  description    = "Allow outgoing TCP except SMTP"
}

resource "gcore_cloud_security_group_rule" "ssh" {
  project_id     = 1
  region_id      = 1
  group_id       = gcore_cloud_security_group.web_server.id
  direction      = "ingress"
  ethertype      = "IPv4"
  protocol       = "tcp"
  port_range_min = 22
  port_range_max = 22
  description    = "Allow SSH"
}

resource "gcore_cloud_security_group_rule" "http" {
  project_id     = 1
  region_id      = 1
  group_id       = gcore_cloud_security_group.web_server.id
  direction      = "ingress"
  ethertype      = "IPv4"
  protocol       = "tcp"
  port_range_min = 80
  port_range_max = 80
  description    = "Allow HTTP"
}

resource "gcore_cloud_security_group_rule" "https" {
  project_id     = 1
  region_id      = 1
  group_id       = gcore_cloud_security_group.web_server.id
  direction      = "ingress"
  ethertype      = "IPv4"
  protocol       = "tcp"
  port_range_min = 443
  port_range_max = 443
  description    = "Allow HTTPS"
}

resource "gcore_cloud_instance" "instance_with_custom_sg" {
  project_id   = 1
  region_id    = 1
  flavor       = "g1-standard-2-4"
  name         = "my-instance"
  ssh_key_name = gcore_cloud_ssh_key.my_key.name

  volumes = [{ volume_id = gcore_cloud_volume.boot_volume.id }]

  interfaces = [{
    type      = "external"
    ip_family = "ipv4"
    security_groups = [{
      id = gcore_cloud_security_group.web_server.id
    }]
  }]

  security_groups = [{
    id = gcore_cloud_security_group.web_server.id
  }]
}
