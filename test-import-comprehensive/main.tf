terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

# Ubuntu 22.04 x64 image ID (from API)
locals {
  ubuntu_image_id = "6343932d-0257-4285-bf89-05060f24095a"
  project_id      = 379987
  region_id       = 76
}

# ===========================================
# BOOT VOLUMES (pre-created)
# ===========================================

resource "gcore_cloud_volume" "minimal_boot" {
  count     = var.test_minimal ? 1 : 0
  name      = "import-test-minimal-boot"
  source    = "image"
  image_id  = local.ubuntu_image_id
  type_name = "ssd_hiiops"
  size      = 10
}

resource "gcore_cloud_volume" "two_volumes_boot" {
  count     = var.test_two_volumes ? 1 : 0
  name      = "import-test-two-volumes-boot"
  source    = "image"
  image_id  = local.ubuntu_image_id
  type_name = "ssd_hiiops"
  size      = 10
}

resource "gcore_cloud_volume" "two_volumes_data" {
  count     = var.test_two_volumes ? 1 : 0
  name      = "import-test-two-volumes-data"
  source    = "new-volume"
  type_name = "standard"
  size      = 5
}

resource "gcore_cloud_volume" "fip_boot" {
  count     = var.test_floating_ip ? 1 : 0
  name      = "import-test-fip-boot"
  source    = "image"
  image_id  = local.ubuntu_image_id
  type_name = "ssd_hiiops"
  size      = 10
}

resource "gcore_cloud_volume" "tags_boot" {
  count     = var.test_tags ? 1 : 0
  name      = "import-test-tags-boot"
  source    = "image"
  image_id  = local.ubuntu_image_id
  type_name = "ssd_hiiops"
  size      = 10
}

resource "gcore_cloud_volume" "servergroup_boot" {
  count     = var.test_servergroup ? 1 : 0
  name      = "import-test-servergroup-boot"
  source    = "image"
  image_id  = local.ubuntu_image_id
  type_name = "ssd_hiiops"
  size      = 10
}

# ===========================================
# Test 1: Minimal instance (just required fields)
# ===========================================

resource "gcore_cloud_instance" "minimal" {
  count = var.test_minimal ? 1 : 0

  name       = "import-test-minimal"
  flavor     = "g1-standard-1-2"
  project_id = local.project_id
  region_id  = local.region_id

  interfaces = [{
    type = "external"
  }]

  volumes = [{
    volume_id  = gcore_cloud_volume.minimal_boot[0].id
    boot_index = 0
  }]
}

# ===========================================
# Test 2: Instance with two volumes (boot + data)
# ===========================================

resource "gcore_cloud_instance" "two_volumes" {
  count = var.test_two_volumes ? 1 : 0

  name       = "import-test-two-volumes"
  flavor     = "g1-standard-1-2"
  project_id = local.project_id
  region_id  = local.region_id

  interfaces = [{
    type = "external"
  }]

  volumes = [
    {
      volume_id  = gcore_cloud_volume.two_volumes_boot[0].id
      boot_index = 0
    },
    {
      volume_id  = gcore_cloud_volume.two_volumes_data[0].id
      boot_index = -1
    }
  ]
}

# ===========================================
# Network for private interface tests
# ===========================================

resource "gcore_cloud_network" "test" {
  count = var.test_floating_ip ? 1 : 0
  name  = "import-test-network"
  type  = "vlan"
}

resource "gcore_cloud_network_subnet" "test" {
  count      = var.test_floating_ip ? 1 : 0
  name       = "import-test-subnet"
  cidr       = "192.168.100.0/24"
  network_id = gcore_cloud_network.test[0].id
}

resource "gcore_cloud_network_router" "test" {
  count = var.test_floating_ip ? 1 : 0
  name  = "import-test-router"
  external_gateway_info = {
    type        = "default"
    enable_snat = true
  }
  interfaces = [
    {
      subnet_id = gcore_cloud_network_subnet.test[0].id
      type      = "subnet"
    }
  ]
}

# ===========================================
# Test 3: Instance with floating IP
# ===========================================

resource "gcore_cloud_instance" "floating_ip" {
  count = var.test_floating_ip ? 1 : 0

  name       = "import-test-fip"
  flavor     = "g1-standard-1-2"
  project_id = local.project_id
  region_id  = local.region_id

  interfaces = [{
    type       = "subnet"
    network_id = gcore_cloud_network.test[0].id
    subnet_id  = gcore_cloud_network_subnet.test[0].id
    floating_ip = {
      source = "new"
    }
  }]

  volumes = [{
    volume_id  = gcore_cloud_volume.fip_boot[0].id
    boot_index = 0
  }]

  depends_on = [gcore_cloud_network_router.test]
}

# ===========================================
# Test 4: Instance with tags
# ===========================================

resource "gcore_cloud_instance" "with_tags" {
  count = var.test_tags ? 1 : 0

  name       = "import-test-tags"
  flavor     = "g1-standard-1-2"
  project_id = local.project_id
  region_id  = local.region_id

  tags = {
    env     = "test"
    project = "import-testing"
    owner   = "claude"
  }

  interfaces = [{
    type = "external"
  }]

  volumes = [{
    volume_id  = gcore_cloud_volume.tags_boot[0].id
    boot_index = 0
  }]
}

# ===========================================
# Servergroup for test 5
# ===========================================

resource "gcore_cloud_placement_group" "test" {
  count  = var.test_servergroup ? 1 : 0
  name   = "import-test-servergroup"
  policy = "soft-anti-affinity"
}

# ===========================================
# Test 5: Instance with servergroup
# ===========================================

resource "gcore_cloud_instance" "with_servergroup" {
  count = var.test_servergroup ? 1 : 0

  name           = "import-test-servergroup"
  flavor         = "g1-standard-1-2"
  project_id     = local.project_id
  region_id      = local.region_id
  servergroup_id = gcore_cloud_placement_group.test[0].servergroup_id

  interfaces = [{
    type = "external"
  }]

  volumes = [{
    volume_id  = gcore_cloud_volume.servergroup_boot[0].id
    boot_index = 0
  }]
}

# ===========================================
# Variables
# ===========================================

variable "test_minimal" {
  default = false
}

variable "test_two_volumes" {
  default = false
}

variable "test_floating_ip" {
  default = false
}

variable "test_tags" {
  default = false
}

variable "test_servergroup" {
  default = false
}

# ===========================================
# Outputs for import IDs
# ===========================================

output "minimal_id" {
  value = var.test_minimal ? gcore_cloud_instance.minimal[0].id : ""
}

output "two_volumes_id" {
  value = var.test_two_volumes ? gcore_cloud_instance.two_volumes[0].id : ""
}

output "floating_ip_id" {
  value = var.test_floating_ip ? gcore_cloud_instance.floating_ip[0].id : ""
}

output "with_tags_id" {
  value = var.test_tags ? gcore_cloud_instance.with_tags[0].id : ""
}

output "with_servergroup_id" {
  value = var.test_servergroup ? gcore_cloud_instance.with_servergroup[0].id : ""
}
