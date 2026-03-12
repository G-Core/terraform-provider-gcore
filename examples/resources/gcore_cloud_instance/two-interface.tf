# Create a private network and subnet
resource "gcore_cloud_network" "network" {
  project_id = 1
  region_id  = 1
  name       = "my-network"
  type       = "vxlan"
}

resource "gcore_cloud_network_subnet" "subnet" {
  project_id      = 1
  region_id       = 1
  name            = "my-subnet"
  cidr            = "192.168.10.0/24"
  network_id      = gcore_cloud_network.network.id
  dns_nameservers = ["8.8.4.4", "1.1.1.1"]
}

# Create an SSH key for instance access
resource "gcore_cloud_ssh_key" "my_key" {
  project_id = 1
  name       = "my-keypair"
  public_key = "ssh-ed25519 ...your public key... user@example.com"
}

# Create a boot volume from an image
resource "gcore_cloud_volume" "boot_volume" {
  project_id = 1
  region_id  = 1
  name       = "my-boot-volume"
  source     = "image"
  image_id   = "6dc4e521-0c72-462f-b2d4-306bcf15e227"
  size       = 20
  type_name  = "ssd_hiiops"
}

# Create an instance with two interfaces: one public, one private
resource "gcore_cloud_instance" "instance_with_two_interfaces" {
  project_id   = 1
  region_id    = 1
  flavor       = "g1-standard-2-4"
  name         = "my-instance"
  ssh_key_name = gcore_cloud_ssh_key.my_key.name

  volumes = [{ volume_id = gcore_cloud_volume.boot_volume.id }]

  interfaces = [
    {
      type      = "external"
      ip_family = "ipv4"
    },
    {
      type       = "subnet"
      network_id = gcore_cloud_network.network.id
      subnet_id  = gcore_cloud_network_subnet.subnet.id
    },
  ]
}
