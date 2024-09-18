resource "gcore_instancev2" "instance_with_two_interface" {
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
    security_groups = [data.gcore_securitygroup.default.id]
  }

  interface {
    type = "subnet"
    name = "my-private-interface"
    security_groups = [data.gcore_securitygroup.default.id]

    network_id = gcore_network.network.id
    subnet_id = gcore_subnet.subnet.id
  }
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
}