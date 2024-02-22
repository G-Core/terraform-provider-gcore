provider gcore {
  permanent_api_token = "251$d3361.............1b35f26d8"
}

resource "gcore_network" "network" {
  name       = "network_example"
  mtu        = 1450
  type       = "vxlan"
  region_id  = 1
  project_id = 1
}

resource "gcore_subnet" "subnet" {
  name            = "subnet_example"
  cidr            = "192.168.10.0/24"
  network_id      = gcore_network.network.id
  dns_nameservers = ["8.8.4.4", "1.1.1.1"]

  host_routes {
    destination = "10.0.3.0/24"
    nexthop     = "10.0.0.13"
  }

  gateway_ip = "192.168.10.1"
  region_id  = 1
  project_id = 1
}

resource "gcore_volume" "first_volume" {
  name       = "boot volume"
  type_name  = "ssd_hiiops"
  size       = 5
  image_id   = "f4ce3d30-e29c-4cfd-811f-46f383b6081f"
  region_id  = 1
  project_id = 1
}

resource "gcore_volume" "second_volume" {
  name       = "second volume"
  type_name  = "ssd_hiiops"
  size       = 5
  region_id  = 1
  project_id = 1
}

resource "gcore_instance" "instance" {
  flavor_id = "g1-standard-2-4"
  name      = "test"

  volume {
    source     = "existing-volume"
    volume_id  = gcore_volume.first_volume.id
    boot_index = 0
  }

  volume {
    source     = "existing-volume"
    volume_id  = gcore_volume.second_volume.id
    boot_index = 1
  }

  interface {
    type       = "subnet"
    network_id = gcore_network.network.id
    subnet_id  = gcore_subnet.subnet.id
    security_groups = ["d75db0b2-58f1-4a11-88c6-a932bb897310"]
  }
    //port_id = null
    //ip_address = null
    //fip_source = null
    //existing_fip_id = null
  }

  //deprecated, use metadata_map instead
  //metadata {
  //  key = "some_key"
  //  value = "some_data"
  //}
  metadata_map = {
    some_key = "some_value"
    stage    = "dev"
  }

  configuration {
    key   = "some_key"
    value = "some_data"
  }

  region_id  = 1
  project_id = 1
}

//***
// another one example with one interface to private network and fip to internet
//***

resource "gcore_reservedfixedip" "fixed_ip" {
  project_id       = 1
  region_id        = 1
  type             = "ip_address"
  network_id       = "faf6507b-1ff1-4ebf-b540-befd5c09fe06"
  fixed_ip_address = "192.168.13.6"
  is_vip           = false
}

resource "gcore_volume" "first_volume" {
  name       = "boot volume"
  type_name  = "ssd_hiiops"
  size       = 10
  image_id   = "6dc4e061-6fab-41f3-91a3-0ba848fb32d9"
  project_id = 1
  region_id  = 1
}

resource "gcore_floatingip" "fip" {
  project_id       = 1
  region_id        = 1
  fixed_ip_address = gcore_reservedfixedip.fixed_ip.fixed_ip_address
  port_id          = gcore_reservedfixedip.fixed_ip.port_id
}


resource "gcore_instance" "v" {
  project_id = 1
  region_id  = 1
  name       = "hello"
  flavor_id  = "g1-standard-1-2"

  volume {
    source     = "existing-volume"
    volume_id  = gcore_volume.first_volume.id
    boot_index = 0
  }

  interface {
    type            = "reserved_fixed_ip"
    port_id         = gcore_reservedfixedip.fixed_ip.port_id
    fip_source      = "existing"
    existing_fip_id = gcore_floatingip.fip.id
    security_groups = ["ada84751-fcca-4491-9249-2dfceb321616"]
  }
}
