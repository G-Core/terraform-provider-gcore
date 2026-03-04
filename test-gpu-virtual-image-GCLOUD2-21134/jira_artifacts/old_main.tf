terraform {
  required_providers {
    gcore = {
      source  = "G-Core/gcore"
      version = "~> 0.12"
    }
  }
}

provider "gcore" {
  permanent_api_token = var.gcore_api_token
}

variable "gcore_api_token" {
  type        = string
  description = "Gcore API token"
  sensitive   = true
}

variable "project_id" {
  type        = number
  description = "Gcore project ID"
}

variable "region_id" {
  type        = number
  description = "Gcore region ID"
}

# GPU Virtual Image resource - tests async create with task polling
resource "gcore_gpu_virtual_image" "test" {
  project_id   = var.project_id
  region_id    = var.region_id
  name         = "test-gpu-virtual-image-GCLOUD2-21134"
  url          = "http://mirror.noris.net/cirros/0.4.0/cirros-0.4.0-x86_64-disk.img"
  architecture = "x86_64"
  os_type      = "linux"
  os_distro    = "cirros"
  os_version   = "0.4.0"
  ssh_key      = "allow"
}

output "image_id" {
  value       = gcore_gpu_virtual_image.test.id
  description = "The ID of the created GPU virtual image"
}
