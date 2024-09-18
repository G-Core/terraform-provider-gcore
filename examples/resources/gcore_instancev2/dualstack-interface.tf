resource "gcore_instancev2" "instance_with_dualstack" {
  flavor_id     = "g1-standard-2-4"
  name          = "my-instance"
  keypair_name  = "my-keypair"

  volume {
    volume_id  = gcore_volume.boot_volume.id
    boot_index = 0
  }

  interface {
    type      = "external"
    ip_family = "dual"
    name      = "my-external-interface"
    security_groups = [data.gcore_securitygroup.default.id]
  }

  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
}

output "addresses" {
  value = gcore_instancev2.instance_with_dualstack.addresses
}