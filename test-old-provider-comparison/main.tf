terraform {
  required_providers {
    gcore = {
      source  = "gcore/gcore"
    }
  }
}

provider "gcore" {
  permanent_api_token = var.gcore_api_key
}

variable "gcore_api_key" {
  type      = string
  sensitive = true
}

variable "project_id" {
  type    = number
  default = 379987
}

variable "region_id" {
  type    = number
  default = 76
}

variable "flavor_id" {
  description = "Instance flavor ID"
  type        = string
  default     = "g1-standard-1-2"
}

variable "volume_size" {
  description = "Boot volume size in GiB"
  type        = number
  default     = 15
}

# Create a volume for the instance
resource "gcore_volume" "test_boot_volume" {
  name       = "test-old-provider-resize-boot"
  type_name  = "standard"
  size       = var.volume_size
  image_id   = "8a59956d-ada3-42c6-aaf3-e2c38ba70f99" # ubuntu-24.04-x64
  project_id = var.project_id
  region_id  = var.region_id
}

# Create the instance using the old provider schema
resource "gcore_instancev2" "test_instance" {
  name       = "test-old-provider-resize"
  flavor_id  = var.flavor_id
  project_id = var.project_id
  region_id  = var.region_id

  volume {
    volume_id = gcore_volume.test_boot_volume.id
  }

  interface {
    type            = "external"
    name            = "eth0"
    security_groups = []
  }

  metadata_map = {
    purpose = "old-provider-resize-test"
  }
}

output "instance_id" {
  value = gcore_instancev2.test_instance.id
}

output "instance_flavor_id" {
  value = gcore_instancev2.test_instance.flavor_id
}

output "volume_id" {
  value = gcore_volume.test_boot_volume.id
}

output "volume_size" {
  value = gcore_volume.test_boot_volume.size
}
