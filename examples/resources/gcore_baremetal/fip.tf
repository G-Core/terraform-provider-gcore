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

resource "gcore_baremetal" "baremetal_with_floating_ip" {
  flavor_id     = "bm1-infrastructure-small"
  name          = "my-baremetal"
  keypair_name  = "my-keypair"
  image_id      = data.gcore_image.ubuntu.id

  interface {
    type    = "reserved_fixed_ip"
    port_id = gcore_reservedfixedip.fixed_ip.port_id

    existing_fip_id = gcore_floatingip.floating_ip.id
  }

  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
}