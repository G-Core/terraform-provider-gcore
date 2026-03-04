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

variable "api_key" {
  type        = string
  description = "Gcore API key"
  sensitive   = true
}

variable "project_id" {
  type        = number
  description = "Project ID"
  default     = 379987
}

variable "region_id" {
  type        = number
  description = "Region ID"
  default     = 76
}

# Create a boot volume from an image
resource "gcore_cloud_volume" "boot" {
  project_id = var.project_id
  region_id  = var.region_id
  name       = "tf-test-boot-volume"
  source     = "image"
  image_id   = "f84ddba3-7a5a-4199-931a-250e981d16fb" # Ubuntu 25.04
  size       = 10
  type_name  = "standard"
}

# Create an instance with the existing volume
resource "gcore_cloud_instance" "test" {
  project_id = var.project_id
  region_id  = var.region_id
  name       = "tf-test-instance-existing-volume"
  flavor     = "g1-standard-1-2"

  # Attach existing volumes by ID
  volumes = [{
    volume_id  = gcore_cloud_volume.boot.id
    boot_index = 0
  }]

  interfaces = [{
    type      = "external"
    ip_family = "dual"
  }]
}

output "instance_id" {
  value = gcore_cloud_instance.test.id
}

output "volume_id" {
  value = gcore_cloud_volume.boot.id
}
