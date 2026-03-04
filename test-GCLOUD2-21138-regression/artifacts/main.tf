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

# ===========================================
# MAIN INSTANCE - Simple external interface
# ===========================================

resource "gcore_cloud_instance" "test" {
  name   = var.instance_name
  flavor = var.flavor

  volumes = [
    {
      volume_id  = gcore_cloud_volume.boot.id
      boot_index = 0
    }
  ]

  interfaces = [
    {
      type = "external"
    }
  ]

  tags = var.tags

  servergroup_id = var.servergroup_id != "" ? var.servergroup_id : null

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
  value = gcore_cloud_instance.test.id
}

output "instance_status" {
  value = gcore_cloud_instance.test.status
}

output "instance_addresses" {
  value = gcore_cloud_instance.test.addresses
}

output "instance_interfaces" {
  value = gcore_cloud_instance.test.interfaces
}

output "instance_volumes" {
  value = gcore_cloud_instance.test.volumes
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
