# Upload a custom image for GPU virtual clusters
resource "gcore_cloud_gpu_virtual_cluster_image" "ubuntu" {
  project_id   = 1
  region_id    = 1
  name         = "ubuntu-gpu-virtual"
  url          = "https://cloud-images.ubuntu.com/releases/22.04/release/ubuntu-22.04-server-cloudimg-amd64.img"
  architecture = "x86_64"
  os_type      = "linux"
  ssh_key      = "allow"
}
