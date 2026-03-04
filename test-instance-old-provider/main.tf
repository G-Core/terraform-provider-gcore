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

# Create a volume for the instance
resource "gcore_volume" "test_boot_volume" {
  name       = "test-instance-old-provider-boot"
  type_name  = "standard"
  size       = 10
  image_id   = "8a59956d-ada3-42c6-aaf3-e2c38ba70f99" # ubuntu-24.04-x64
  project_id = var.project_id
  region_id  = var.region_id
}

# Create the instance using the old provider schema
resource "gcore_instancev2" "test_instance" {
  name       = "test-instance-old-provider"
  flavor_id  = "g1-standard-1-2"
  project_id = var.project_id
  region_id  = var.region_id

  volume {
    volume_id = gcore_volume.test_boot_volume.id
  }

  interface {
    type           = "subnet"
    name           = "eth0"
    # Using known network and subnet IDs from test-network-members
    network_id     = "07cbeb68-c7b9-4640-b143-41b029bc57f4"
    subnet_id      = "17d3ae28-21b5-49a7-a688-a4bbc52dd386"
    security_groups = []
  }

  metadata_map = {
    purpose = "terraform-migration-test"
    created_by = "old-provider-test"
  }
}

# Outputs to capture state information
output "instance_id" {
  value = gcore_instancev2.test_instance.id
}

output "instance_name" {
  value = gcore_instancev2.test_instance.name
}

output "instance_flavor_id" {
  value = gcore_instancev2.test_instance.flavor_id
}

output "instance_status" {
  value = gcore_instancev2.test_instance.status
}

output "instance_addresses" {
  value = gcore_instancev2.test_instance.addresses
}

output "volume_id" {
  value = gcore_volume.test_boot_volume.id
}
