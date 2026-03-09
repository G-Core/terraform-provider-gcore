resource "gcore_cloud_volume" "boot_volume" {
  project_id = 1
  region_id  = 1

  name      = "my-boot-volume"
  type_name = "ssd_hiiops"
  size      = 5
  image_id  = "your-ubuntu-image-id"
}

resource "gcore_cloud_volume" "boot_volume2" {
  project_id = 1
  region_id  = 1

  name      = "my-boot-volume2"
  type_name = "ssd_hiiops"
  size      = 5
  image_id  = "your-ubuntu-image-id"
}

resource "gcore_cloud_instance" "instance" {
  project_id = 1
  region_id  = 1

  flavor = "g1-standard-2-4"
  name   = "my-instance"

  volumes = [{
    volume_id  = gcore_cloud_volume.boot_volume.id
    boot_index = 0
  }]

  interfaces = [{
    type = "external"
    name = "my-external-interface"
  }]
}

resource "gcore_cloud_instance" "instance2" {
  project_id = 1
  region_id  = 1

  flavor = "g1-standard-2-4"
  name   = "my-instance2"

  volumes = [{
    volume_id  = gcore_cloud_volume.boot_volume2.id
    boot_index = 0
  }]

  interfaces = [{
    type = "external"
    name = "my-external-interface"
  }]
}

resource "gcore_cloud_placement_group" "servergroup" {
  project_id = 1
  region_id  = 1

  name   = "default"
  policy = "affinity"

  instances = [
    { instance_id = gcore_cloud_instance.instance.id },
    { instance_id = gcore_cloud_instance.instance2.id },
  ]
}
