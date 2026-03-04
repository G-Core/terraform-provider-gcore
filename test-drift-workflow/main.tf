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

# Test AWS-style drift detection
resource "gcore_cloud_security_group" "drift_test" {
  project_id = 379987
  region_id  = 76

  security_group = {
    name        = "drift-test-workflow"
    description = "Testing AWS-style drift detection"
  }

  # NO security_group_rules specified - user wants ZERO inline rules
  # Backend-created rules will show as drift
}

output "security_group_id" {
  value = gcore_cloud_security_group.drift_test.id
}

output "rules_count" {
  value = length(gcore_cloud_security_group.drift_test.security_group_rules)
}
