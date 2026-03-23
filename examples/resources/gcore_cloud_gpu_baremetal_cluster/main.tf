# Prerequisite resources for GPU bare metal cluster examples

resource "gcore_cloud_network" "network" {
  project_id = 1
  region_id  = 1
  name       = "my-network"
  type       = "vlan"
}

resource "gcore_cloud_network_subnet" "subnet" {
  project_id = 1
  region_id  = 1
  name       = "my-subnet"
  cidr       = "192.168.10.0/24"
  network_id = gcore_cloud_network.network.id
}

resource "gcore_cloud_ssh_key" "keypair" {
  project_id = 1
  name       = "my-keypair"
  public_key = "ssh-ed25519 ...your public key... gcore@gcore.com"
}

data "gcore_cloud_file_share" "vast" {
  project_id = 1
  region_id  = 1
  find_one_by = {
    name = "my-files-share"
  }
}
