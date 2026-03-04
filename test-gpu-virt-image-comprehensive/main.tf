terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

# Test 1: Minimal config - only required fields
resource "gcore_cloud_gpu_virtual_cluster_image" "minimal" {
  count = var.test_minimal ? 1 : 0

  name = var.minimal_name
  url  = var.image_url

  project_id = var.project_id
  region_id  = var.region_id
}

# Test 2: Full config - all optional fields
resource "gcore_cloud_gpu_virtual_cluster_image" "full" {
  count = var.test_full ? 1 : 0

  name             = var.full_name
  url              = var.image_url
  architecture     = "x86_64"
  os_type          = "linux"
  os_distro        = "ubuntu"
  os_version       = "22.04"
  ssh_key          = "allow"
  hw_firmware_type = "bios"
  cow_format       = false

  project_id = var.project_id
  region_id  = var.region_id
}

# Test 3: cow_format = true
resource "gcore_cloud_gpu_virtual_cluster_image" "cow_format_true" {
  count = var.test_cow_format ? 1 : 0

  name       = var.cow_format_name
  url        = var.image_url
  cow_format = true

  project_id = var.project_id
  region_id  = var.region_id
}

variable "project_id" {
  type = number
}

variable "region_id" {
  type = number
}

variable "image_url" {
  type    = string
  default = "https://cloud-images.ubuntu.com/minimal/releases/jammy/release/ubuntu-22.04-minimal-cloudimg-amd64.img"
}

variable "minimal_name" {
  type    = string
  default = "tf-test-gpu-img-minimal"
}

variable "full_name" {
  type    = string
  default = "tf-test-gpu-img-full"
}

variable "cow_format_name" {
  type    = string
  default = "tf-test-gpu-img-cow"
}

variable "test_minimal" {
  type    = bool
  default = false
}

variable "test_full" {
  type    = bool
  default = false
}

variable "test_cow_format" {
  type    = bool
  default = false
}

output "minimal_id" {
  value = var.test_minimal ? gcore_cloud_gpu_virtual_cluster_image.minimal[0].id : null
}

output "full_id" {
  value = var.test_full ? gcore_cloud_gpu_virtual_cluster_image.full[0].id : null
}

output "cow_format_id" {
  value = var.test_cow_format ? gcore_cloud_gpu_virtual_cluster_image.cow_format_true[0].id : null
}
