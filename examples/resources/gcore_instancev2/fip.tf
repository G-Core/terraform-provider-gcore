resource "gcore_reservedfixedip" "fixed_ip" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
  type       = "subnet"
  network_id = gcore_network.network.id
  subnet_id  = gcore_subnet.subnet.id
}

resource "gcore_floatingip" "floating_ip" {
  project_id       = data.gcore_project.project.id
  region_id        = data.gcore_region.region.id
  fixed_ip_address = gcore_reservedfixedip.fixed_ip.fixed_ip_address
  port_id          = gcore_reservedfixedip.fixed_ip.port_id
}

resource "gcore_instancev2" "instance_with_floating_ip" {
  flavor_id     = "g1-standard-2-4"
  name          = "my-instance"
  keypair_name  = "my-keypair"

  volume {
    volume_id  = gcore_volume.boot_volume.id
    boot_index = 0
  }

  interface {
    type    = "reserved_fixed_ip"
    name    = "my-floating-ip-interface"
    port_id = gcore_reservedfixedip.fixed_ip.port_id

    existing_fip_id = gcore_floatingip.floating_ip.id
    security_groups = [gcore_securitygroup.default.id]
  }

  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
}