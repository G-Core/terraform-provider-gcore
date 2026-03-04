terraform {
  required_providers {
    gcore = { source = "gcore/gcore" }
  }
}

provider "gcore" {}

locals {
  image_id = "f84ddba3-7a5a-4199-931a-250e981d16fb"  # Ubuntu 25.04
}

# Use existing network (to avoid subnet resource bug)
variable "network_id" {
  default = "cd2c62cd-9763-4766-8d36-6066ed92b3e3"  # From previous tests
}

variable "subnet_id" {
  default = "8dfa65b6-bdcc-4fa8-ae53-1720b0b88a63"  # From previous tests
}

# Boot volume
resource "gcore_cloud_volume" "boot" {
  name      = "test-comprehensive-boot"
  source    = "image"
  image_id  = local.image_id
  size      = 10
  type_name = "standard"
}

# Main instance under test
resource "gcore_cloud_instance" "test" {
  name   = var.instance_name
  flavor = var.flavor

  volumes = [{
    volume_id  = gcore_cloud_volume.boot.id
    boot_index = 0
  }]

  # Interfaces - test floating IP attach/detach
  interfaces = [
    {
      type       = "subnet"
      network_id = var.network_id
      subnet_id  = var.subnet_id
      floating_ip = var.floating_ip_id != null ? {
        source               = "existing"
        existing_floating_id = var.floating_ip_id
      } : null
    }
  ]

  # Test tags (create, update, delete)
  tags = var.tags

  # Test servergroup (add/remove without replacement)
  servergroup_id = var.servergroup_id
}

variable "instance_name" {
  default = "test-comprehensive-instance"
}

variable "flavor" {
  default = "g1-standard-1-2"
}

variable "tags" {
  description = "Instance tags"
  type        = map(string)
  default     = {}
}

variable "servergroup_id" {
  description = "Placement group ID"
  type        = string
  default     = null
}

variable "floating_ip_id" {
  description = "Floating IP ID to attach to subnet interface"
  type        = string
  default     = null
}

# Outputs for verification
output "instance_id" {
  value = gcore_cloud_instance.test.id
}

output "addresses" {
  value = gcore_cloud_instance.test.addresses
}

output "interfaces" {
  value = gcore_cloud_instance.test.interfaces
}

output "tags" {
  value = gcore_cloud_instance.test.tags
}

output "servergroup_id" {
  value = gcore_cloud_instance.test.servergroup_id
}

output "interface_floating_ip" {
  value = try(gcore_cloud_instance.test.interfaces[0].floating_ip, null)
}
