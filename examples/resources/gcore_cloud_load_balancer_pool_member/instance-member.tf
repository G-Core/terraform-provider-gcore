resource "gcore_cloud_network" "instance_member_private_network" {
  project_id = 1
  region_id  = 1

  name = "my-private-network"
}

resource "gcore_cloud_network_subnet" "instance_member_private_subnet" {
  project_id = 1
  region_id  = 1

  cidr       = "10.0.0.0/24"
  name       = "my-private-network-subnet"
  network_id = gcore_cloud_network.instance_member_private_network.id
}

resource "gcore_cloud_reserved_fixed_ip" "instance_member_fixed_ip" {
  project_id = 1
  region_id  = 1

  type             = "ip_address"
  network_id       = gcore_cloud_network.instance_member_private_network.id
  subnet_id        = gcore_cloud_network_subnet.instance_member_private_subnet.id
  fixed_ip_address = "10.0.0.11"
  is_vip           = false
}

resource "gcore_cloud_volume" "instance_member_volume" {
  project_id = 1
  region_id  = 1

  name      = "boot volume"
  type_name = "ssd_hiiops"
  size      = 10
  image_id  = "your-ubuntu-image-id"
}

resource "gcore_cloud_instance" "instance_member" {
  project_id = 1
  region_id  = 1

  name_template = "ed-c16-{ip_octets}"
  flavor        = "g1-standard-1-2"

  volumes = [{
    volume_id  = gcore_cloud_volume.instance_member_volume.id
    boot_index = 0
  }]

  interfaces = [{
    type            = "reserved_fixed_ip"
    name            = "my-private-network-interface"
    port_id         = gcore_cloud_reserved_fixed_ip.instance_member_fixed_ip.port_id
    security_groups = [{ id = "your-security-group-id" }]
  }]
}

resource "gcore_cloud_load_balancer_pool_member" "instance_member" {
  project_id = 1
  region_id  = 1

  pool_id = gcore_cloud_load_balancer_pool.http.id

  instance_id   = gcore_cloud_instance.instance_member.id
  address       = gcore_cloud_reserved_fixed_ip.instance_member_fixed_ip.fixed_ip_address
  protocol_port = 80
  weight        = 1
}
