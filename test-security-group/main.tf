terraform {
  required_providers {
    gcore = {
      source  = "gcore/gcore"
      version = "~> 1.0"
    }
  }
}

provider "gcore" {
  # Reads from environment variables:
  # - GCORE_API_KEY
  # - GCORE_CLOUD_PROJECT_ID
  # - GCORE_CLOUD_REGION_ID
}

# Test security group - using nested pattern (security_group is now optional)
resource "gcore_cloud_security_group" "test" {
  project_id = 379987
  region_id  = 76

  security_group = {
    name        = "test-sg-after-rebase"
    description = "Test security group after rebase - verifying no drift"
    security_group_rules = [
      {
        direction        = "egress"
        ethertype        = "IPv4"
        protocol         = "tcp"
        port_range_min   = 443
        port_range_max   = 443
        remote_ip_prefix = "0.0.0.0/0"
      }
    ]
  }
}

output "security_group_id" {
  value = gcore_cloud_security_group.test.id
}

output "security_group_name" {
  value = gcore_cloud_security_group.test.name
}
