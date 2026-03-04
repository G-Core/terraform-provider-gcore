terraform {
  required_providers {
    gcore = {
      source  = "gcore/gcore"
      version = "~> 1.0"
    }
  }
}

provider "gcore" {
  # Credentials loaded from environment:
  # - GCORE_API_KEY
  # - GCORE_CLOUD_PROJECT_ID
  # - GCORE_CLOUD_REGION_ID
}

# Basic instance with simple network interface
# Using a minimal configuration that should work
resource "gcore_cloud_instance" "test" {
  name       = "test-instance-andpoll"
  flavor     = "g1-standard-1-2"
  project_id = 379987
  region_id  = 76

  interfaces = [
    {
      type = "external"
    }
  ]

  volumes = [
    {
      source     = "image"
      image_id   = "f4ce3d30-e29c-4cfd-811f-46f383b6081f"  # Placeholder - update with valid ID
      size       = 10
      boot_index = 0
    }
  ]

  tags = {
    environment = "test"
    managed_by  = "terraform"
    test_type   = "andpoll_migration"
  }
}

# Output instance details
output "instance_id" {
  value = gcore_cloud_instance.test.id
}

output "instance_name" {
  value = gcore_cloud_instance.test.name
}

output "instance_status" {
  value = gcore_cloud_instance.test.vm_state
}

output "instance_addresses" {
  value = gcore_cloud_instance.test.addresses
}
