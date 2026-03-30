# Create an SSH key for baremetal server access
resource "gcore_cloud_ssh_key" "my_key" {
  project_id = 1
  name       = "my-keypair"
  public_key = "ssh-ed25519 ...your public key... user@example.com"
}

# Reserve a public IP address
resource "gcore_cloud_reserved_fixed_ip" "external_fixed_ip" {
  project_id = 1
  region_id  = 1
  type       = "external"
}

# Create a baremetal server using the reserved public IP
resource "gcore_cloud_baremetal_server" "server_with_reserved_address" {
  project_id   = 1
  region_id    = 1
  flavor       = "bm1-infrastructure-small"
  name         = "my-bare-metal"
  image_id     = "0f25a566-91a4-4507-aa42-bdd732fb998d"
  ssh_key_name = gcore_cloud_ssh_key.my_key.name

  interfaces = [{
    type    = "reserved_fixed_ip"
    port_id = gcore_cloud_reserved_fixed_ip.external_fixed_ip.port_id
  }]
}
