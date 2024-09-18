resource "gcore_reservedfixedip" "external_fixed_ip" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
  type       = "external"
}

resource "gcore_instancev2" "instance_with_reserved_address" {
  flavor_id     = "g1-standard-2-4"
  name          = "my-instance"
  keypair_name  = "my-keypair"

  volume {
    volume_id  = gcore_volume.boot_volume.id
    boot_index = 0
  }

  interface {
    type    = "reserved_fixed_ip"
    name    = "my-reserved-public-interface"
    port_id = gcore_reservedfixedip.external_fixed_ip.port_id
    security_groups = [data.gcore_securitygroup.default.id]
  }

  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
}