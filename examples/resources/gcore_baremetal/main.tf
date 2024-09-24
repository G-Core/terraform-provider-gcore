provider gcore {
  permanent_api_token = "251$d3361.............1b35f26d8"
}

data "gcore_project" "project" {
  name = "Default"
}

data "gcore_region" "region" {
  name = "Luxembourg-2"
}

resource "gcore_network" "network" {
  name       = "my-network"
  type       = "vlan"
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
}

resource "gcore_subnet" "subnet" {
  name            = "my-subnet"
  cidr            = "192.168.10.0/24"
  network_id      = gcore_network.network.id

  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
}

data "gcore_image" "ubuntu" {
  name         = "ubuntu-24.04-x64"
  is_baremetal = true
  region_id    = data.gcore_region.region.id
  project_id   = data.gcore_project.project.id
}

resource "gcore_keypair" "keypair" {
  project_id  = data.gcore_project.project.id
  sshkey_name = "my-keypair"
  public_key  = "ssh-ed25519 ...your public key... gcore@gcore.com"
}

data "gcore_image" "windows" {
  name         = "windows-server-standard-2022-ironic"
  is_baremetal = true
  region_id    = data.gcore_region.region.id
  project_id   = data.gcore_project.project.id
}

