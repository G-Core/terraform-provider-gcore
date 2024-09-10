resource "gcore_securitygroup" "web_server_security_group" {
  name       = "web server only"
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  security_group_rules {
    direction      = "egress"
    ethertype      = "IPv4"
    protocol       = "tcp"
    port_range_min = 1
    port_range_max = 24
  }

  security_group_rules {
    direction      = "egress"
    ethertype      = "IPv4"
    protocol       = "tcp"
    port_range_min = 26
    port_range_max = 65535
  }

  security_group_rules {
    direction      = "ingress"
    ethertype      = "IPv4"
    protocol       = "tcp"
    port_range_min = 22
    port_range_max = 22
  }

  security_group_rules {
    direction      = "ingress"
    ethertype      = "IPv4"
    protocol       = "tcp"
    port_range_min = 80
    port_range_max = 80
  }

  security_group_rules {
    direction      = "ingress"
    ethertype      = "IPv4"
    protocol       = "tcp"
    port_range_min = 443
    port_range_max = 443
  }

}

resource "gcore_instancev2" "instance" {
  flavor_id     = "g1-standard-2-4"
  name          = "my-instance"
  keypair_name  = "my-keypair"

  volume {
    source     = "existing-volume"
    volume_id  = gcore_volume.boot_volume.id
    boot_index = 0
  }

  interface {
    type = "external"
    name = "my-external-interface"

    security_groups = [gcore_securitygroup.web_server_security_group.id]
  }

  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
}
