# Create an SSH key for baremetal server access
resource "gcore_cloud_ssh_key" "my_key" {
  project_id = 1
  name       = "my-keypair"
  public_key = "ssh-ed25519 ...your public key... user@example.com"
}

# Create a baremetal server with a single external interface
resource "gcore_cloud_baremetal_server" "server" {
  project_id   = 1
  region_id    = 1
  flavor       = "bm1-infrastructure-small"
  name         = "my-bare-metal"
  image_id     = "0f25a566-91a4-4507-aa42-bdd732fb998d"
  ssh_key_name = gcore_cloud_ssh_key.my_key.name

  interfaces = [{
    type      = "external"
    ip_family = "ipv4"
  }]
}
