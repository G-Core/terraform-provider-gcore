# Create a private network and subnet
resource "gcore_cloud_network" "network" {
  project_id = 1
  region_id  = 1
  name       = "my-network"
  type       = "vlan"
}

resource "gcore_cloud_network_subnet" "subnet" {
  project_id      = 1
  region_id       = 1
  name            = "my-subnet"
  cidr            = "192.168.10.0/24"
  network_id      = gcore_cloud_network.network.id
  dns_nameservers = ["8.8.4.4", "1.1.1.1"]
}

# Create an SSH key for baremetal server access
resource "gcore_cloud_ssh_key" "my_key" {
  project_id = 1
  name       = "my-keypair"
  public_key = "ssh-ed25519 ...your public key... user@example.com"
}

# Reserve a fixed IP on the private subnet
resource "gcore_cloud_reserved_fixed_ip" "fixed_ip" {
  project_id = 1
  region_id  = 1
  type       = "subnet"
  network_id = gcore_cloud_network.network.id
  subnet_id  = gcore_cloud_network_subnet.subnet.id
}

# Create a floating IP and associate it with the fixed IP
resource "gcore_cloud_floating_ip" "floating_ip" {
  project_id       = 1
  region_id        = 1
  fixed_ip_address = gcore_cloud_reserved_fixed_ip.fixed_ip.fixed_ip_address
  port_id          = gcore_cloud_reserved_fixed_ip.fixed_ip.port_id
}

# Create a baremetal server with floating IP for external access
resource "gcore_cloud_baremetal_server" "server_with_floating_ip" {
  project_id   = 1
  region_id    = 1
  flavor       = "bm1-infrastructure-small"
  name         = "my-bare-metal"
  image_id     = "0f25a566-91a4-4507-aa42-bdd732fb998d"
  ssh_key_name = gcore_cloud_ssh_key.my_key.name

  interfaces = [{
    type    = "reserved_fixed_ip"
    port_id = gcore_cloud_reserved_fixed_ip.fixed_ip.port_id

    floating_ip = {
      source               = "existing"
      existing_floating_id = gcore_cloud_floating_ip.floating_ip.id
    }
  }]
}
