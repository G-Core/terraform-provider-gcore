resource "gcore_instancev2" "instance-with-one-interface" {
  flavor_id     = "g1-standard-2-4"
  name          = "my-instance"
  keypair_name  = "my-keypair"

  volume {
    volume_id  = gcore_volume.boot_volume.id
    boot_index = 0
  }

  interface {
    type = "external"
    name = "my-external-interface"
    security_groups = [gcore_securitygroup.default.id]
  }

  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
}