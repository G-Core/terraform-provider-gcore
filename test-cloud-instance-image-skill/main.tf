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
  default = 379987
}

variable "region_id" {
  type    = number
  default = 76  # Luxembourg
}

variable "image_name" {
  type    = string
  default = "qa-test-image-gcloud2-22615"
}

variable "tags" {
  type    = map(string)
  default = null
}

resource "gcore_cloud_instance_image" "test_image" {
  project_id = var.project_id
  region_id  = var.region_id

  name = var.image_name
  url  = "http://mirror.noris.net/cirros/0.4.0/cirros-0.4.0-x86_64-disk.img"

  os_distro  = "cirros"
  os_version = "0.4.0"
  os_type    = "linux"

  architecture = "x86_64"
  ssh_key      = "allow"

  tags = var.tags
}

output "image_id" {
  value = gcore_cloud_instance_image.test_image.id
}

output "image_status" {
  value = gcore_cloud_instance_image.test_image.status
}

output "image_size" {
  value = gcore_cloud_instance_image.test_image.size
}

output "image_tags" {
  value = gcore_cloud_instance_image.test_image.tags
}
