terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

# ============================================================================
# Local values
# ============================================================================
locals {
  # Ubuntu 22.04 x64 image ID from Luxembourg-2 region
  ubuntu_image_id = "6343932d-0257-4285-bf89-05060f24095a"
}

# ============================================================================
# Boot volume for instance (tests JIRA volume_id serialization bug)
# ============================================================================
resource "gcore_cloud_volume" "boot" {
  name      = "test-instance-skill-boot"
  size      = 5
  type_name = "standard"
  source    = "image"
  image_id  = local.ubuntu_image_id
}

# ============================================================================
# Test 7: Combined update - CRITICAL JIRA GCLOUD2-21138 bug test
# This tests that BOTH name AND flavor changes are applied in single apply.
# The bug was: early return after specialized endpoints prevented PATCH for name.
# ============================================================================
resource "gcore_cloud_instance" "test_external" {
  # CHANGE 1: Update name (requires standard PATCH)
  name = var.instance_name

  # CHANGE 2: Update flavor (requires specialized /changeflavor endpoint)
  flavor = var.flavor

  interfaces = [{
    type = "external"
  }]

  # Note: boot_index is NOT specified here to avoid import drift.
  # API doesn't return boot_index, so import can't populate it in state.
  # If you need boot_index, specify it only for initial creation.
  volumes = [{
    volume_id = gcore_cloud_volume.boot.id
  }]
}

# ============================================================================
# Variables for testing different scenarios
# ============================================================================
variable "instance_name" {
  default = "test-instance-import"
}

variable "flavor" {
  default = "g1-standard-1-2"
}

# ============================================================================
# Outputs for verification
# ============================================================================
output "instance_id" {
  value = gcore_cloud_instance.test_external.id
}

output "instance_name" {
  value = gcore_cloud_instance.test_external.name
}

output "instance_flavor" {
  value = gcore_cloud_instance.test_external.flavor
}

output "instance_status" {
  value = gcore_cloud_instance.test_external.status
}

output "boot_volume_id" {
  value = gcore_cloud_volume.boot.id
}

output "interface_port_id" {
  value = gcore_cloud_instance.test_external.interfaces[0].port_id
}

output "interface_ip_address" {
  value = gcore_cloud_instance.test_external.interfaces[0].ip_address
}
