terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {
  # Uses environment variables
}

# Security group WITHOUT any inline rules - testing AWS-style drift detection
resource "gcore_cloud_security_group" "test" {
  project_id = 379987
  region_id  = 76

  security_group = {
    name        = "test-rules-only"
    description = "Test: rules only via separate resources"
    # NO security_group_rules field specified - user wants ZERO inline rules
  }
}

output "security_group_id" {
  value = gcore_cloud_security_group.test.id
}
