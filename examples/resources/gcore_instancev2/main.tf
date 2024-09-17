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
  type       = "vxlan"
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
  name       = "ubuntu-22.04-x64"
  region_id  = data.gcore_region.region.id
  project_id = data.gcore_project.project.id
}

resource "gcore_volume" "boot_volume" {
  name       = "my-boot-volume"
  type_name  = "ssd_hiiops"
  size       = 5
  image_id   = data.gcore_image.ubuntu.id
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
}

resource "gcore_keypair" "my_keypair" {
  project_id  = data.gcore_project.project.id
  sshkey_name = "my-keypair"
  public_key  = "ssh-ed25519 ...your public key... gcore@gcore.com"
}

data "gcore_securitygroup" "default" {
  name       = "default"
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
}
