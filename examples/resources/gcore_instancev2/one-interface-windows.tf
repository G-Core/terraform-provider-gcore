data "gcore_image" "windows" {
  name       = "windows-server-2022"
  region_id  = data.gcore_region.region.id
  project_id = data.gcore_project.project.id
}

resource "gcore_volume" "boot_volume_windows" {
  name       = "my-windows-boot-volume"
  type_name  = "ssd_hiiops"
  size       = 50
  image_id   = data.gcore_image.windows.id
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
}

resource "gcore_instancev2" "instance" {
  flavor_id     = "g1w-standard-4-8"
  name          = "my-windows-instance"
  password      = "my-s3cR3tP@ssw0rd"

  volume {
    source     = "existing-volume"
    volume_id  = gcore_volume.boot_volume_windows.id
    boot_index = 0
  }

  interface {
    type = "external"
    name = "my-external-interface"
  }

  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
}