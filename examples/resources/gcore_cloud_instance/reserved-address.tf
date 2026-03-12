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

# Reserve a public IP address
resource "gcore_cloud_reserved_fixed_ip" "external_fixed_ip" {
  project_id = 1
  region_id  = 1
  type       = "external"
}

# Create an instance using the reserved public IP
resource "gcore_cloud_instance" "instance_with_reserved_address" {
  project_id   = 1
  region_id    = 1
  flavor       = "g1-standard-2-4"
  name         = "my-instance"
  ssh_key_name = gcore_cloud_ssh_key.my_key.name

  volumes = [{ volume_id = gcore_cloud_volume.boot_volume.id }]

  interfaces = [{
    type    = "reserved_fixed_ip"
    port_id = gcore_cloud_reserved_fixed_ip.external_fixed_ip.port_id
  }]
}
