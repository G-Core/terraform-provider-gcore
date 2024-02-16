resource "gcore_network" "instance_member_private_network" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  name = "my-private-network"
}

resource "gcore_subnet" "instance_member_private_subnet" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  cidr       = "10.0.0.0/24"
  name       = "my-private-network-subnet"
  network_id = gcore_network.instance_member_private_network.id
}

resource "gcore_reservedfixedip" "instance_member_fixed_ip" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  type             = "ip_address"
  network_id       = gcore_network.instance_member_private_network.id
  subnet_id        = gcore_subnet.instance_member_private_subnet.id
  fixed_ip_address = "10.0.0.11"
  is_vip           = false
}

data "gcore_image" "ubuntu" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  name       = "ubuntu-23.10"
}

resource "gcore_volume" "instance_member_volume" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  name       = "boot volume"
  type_name  = "ssd_hiiops"
  size       = 10
  image_id   = data.gcore_image.ubuntu.id
}


resource "gcore_instance" "instance_member" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  name_template = "ed-c16-{ip_octets}"
  flavor_id  = "g1-standard-1-2"

  volume {
    source     = "existing-volume"
    volume_id  = gcore_volume.instance_member_volume.id
    boot_index = 0
  }

  interface {
    type            = "reserved_fixed_ip"
    port_id         = gcore_reservedfixedip.instance_member_fixed_ip.port_id
  }
}

resource "gcore_lbmember" "instance_member" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  pool_id       = gcore_lbpool.http.id

  instance_id = gcore_instance.instance_member.id
  address       = gcore_reservedfixedip.instance_member_fixed_ip.fixed_ip_address
  protocol_port = 80
  weight        = 1
}
