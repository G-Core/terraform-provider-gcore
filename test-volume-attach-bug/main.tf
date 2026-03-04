terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

locals {
  project_id = 379987
  region_id  = 76
  # Ubuntu 22.04 x64 image ID (from cloud_insts_imgs_ls)
  ubuntu_image_id = "6343932d-0257-4285-bf89-05060f24095a"
}

# Boot volume from image
resource "gcore_cloud_volume" "boot" {
  name       = "test-bug-boot-volume"
  source     = "image"
  image_id   = local.ubuntu_image_id
  size       = 10
  type_name  = "standard"
  project_id = local.project_id
  region_id  = local.region_id
}

# Additional volume to attach later
resource "gcore_cloud_volume" "additional" {
  name       = "test-bug-additional-volume"
  source     = "new-volume"
  size       = 5
  type_name  = "standard"
  project_id = local.project_id
  region_id  = local.region_id
}

# Instance with boot volume only (Phase 3 - testing detach)
resource "gcore_cloud_instance" "test" {
  name       = "test-volume-attach-fix"
  flavor     = "g1-standard-1-2"
  project_id = local.project_id
  region_id  = local.region_id

  # Phase 3: Remove additional volume from instance
  volumes = [
    {
      volume_id  = gcore_cloud_volume.boot.id
      boot_index = 0
    }
  ]

  interfaces = [
    {
      type = "external"
    }
  ]
}

output "instance_id" {
  value = gcore_cloud_instance.test.id
}

output "boot_volume_id" {
  value = gcore_cloud_volume.boot.id
}

output "additional_volume_id" {
  value = gcore_cloud_volume.additional.id
}
