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

resource "gcore_volume" "boot_volume2" {
  name       = "my-boot-volume2"
  type_name  = "ssd_hiiops"
  size       = 5
  image_id   = data.gcore_image.ubuntu.id
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
}

resource "gcore_instancev2" "instance" {
  flavor_id     = "g1-standard-2-4"
  name          = "my-instance"

  volume {
    volume_id  = gcore_volume.boot_volume.id
    boot_index = 0
  }

  interface {
    type = "external"
    name = "my-external-interface"
  }

  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
}

resource "gcore_instancev2" "instance2" {
  flavor_id     = "g1-standard-2-4"
  name          = "my-instance2"

  volume {
    volume_id  = gcore_volume.boot_volume2.id
    boot_index = 0
  }

  interface {
    type = "external"
    name = "my-external-interface"
  }

  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
}

resource "gcore_servergroup" "servergroup" {
  name       = "default"
  policy     = "affinity"
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
  instance {
      instance_id = gcore_instancev2.instance.id
  }

  instance {
      instance_id = gcore_instancev2.instance2.id
  }
}
