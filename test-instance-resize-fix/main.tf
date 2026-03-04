# Test configuration for GCLOUD2-21138: Instance flavor and volume resize
# This tests:
# 1. Flavor changes use /changeflavor endpoint (not replacement)
# 2. Volume size increases use /extend endpoint (not replacement)

terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
      version = ">= 0.1.0"
    }
  }
}

provider "gcore" {
  # Configuration comes from environment variables:
  # - GCORE_API_KEY
  # - GCORE_CLOUD_PROJECT_ID
  # - GCORE_CLOUD_REGION_ID
}

# Main resource under test: Instance with flavor and volume resize capabilities
resource "gcore_cloud_instance" "test" {
  name       = var.instance_name
  flavor     = var.flavor_id
  project_id = var.project_id
  region_id  = var.region_id

  # Network interfaces - using list attribute syntax
  interfaces = [
    {
      type = "external"  # Simple external interface
    }
  ]

  # Boot volume - size can be increased after creation
  volumes = [
    {
      source     = "image"
      image_id   = "8a59956d-ada3-42c6-aaf3-e2c38ba70f99"  # Ubuntu 24.04 x64
      size       = var.volume_size
      boot_index = 0
    }
  ]
}

# Output for verification
output "instance_id" {
  value = gcore_cloud_instance.test.id
  description = "Instance ID - should remain stable across flavor/volume changes"
}

output "instance_flavor" {
  value = gcore_cloud_instance.test.flavor
  description = "Current instance flavor"
}

output "instance_volumes" {
  value = gcore_cloud_instance.test.volumes
  description = "Current instance volumes with sizes"
}
