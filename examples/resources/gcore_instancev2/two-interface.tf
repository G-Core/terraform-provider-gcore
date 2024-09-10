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
  }

  interface {
    type = "subnet"
    name = "my-private-interface"

    network_id = gcore_network.network.id
    subnet_id = gcore_subnet.subnet.id
  }
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
}