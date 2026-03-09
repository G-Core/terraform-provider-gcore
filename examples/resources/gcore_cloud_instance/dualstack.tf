resource "gcore_cloud_instance" "instance_with_dualstack" {
  project_id   = 1
  region_id    = 1
  flavor       = "g1-standard-2-4"
  name         = "my-instance"
  ssh_key_name = gcore_cloud_ssh_key.my_key.name

  volumes = [{ volume_id = gcore_cloud_volume.boot_volume.id }]

  interfaces = [{
    type      = "external"
    ip_family = "dual"
  }]
}

output "addresses" {
  value = gcore_cloud_instance.instance_with_dualstack.addresses
}
