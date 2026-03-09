resource "gcore_cloud_reserved_fixed_ip" "external_fixed_ip" {
  project_id = 1
  region_id  = 1
  type       = "external"
}

resource "gcore_cloud_instance" "instance_with_reserved_address" {
  project_id   = 1
  region_id    = 1
  flavor       = "g1-standard-2-4"
  name         = "my-instance"
  ssh_key_name = gcore_cloud_ssh_key.my_key.name

  volumes = [{ volume_id = gcore_cloud_volume.boot_volume.id }]

  interfaces = [{
    type    = "reserved_fixed_ip"
    port_id = gcore_cloud_reserved_fixed_ip.external_fixed_ip.port_id
  }]
}
