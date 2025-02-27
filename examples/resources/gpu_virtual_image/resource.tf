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
  project_id   = data.gcore_project.project.id
  region_id    = data.gcore_region.region.id
  name         = "my-cirros-image"
  url          = "http://mirror.noris.net/cirros/0.4.0/cirros-0.4.0-x86_64-disk.img"
  architecture = "x86_64"
  os_type      = "linux"
  ssh_key      = "allow"
}
