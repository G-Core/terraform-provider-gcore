resource "gcore_cloud_reserved_fixed_ip" "fixed_ip" {
  project_id = 1
  region_id  = 1
  type       = "subnet"
  network_id = gcore_cloud_network.network.id
  subnet_id  = gcore_cloud_network_subnet.subnet.id
}

resource "gcore_cloud_floating_ip" "floating_ip" {
  project_id       = 1
  region_id        = 1
  fixed_ip_address = gcore_cloud_reserved_fixed_ip.fixed_ip.fixed_ip_address
  port_id          = gcore_cloud_reserved_fixed_ip.fixed_ip.port_id
}

resource "gcore_cloud_instance" "instance_with_floating_ip" {
  project_id   = 1
  region_id    = 1
  flavor       = "g1-standard-2-4"
  name         = "my-instance"
  ssh_key_name = gcore_cloud_ssh_key.my_key.name

  volumes = [{ volume_id = gcore_cloud_volume.boot_volume.id }]

  interfaces = [{
    type    = "reserved_fixed_ip"
    port_id = gcore_cloud_reserved_fixed_ip.fixed_ip.port_id

    floating_ip = {
      source              = "existing"
      existing_floating_id = gcore_cloud_floating_ip.floating_ip.id
    }
  }]
}
