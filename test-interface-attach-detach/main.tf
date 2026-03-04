terraform {
  required_providers {
    gcore = { source = "gcore/gcore" }
  }
}

provider "gcore" {}

locals {
  image_id = "f84ddba3-7a5a-4199-931a-250e981d16fb"  # Ubuntu 25.04

  # Base interface - always present
  base_interfaces = [
    {
      type       = "subnet"
      network_id = var.network_id
      subnet_id  = var.subnet_id
      floating_ip = var.attach_floating_ip && var.create_floating_ip ? {
        source               = "existing"
        existing_floating_id = gcore_cloud_floating_ip.test[0].id
      } : null
    }
  ]

  # External interface - conditional
  external_interface = var.add_second_interface ? [
    {
      type        = "external"
      floating_ip = null
      network_id  = null
      subnet_id   = null
    }
  ] : []
}

# Boot volume from image
resource "gcore_cloud_volume" "boot" {
  name      = "test-mitm-volume"
  source    = "image"
  image_id  = local.image_id
  size      = 10
  type_name = "standard"
}

# Create a floating IP to attach
resource "gcore_cloud_floating_ip" "test" {
  count = var.create_floating_ip ? 1 : 0
}

# Instance - test interface attach/detach
resource "gcore_cloud_instance" "test" {
  name   = "test-mitm-instance"
  flavor = "g1-standard-1-2"

  volumes = [{
    volume_id  = gcore_cloud_volume.boot.id
    boot_index = 0
  }]

  interfaces = concat(local.base_interfaces, local.external_interface)
}

variable "network_id" {
  default = "cd2c62cd-9763-4766-8d36-6066ed92b3e3"
}

variable "subnet_id" {
  default = "8dfa65b6-bdcc-4fa8-ae53-1720b0b88a63"
}

variable "add_second_interface" {
  description = "Add external interface as second interface"
  type        = bool
  default     = false
}

variable "create_floating_ip" {
  description = "Whether to create a floating IP"
  type        = bool
  default     = true
}

variable "attach_floating_ip" {
  description = "Whether to attach the floating IP to the instance"
  type        = bool
  default     = false
}

output "instance_id" {
  value = gcore_cloud_instance.test.id
}

output "floating_ip_id" {
  value = var.create_floating_ip ? gcore_cloud_floating_ip.test[0].id : null
}

output "interfaces" {
  value = gcore_cloud_instance.test.interfaces
}
