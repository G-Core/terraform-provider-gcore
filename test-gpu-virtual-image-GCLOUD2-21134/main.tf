terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

variable "project_id" {
  type    = number
  default = 6
}

variable "region_id" {
  type    = number
  default = 76
}

# Test creating a new GPU virtual cluster image
resource "gcore_cloud_gpu_virtual_cluster_image" "test_image" {
  project_id = var.project_id
  region_id  = var.region_id

  name       = "image-luxembourg-2-test"
  url        = "http://mirror.noris.net/cirros/0.4.0/cirros-0.4.0-x86_64-disk.img"

  #architecture     = "x86_64"
  #os_type          = "linux"
  #os_distro        = "cirros"
  #os_version       = "0.4.0"
  #ssh_key          = "allow"
  #hw_firmware_type = "bios"
  #cow_format       = false
}


output "created_image_id" {
  value = gcore_cloud_gpu_virtual_cluster_image.test_image.id
}

output "created_image_name" {
  value = gcore_cloud_gpu_virtual_cluster_image.test_image.name
}

output "created_image_status" {
  value = gcore_cloud_gpu_virtual_cluster_image.test_image.status
}
