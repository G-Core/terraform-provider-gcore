# Test Configuration for Remaining Union Variants
# Tests: any_subnet, vm_state change, flavor change, MITM verification

terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

locals {
  image_id         = "f84ddba3-7a5a-4199-931a-250e981d16fb" # Ubuntu 25.04
  flavor_initial   = "g1-standard-1-2"
  flavor_upgraded  = "g1-standard-2-4"
  network_id       = "cd2c62cd-9763-4766-8d36-6066ed92b3e3"
  subnet_id        = "4f144cc7-c377-445d-9c23-fa6576f1b945"
  security_group_1 = "9fa59dfb-df95-4860-965b-455556cbe7eb" # default
}

# Boot volume for the test instance
resource "gcore_cloud_volume" "test_boot" {
  name       = "tf-test-boot-remaining"
  size       = 10
  source     = "image"
  image_id   = local.image_id
  type_name  = "standard"
}

# Test 3: any_subnet interface variant
# This tests the any_subnet union type which auto-selects a subnet from the network
resource "gcore_cloud_instance" "test_any_subnet" {
  name   = "tf-test-instance-any-subnet"
  flavor = var.flavor

  interfaces = [{
    type       = "any_subnet"
    network_id = local.network_id
    # Note: subnet_id NOT specified - API will auto-select
  }]

  volumes = [{
    volume_id  = gcore_cloud_volume.test_boot.id
    boot_index = 0
  }]
}

variable "flavor" {
  description = "Instance flavor - can be changed to test resize"
  default     = "g1-standard-1-2"
}

variable "vm_state" {
  description = "VM state - active or stopped"
  default     = "active"
}

output "instance_id" {
  value = gcore_cloud_instance.test_any_subnet.id
}

output "instance_status" {
  value = gcore_cloud_instance.test_any_subnet.status
}

output "instance_vm_state" {
  value = gcore_cloud_instance.test_any_subnet.vm_state
}

output "interface_ip_address" {
  value = gcore_cloud_instance.test_any_subnet.interfaces[0].ip_address
}

output "interface_port_id" {
  value = gcore_cloud_instance.test_any_subnet.interfaces[0].port_id
}
