terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {
  # Credentials from .env
}

resource "gcore_cloud_floating_ip" "test_fip" {
  project_id = 379987
  region_id  = 76
  # Don't specify port_id - let it auto-assign
}

output "fip_id" {
  value = gcore_cloud_floating_ip.test_fip.id
}

output "fip_status" {
  value = gcore_cloud_floating_ip.test_fip.status  
}

output "fip_created_at" {
  value = gcore_cloud_floating_ip.test_fip.created_at
}

output "fip_floating_ip_address" {
  value = gcore_cloud_floating_ip.test_fip.floating_ip_address
}
