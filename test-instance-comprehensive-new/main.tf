# Comprehensive Cloud Instance Testing
# Tests: Create, Import, Drift, Update (name, tags, flavor), Volume attach/detach

terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

locals {
  image_id         = "6343932d-0257-4285-bf89-05060f24095a" # ubuntu-22.04-x64
  network_id       = "cd2c62cd-9763-4766-8d36-6066ed92b3e3"
  subnet_id        = "4f144cc7-c377-445d-9c23-fa6576f1b945"
  security_group_1 = "9fa59dfb-df95-4860-965b-455556cbe7eb" # default
}

# Boot volume for the test instance
resource "gcore_cloud_volume" "test_boot" {
  name      = "tf-test-boot-comprehensive"
  size      = 10
  source    = "image"
  image_id  = local.image_id
  type_name = "standard"
}

# Second volume for attach/detach testing
resource "gcore_cloud_volume" "test_data" {
  name      = "tf-test-data-comprehensive"
  size      = 5
  source    = "new-volume"
  type_name = "standard"
}

# Test instance with subnet interface
resource "gcore_cloud_instance" "test" {
  name   = var.instance_name
  flavor = var.flavor

  interfaces = [{
    type       = "subnet"
    network_id = local.network_id
    subnet_id  = local.subnet_id
  }]

  volumes = var.attach_data_volume ? [
    {
      volume_id  = gcore_cloud_volume.test_boot.id
      boot_index = 0
    },
    {
      volume_id  = gcore_cloud_volume.test_data.id
      boot_index = 1
    }
  ] : [
    {
      volume_id  = gcore_cloud_volume.test_boot.id
      boot_index = 0
    }
  ]
}

variable "instance_name" {
  description = "Instance name - for update testing"
  default     = "tf-test-comprehensive-v1"
}

variable "flavor" {
  description = "Instance flavor - for resize testing"
  default     = "g1-standard-1-2"
}

variable "attach_data_volume" {
  description = "Whether to attach data volume"
  type        = bool
  default     = false
}

output "instance_id" {
  value = gcore_cloud_instance.test.id
}

output "instance_name" {
  value = gcore_cloud_instance.test.name
}

output "instance_status" {
  value = gcore_cloud_instance.test.status
}

output "instance_vm_state" {
  value = gcore_cloud_instance.test.vm_state
}

output "instance_flavor" {
  value = gcore_cloud_instance.test.flavor
}

output "volumes_attached" {
  value = [for v in gcore_cloud_instance.test.volumes : v.volume_id]
}

output "interface_ip" {
  value = gcore_cloud_instance.test.interfaces[0].ip_address
}

output "interface_port_id" {
  value = gcore_cloud_instance.test.interfaces[0].port_id
}

output "boot_volume_id" {
  value = gcore_cloud_volume.test_boot.id
}

output "data_volume_id" {
  value = gcore_cloud_volume.test_data.id
}
