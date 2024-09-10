data "gcore_image" "windows" {
  name       = "windows-server-2022"
  region_id  = data.gcore_region.region.id
  project_id = data.gcore_project.project.id
}

resource "gcore_volume" "boot_volume_windows" {
  name       = "windows boot volume"
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
  user_data     = "PHBvd2Vyc2hlbGw+CiMgQmUgc3VyZSB0byBzZXQgdGhlIHVzZXJuYW1lIGFuZCBwYXNzd29yZCBvbiB0aGVzZSB0d28gbGluZXMuIE9mIGNvdXJzZSB0aGlzIGlzIG5vdCBhIGdvb2QKIyBzZWN1cml0eSBwcmFjdGljZSB0byBpbmNsdWRlIGEgcGFzc3dvcmQgYXQgY29tbWFuZCBsaW5lLgokVXNlciA9ICJTZWNvbmRVc2VyIgokUGFzc3dvcmQgPSBDb252ZXJ0VG8tU2VjdXJlU3RyaW5nICJzM2NSM3RQQHNzdzByZCIgLUFzUGxhaW5UZXh0IC1Gb3JjZQpOZXctTG9jYWxVc2VyICRVc2VyIC1QYXNzd29yZCAkUGFzc3dvcmQKQWRkLUxvY2FsR3JvdXBNZW1iZXIgLUdyb3VwICJSZW1vdGUgRGVza3RvcCBVc2VycyIgLU1lbWJlciAkVXNlcgpBZGQtTG9jYWxHcm91cE1lbWJlciAtR3JvdXAgIkFkbWluaXN0cmF0b3JzIiAtTWVtYmVyICRVc2VyCjwvcG93ZXJzaGVsbD4="

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