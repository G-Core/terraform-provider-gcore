terraform {
  required_providers {
    gcore = { source = "gcore/gcore" }
  }
}

provider "gcore" {}

locals {
  image_id = "f84ddba3-7a5a-4199-931a-250e981d16fb"  # Ubuntu 25.04
}

# Boot volume from image
resource "gcore_cloud_volume" "boot" {
  name      = "test-bug2-tags-volume"
  source    = "image"
  image_id  = local.image_id
  size      = 10
  type_name = "standard"
}

# Instance - test tags CRUD operations
resource "gcore_cloud_instance" "test" {
  name   = "test-bug2-tags"
  flavor = "g1-standard-1-2"

  volumes = [{
    volume_id  = gcore_cloud_volume.boot.id
    boot_index = 0
  }]

  interfaces = [{
    type       = "subnet"
    network_id = var.network_id
    subnet_id  = var.subnet_id
  }]

  # Tags - the main focus of this test
  tags = var.tags
}

variable "network_id" {
  default = "cd2c62cd-9763-4766-8d36-6066ed92b3e3"
}

variable "subnet_id" {
  default = "8dfa65b6-bdcc-4fa8-ae53-1720b0b88a63"
}

variable "tags" {
  description = "Instance tags - test create, update, delete"
  type        = map(string)
  default     = {}
}

output "instance_id" {
  value = gcore_cloud_instance.test.id
}

output "tags" {
  value = gcore_cloud_instance.test.tags
}
