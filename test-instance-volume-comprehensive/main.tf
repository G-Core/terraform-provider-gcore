terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {
  api_key = var.api_key
}

variable "api_key" {
  type        = string
  description = "Gcore API key"
  sensitive   = true
}

variable "project_id" {
  type    = number
  default = 379987
}

variable "region_id" {
  type    = number
  default = 76
}

variable "test_data_volume" {
  type    = bool
  default = false
}

# Boot volume from image
resource "gcore_cloud_volume" "boot" {
  project_id = var.project_id
  region_id  = var.region_id
  name       = "tf-test-boot-volume-comprehensive"
  source     = "image"
  image_id   = "f84ddba3-7a5a-4199-931a-250e981d16fb" # Ubuntu 25.04
  size       = 10
  type_name  = "standard"
}

# Data volume (optional - for Test 3)
resource "gcore_cloud_volume" "data" {
  count      = var.test_data_volume ? 1 : 0
  project_id = var.project_id
  region_id  = var.region_id
  name       = "tf-test-data-volume-comprehensive"
  source     = "new-volume"
  size       = 10
  type_name  = "standard"
}

# Instance with existing volume(s)
resource "gcore_cloud_instance" "test" {
  project_id = var.project_id
  region_id  = var.region_id
  name       = "tf-test-instance-volumes-comprehensive"
  flavor     = "g1-standard-1-2"

  volumes = var.test_data_volume ? [
    {
      volume_id      = gcore_cloud_volume.boot.id
      boot_index     = 0
      attachment_tag = "boot"
    },
    {
      volume_id      = gcore_cloud_volume.data[0].id
      boot_index     = -1
      attachment_tag = "data"
    }
  ] : [
    {
      volume_id      = gcore_cloud_volume.boot.id
      boot_index     = 0
      attachment_tag = "boot"
    }
  ]

  interfaces = [{
    type      = "external"
    ip_family = "dual"
  }]
}

output "instance_id" {
  value = gcore_cloud_instance.test.id
}

output "boot_volume_id" {
  value = gcore_cloud_volume.boot.id
}

output "data_volume_id" {
  value = var.test_data_volume ? gcore_cloud_volume.data[0].id : null
}
