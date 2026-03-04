terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

locals {
  image_id = "f84ddba3-7a5a-4199-931a-250e981d16fb"  # Ubuntu 25.04
}

# ============================================================================
# Test for Bug 1: New Floating IP on Private Interface
# When adding floating_ip { source = "new" } to a private interface,
# TF sends PATCH with empty body instead of creating FIP via POST /v1/floatingips
# ============================================================================

# Private network for testing
resource "gcore_cloud_network" "private" {
  name          = "test-kirill-private-net"
  create_router = false  # Avoid router quota limit
}

resource "gcore_cloud_network_subnet" "private" {
  name                      = "test-kirill-private-subnet"
  network_id                = gcore_cloud_network.private.id
  cidr                      = "10.100.0.0/24"
  enable_dhcp               = true
  gateway_ip                = "10.100.0.1"
  dns_nameservers           = ["8.8.8.8", "8.8.4.4"]
  connect_to_network_router = false  # No router
}

# Boot volume for the FIP test instance
resource "gcore_cloud_volume" "fip_boot" {
  name      = "test-kirill-fip-boot"
  size      = 5
  type_name = "standard"
  source    = "image"
  image_id  = local.image_id
}

# Instance with private interface
# Step 1: Create without FIP (attach_new_fip=false)
# Step 2: Update with FIP (attach_new_fip=true) - this triggers the bug
resource "gcore_cloud_instance" "test_fip" {
  name   = "test-kirill-fip"
  flavor = "g1-standard-1-2"

  interfaces = [
    {
      type       = "subnet"
      network_id = gcore_cloud_network.private.id
      subnet_id  = gcore_cloud_network_subnet.private.id
      # Bug 1: When we add floating_ip with source="new" in update,
      # TF sends empty PATCH instead of POST /v1/floatingips
      floating_ip = var.attach_new_fip ? {
        source = "new"
      } : null
    }
  ]

  volumes = [
    {
      volume_id  = gcore_cloud_volume.fip_boot.id
      boot_index = 0
    }
  ]
}

# ============================================================================
# Test for Bug 2: Import with Two Volumes Shows Drift
# After importing VM with 2 volumes, boot_index shows drift from -1 to 0
# ============================================================================

# Pre-create volumes for the two-volume test
resource "gcore_cloud_volume" "boot" {
  count     = var.test_two_volumes ? 1 : 0
  name      = "test-kirill-boot-vol"
  size      = 5
  type_name = "standard"
  source    = "image"
  image_id  = local.image_id
}

resource "gcore_cloud_volume" "data" {
  count     = var.test_two_volumes ? 1 : 0
  name      = "test-kirill-data-vol"
  size      = 10
  type_name = "standard"
  source    = "new-volume"
}

# Instance with two volumes for import test
resource "gcore_cloud_instance" "test_import" {
  count  = var.test_two_volumes ? 1 : 0
  name   = "test-kirill-import"
  flavor = "g1-standard-1-2"

  interfaces = [
    {
      type = "external"
    }
  ]

  # Boot volume (boot_index=0) and Data volume (boot_index=-1)
  volumes = [
    {
      volume_id  = gcore_cloud_volume.boot[0].id
      boot_index = 0
    },
    {
      volume_id  = gcore_cloud_volume.data[0].id
      boot_index = -1
    }
  ]
}

# Variables
variable "attach_new_fip" {
  description = "Set to true to trigger Bug 1 (attach new FIP)"
  type        = bool
  default     = false
}

variable "test_two_volumes" {
  description = "Set to true to test Bug 2 (import with two volumes)"
  type        = bool
  default     = false
}

# Outputs for import test
output "instance_id" {
  value = var.test_two_volumes ? gcore_cloud_instance.test_import[0].id : null
}

output "project_id" {
  value = 379987
}

output "region_id" {
  value = 76
}
