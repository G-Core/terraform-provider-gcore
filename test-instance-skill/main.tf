# Comprehensive test configuration for gcore_cloud_instance
# Based on JIRA GCLOUD2-21138 bugs and requirements
# Tests: creation, flavor change, volume extend, tags, floating IPs, import
#
# DESIGN NOTE: The new provider uses a "bring your own volume" approach.
# Volumes must be created separately using gcore_cloud_volume, then attached.

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

# ============================================================================
# Variables
# ============================================================================

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

variable "instance_name" {
  type    = string
  default = "pr35-comprehensive-test"
}

variable "flavor" {
  type    = string
  default = "g1-standard-1-2"  # Initial flavor, change to g1-standard-2-4 for flavor test
}

variable "volume_size" {
  type    = number
  default = 10  # Extended from 5 to 10
}

# Known Ubuntu 22.04 image ID for Luxembourg-2 region
variable "image_id" {
  type    = string
  default = "6343932d-0257-4285-bf89-05060f24095a"
}

# Tags for testing (JIRA: JSON Merge Patch for tags) - Now a map of strings
variable "tags" {
  type    = map(string)
  default = null
}

# Floating IP config (JIRA: FIP on interface)
variable "enable_floating_ip" {
  type    = bool
  default = false
}

# ============================================================================
# Boot Volume - Created first, then attached to instance
# ============================================================================

resource "gcore_cloud_volume" "boot" {
  project_id = var.project_id
  region_id  = var.region_id

  name      = "${var.instance_name}-boot"
  size      = var.volume_size
  source    = "image"
  image_id  = var.image_id
  type_name = "standard"
}

# ============================================================================
# Main Instance Resource
# ============================================================================

resource "gcore_cloud_instance" "test" {
  project_id = var.project_id
  region_id  = var.region_id

  name   = var.instance_name
  flavor = var.flavor

  # External interface - simple setup for initial creation test
  interfaces = [
    {
      type = "external"
    }
  ]

  # Attach the pre-created boot volume
  volumes = [
    {
      volume_id  = gcore_cloud_volume.boot.id
      boot_index = 0
    }
  ]

  # Tags (JIRA: JSON Merge Patch update/delete)
  # BUG: tags type *map[string]types.String can't handle unknown values
  # tags = var.tags
}

# ============================================================================
# Outputs for verification
# ============================================================================

output "instance_id" {
  description = "Instance ID for import test"
  value       = gcore_cloud_instance.test.id
}

output "instance_status" {
  value = gcore_cloud_instance.test.status
}

output "instance_vm_state" {
  value = gcore_cloud_instance.test.vm_state
}

output "instance_flavor" {
  description = "Current flavor - verify in-place change"
  value       = gcore_cloud_instance.test.flavor
}

output "instance_addresses" {
  value = gcore_cloud_instance.test.addresses
}

output "instance_interfaces" {
  description = "Interface details including port_id"
  value       = gcore_cloud_instance.test.interfaces
}

output "boot_volume_id" {
  description = "Boot volume ID - should have boot_index=0"
  value       = gcore_cloud_instance.test.volumes[0].volume_id
}

output "boot_volume_boot_index" {
  description = "Boot index - should be 0, no drift after import"
  value       = gcore_cloud_instance.test.volumes[0].boot_index
}

output "instance_tags" {
  description = "Current tags - verify update/delete"
  value       = gcore_cloud_instance.test.tags
}

output "volume_size" {
  description = "Boot volume size - verify extend operation"
  value       = gcore_cloud_volume.boot.size
}
