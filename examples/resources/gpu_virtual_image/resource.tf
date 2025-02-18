provider "gcore" {
  permanent_api_token = "251$d3361.............1b35f26d8"
}

data "gcore_project" "project" {
  name = "Default"
}

data "gcore_region" "region" {
  name = "Luxembourg-2"
}

resource "gcore_gpu_virtual_image" "example" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
  name         = "my-ubuntu-23.10-x64"
  url          = "https://cloud-images.ubuntu.com/releases/23.10/release/ubuntu-23.10-server-amd64.qcow2"
  architecture = "x86_64"
  os_type      = "linux"
  ssh_key      = "allow"
} 