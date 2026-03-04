terraform {
  required_providers {
    gcore = {
      source  = "gcore/gcore"
    }
  }
}

provider "gcore" {}

# Variables for test configuration
variable "instance_name" {
  description = "Instance name"
  type        = string
  default     = "regression-test-instance"
}

variable "flavor" {
  description = "Flavor for the instance"
  type        = string
  default     = "g1-standard-1-2"
}

variable "boot_volume_size" {
  description = "Boot volume size in GB"
  type        = number
  default     = 5
}

variable "tags" {
  description = "Instance tags"
  type        = map(string)
  default     = {
    env     = "test"
    project = "regression"
  }
}

variable "servergroup_id" {
  description = "Servergroup ID to add instance to"
  type        = string
  default     = ""
}

variable "attach_new_fip" {
  description = "Whether to attach a new floating IP to private interface"
  type        = bool
  default     = false
}

variable "test_two_volumes" {
  description = "Whether to test with two volumes (for import drift test)"
  type        = bool
  default     = false
}

# ===========================================
# NETWORK INFRASTRUCTURE
# ===========================================

resource "gcore_cloud_network" "private" {
  name = "regression-test-network"
  type = "vlan"
}

resource "gcore_cloud_network_subnet" "private" {
  name       = "regression-test-subnet"
  cidr       = "192.168.100.0/24"
  network_id = gcore_cloud_network.private.id
  enable_dhcp = true
  dns_nameservers = ["8.8.8.8", "8.8.4.4"]
}

resource "gcore_cloud_network_router" "main" {
  name = "regression-test-router"
  external_gateway_info = {
    type        = "default"
    enable_snat = true
  }
  interfaces = [
    {
      subnet_id = gcore_cloud_network_subnet.private.id
      type      = "subnet"
    }
  ]
}

# ===========================================
# VOLUMES
# ===========================================

resource "gcore_cloud_volume" "boot" {
  name      = "regression-boot-volume"
  source    = "image"
  image_id  = "6343932d-0257-4285-bf89-05060f24095a"  # ubuntu-22.04-x64
  type_name = "ssd_hiiops"
  size      = var.boot_volume_size
}

# Second volume for two-volume import test (Bug 2)
resource "gcore_cloud_volume" "data" {
  count     = var.test_two_volumes ? 1 : 0
  name      = "regression-data-volume"
  source    = "new-volume"
  type_name = "ssd_hiiops"
  size      = 5
}

# ===========================================
# MAIN INSTANCE - Private interface with optional FIP
# ===========================================

resource "gcore_cloud_instance" "test" {
  # Only create when NOT testing two volumes (Bug 1 test uses this)
  count  = var.test_two_volumes ? 0 : 1
  name   = var.instance_name
  flavor = var.flavor

  # Single volume for Bug 1 test
  volumes = [
    {
      volume_id  = gcore_cloud_volume.boot.id
      boot_index = 0
    }
  ]

  # Private interface with optional floating IP (Bug 1 test)
  interfaces = [
    {
      type       = "subnet"
      network_id = gcore_cloud_network.private.id
      subnet_id  = gcore_cloud_network_subnet.private.id
      floating_ip = var.attach_new_fip ? {
        source = "new"
      } : null
    }
  ]

  tags = var.tags

  servergroup_id = var.servergroup_id != "" ? var.servergroup_id : null

  depends_on = [
    gcore_cloud_network_router.main
  ]
}

# ===========================================
# INSTANCE FOR IMPORT TEST (Bug 2)
# ===========================================

resource "gcore_cloud_instance" "test_import" {
  count  = var.test_two_volumes ? 1 : 0
  name   = "regression-import-test-instance"
  flavor = var.flavor

  # Two volumes: boot (boot_index=0) first, then data (boot_index=-1)
  volumes = [
    {
      volume_id  = gcore_cloud_volume.boot.id
      boot_index = 0
    },
    {
      volume_id  = gcore_cloud_volume.data[0].id
      boot_index = -1
    }
  ]

  interfaces = [
    {
      type = "external"
    }
  ]

  depends_on = [
    gcore_cloud_network_router.main
  ]
}

# ===========================================
# PLACEMENT GROUP
# ===========================================

resource "gcore_cloud_placement_group" "test" {
  name   = "regression-test-servergroup"
  policy = "soft-anti-affinity"
}

# ===========================================
# OUTPUTS
# ===========================================

output "instance_id" {
  value = var.test_two_volumes ? (length(gcore_cloud_instance.test_import) > 0 ? gcore_cloud_instance.test_import[0].id : null) : (length(gcore_cloud_instance.test) > 0 ? gcore_cloud_instance.test[0].id : null)
}

output "instance_status" {
  value = var.test_two_volumes ? (length(gcore_cloud_instance.test_import) > 0 ? gcore_cloud_instance.test_import[0].status : null) : (length(gcore_cloud_instance.test) > 0 ? gcore_cloud_instance.test[0].status : null)
}

output "instance_addresses" {
  value = var.test_two_volumes ? (length(gcore_cloud_instance.test_import) > 0 ? gcore_cloud_instance.test_import[0].addresses : null) : (length(gcore_cloud_instance.test) > 0 ? gcore_cloud_instance.test[0].addresses : null)
}

output "instance_interfaces" {
  value = var.test_two_volumes ? (length(gcore_cloud_instance.test_import) > 0 ? gcore_cloud_instance.test_import[0].interfaces : null) : (length(gcore_cloud_instance.test) > 0 ? gcore_cloud_instance.test[0].interfaces : null)
}

output "instance_volumes" {
  value = var.test_two_volumes ? (length(gcore_cloud_instance.test_import) > 0 ? gcore_cloud_instance.test_import[0].volumes : null) : (length(gcore_cloud_instance.test) > 0 ? gcore_cloud_instance.test[0].volumes : null)
}

output "router_id" {
  value = gcore_cloud_network_router.main.id
}

output "placement_group_id" {
  value = gcore_cloud_placement_group.test.servergroup_id
}

output "boot_volume_id" {
  value = gcore_cloud_volume.boot.id
}

output "subnet_id" {
  value = gcore_cloud_network_subnet.private.id
}
