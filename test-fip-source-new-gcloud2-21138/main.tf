terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

# Known working Ubuntu 22.04 image ID for region 76
locals {
  ubuntu_image_id = "6343932d-0257-4285-bf89-05060f24095a"
  # Using existing network/subnet from regression test to avoid quota limits
  network_id = "ca86be86-793c-450b-9dc6-429b61cde52f"
  subnet_id  = "bc24a87d-d24a-4e6f-acac-12bb04c74d09"
}

# Boot volume - must be created separately in new provider
resource "gcore_cloud_volume" "boot" {
  name       = "test-fip-boot-volume-2"
  source     = "image"
  image_id   = local.ubuntu_image_id
  type_name  = "standard"
  size       = 10
  region_id  = 76
  project_id = 379987
}

# Test instance with floating_ip source="new" on PRIVATE interface
# GCLOUD2-21138: Tests that updates don't attempt duplicate FIP creation
resource "gcore_cloud_instance" "test" {
  name       = var.instance_name
  flavor     = var.flavor
  region_id  = 76
  project_id = 379987

  # Attach the boot volume
  volumes = [{
    volume_id  = gcore_cloud_volume.boot.id
    boot_index = 0
  }]

  # Private interface (subnet type) with source="new" floating IP
  # This is the scenario from GCLOUD2-21138 bug
  interfaces = [{
    type       = "subnet"
    network_id = local.network_id
    subnet_id  = local.subnet_id

    # This is the key test: source="new" should create FIP once,
    # and NOT attempt to create again on subsequent updates
    floating_ip = {
      source = "new"
    }
  }]
}

variable "instance_name" {
  description = "Instance name - change this to test updates"
  default     = "test-fip-source-new"
}

variable "flavor" {
  description = "Instance flavor - change this to test resize"
  default     = "g1-standard-1-2"
}

output "instance_id" {
  value = gcore_cloud_instance.test.id
}

output "instance_name" {
  value = gcore_cloud_instance.test.name
}

output "interfaces" {
  value = gcore_cloud_instance.test.interfaces
}
