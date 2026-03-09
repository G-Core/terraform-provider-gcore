# Upload a custom instance image from URL
resource "gcore_cloud_instance_image" "ubuntu" {
  project_id = 1
  region_id  = 1
  name       = "ubuntu-22.04-custom"
  url        = "https://cloud-images.ubuntu.com/releases/22.04/release/ubuntu-22.04-server-cloudimg-amd64.img"
  os_distro  = "ubuntu"
  os_type    = "linux"
  os_version = "22.04"
  ssh_key    = "allow"
}
