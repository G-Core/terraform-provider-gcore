resource "gcore_cloud_volume" "boot_volume_windows" {
  project_id = 1
  region_id  = 1
  name       = "my-windows-boot-volume"
  source     = "image"
  image_id   = "a2c1681c-94e0-4aab-8fa3-09a8e662d4c0"
  size       = 50
  type_name  = "ssd_hiiops"
}

resource "gcore_cloud_instance" "windows_instance" {
  project_id = 1
  region_id  = 1
  flavor     = "g1w-standard-4-8"
  name       = "my-windows-instance"
  password   = "my-s3cR3tP@ssw0rd"

  volumes = [{ volume_id = gcore_cloud_volume.boot_volume_windows.id }]

  interfaces = [{
    type      = "external"
    ip_family = "ipv4"
  }]
}
