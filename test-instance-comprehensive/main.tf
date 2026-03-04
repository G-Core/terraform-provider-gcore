# Comprehensive Instance Testing - All Volume Union Types
# Tests the volume_id fix and all discriminated union variants

terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

locals {
  project_id = [379987]
  region_id  = 76
}

# Data source for region
data "gcore_cloud_region" "rg" {
  region_id = local.region_id
}

# Ubuntu image ID (pre-discovered)
locals {
  ubuntu_image_id = "f84ddba3-7a5a-4199-931a-250e981d16fb"
}

# =============================================================================
# Test 1: Volume from Image (most common use case)
# =============================================================================
resource "gcore_cloud_instance" "from_image" {
  count = var.test_image_volume ? 1 : 0

  project_id = local.project_id[0]
  region_id  = data.gcore_cloud_region.rg.id

  name   = "test-from-image"
  flavor = "g1-standard-1-2"

  interfaces = [{
    type = "external"
  }]

  volumes = [{
    source     = "image"
    image_id   = local.ubuntu_image_id
    size       = 10
    boot_index = 0
  }]
}

# =============================================================================
# Test 2: New Volume (empty volume)
# =============================================================================
resource "gcore_cloud_instance" "new_volume" {
  count = var.test_new_volume ? 1 : 0

  project_id = local.project_id[0]
  region_id  = data.gcore_cloud_region.rg.id

  name   = "test-new-volume"
  flavor = "g1-standard-1-2"

  interfaces = [{
    type = "external"
  }]

  volumes = [{
    source     = "new-volume"
    size       = 10
    type_name  = "standard"
    boot_index = 0
  }]
}

# =============================================================================
# Test 3: Existing Volume (tests the volume_id fix)
# First create a volume, then attach it to an instance
# =============================================================================
resource "gcore_cloud_volume" "pre_created" {
  count = var.test_existing_volume ? 1 : 0

  project_id = local.project_id[0]
  region_id  = data.gcore_cloud_region.rg.id

  name      = "pre-created-volume"
  size      = 10
  type_name = "standard"
  source    = "image"
  image_id  = local.ubuntu_image_id
}

resource "gcore_cloud_instance" "existing_volume" {
  count = var.test_existing_volume ? 1 : 0

  project_id = local.project_id[0]
  region_id  = data.gcore_cloud_region.rg.id

  name   = "test-existing-volume"
  flavor = "g1-standard-1-2"

  interfaces = [{
    type = "external"
  }]

  volumes = [{
    source    = "existing-volume"
    volume_id = gcore_cloud_volume.pre_created[0].id
    boot_index = 0
  }]

  depends_on = [gcore_cloud_volume.pre_created]
}

# =============================================================================
# Test 4: Multiple Volumes (mixed types)
# =============================================================================
resource "gcore_cloud_instance" "mixed_volumes" {
  count = var.test_mixed_volumes ? 1 : 0

  project_id = local.project_id[0]
  region_id  = data.gcore_cloud_region.rg.id

  name   = "test-mixed-volumes"
  flavor = "g1-standard-1-2"

  interfaces = [{
    type = "external"
  }]

  volumes = [
    {
      source     = "image"
      image_id   = local.ubuntu_image_id
      size       = 10
      boot_index = 0
    },
    {
      source     = "new-volume"
      size       = 20
      type_name  = "ssd_hiiops"
      boot_index = -1
    }
  ]
}

# =============================================================================
# Test 5: Interface Types (external + subnet)
# =============================================================================
resource "gcore_cloud_network" "test" {
  count = var.test_subnet_interface ? 1 : 0

  project_id = local.project_id[0]
  region_id  = data.gcore_cloud_region.rg.id

  name = "test-network"
}

resource "gcore_cloud_network_subnet" "test" {
  count = var.test_subnet_interface ? 1 : 0

  project_id = local.project_id[0]
  region_id  = data.gcore_cloud_region.rg.id

  name       = "test-subnet"
  network_id = gcore_cloud_network.test[0].id
  cidr       = "10.0.1.0/24"
}

resource "gcore_cloud_instance" "subnet_interface" {
  count = var.test_subnet_interface ? 1 : 0

  project_id = local.project_id[0]
  region_id  = data.gcore_cloud_region.rg.id

  name   = "test-subnet-interface"
  flavor = "g1-standard-1-2"

  interfaces = [
    {
      type       = "subnet"
      network_id = gcore_cloud_network.test[0].id
      subnet_id  = gcore_cloud_network_subnet.test[0].id
    }
  ]

  volumes = [{
    source     = "image"
    image_id   = local.ubuntu_image_id
    size       = 10
    boot_index = 0
  }]

  depends_on = [gcore_cloud_network_subnet.test]
}

# =============================================================================
# Variables for selective testing
# =============================================================================
variable "test_image_volume" {
  description = "Test volume from image"
  type        = bool
  default     = false
}

variable "test_new_volume" {
  description = "Test new empty volume"
  type        = bool
  default     = false
}

variable "test_existing_volume" {
  description = "Test existing volume attachment (volume_id fix)"
  type        = bool
  default     = false
}

variable "test_mixed_volumes" {
  description = "Test multiple volumes of different types"
  type        = bool
  default     = false
}

variable "test_subnet_interface" {
  description = "Test subnet interface type"
  type        = bool
  default     = false
}

# =============================================================================
# Outputs
# =============================================================================
output "from_image_id" {
  value = var.test_image_volume ? gcore_cloud_instance.from_image[0].id : null
}

output "new_volume_id" {
  value = var.test_new_volume ? gcore_cloud_instance.new_volume[0].id : null
}

output "existing_volume_instance_id" {
  value = var.test_existing_volume ? gcore_cloud_instance.existing_volume[0].id : null
}

output "mixed_volumes_id" {
  value = var.test_mixed_volumes ? gcore_cloud_instance.mixed_volumes[0].id : null
}

output "subnet_interface_id" {
  value = var.test_subnet_interface ? gcore_cloud_instance.subnet_interface[0].id : null
}
