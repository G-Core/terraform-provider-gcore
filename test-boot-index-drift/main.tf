terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

variable "project_id" {
  default = 379987
}

variable "region_id" {
  default = 76
}

# Step 1: Create volumes via Terraform
resource "gcore_cloud_volume" "boot" {
  project_id = var.project_id
  region_id  = var.region_id
  name       = "boot-vol-drift-test"
  size       = 5
  source     = "image"
  image_id   = "8a59956d-ada3-42c6-aaf3-e2c38ba70f99"  # Ubuntu 24.04
  type_name  = "standard"
}

resource "gcore_cloud_volume" "data" {
  project_id = var.project_id
  region_id  = var.region_id
  name       = "data-vol-drift-test"
  size       = 5
  source     = "new-volume"
  type_name  = "standard"
}

# Step 2: Create instance with EXPLICIT boot_index values
# (simulating creation outside TF where API sets correct values)
resource "gcore_cloud_instance" "test" {
  project_id = var.project_id
  region_id  = var.region_id
  name       = "boot-index-drift-test"
  flavor     = "g1-standard-1-2"

  # CRITICAL: Only volume_id specified, NO boot_index
  # This is exactly what QA does after import - should NOT cause drift now
  volumes = [
    { volume_id = gcore_cloud_volume.boot.id },
    { volume_id = gcore_cloud_volume.data.id }
  ]

  interfaces = [{ type = "external" }]
}

output "instance_id" {
  value = gcore_cloud_instance.test.id
}

output "boot_volume_id" {
  value = gcore_cloud_volume.boot.id
}

output "data_volume_id" {
  value = gcore_cloud_volume.data.id
}
