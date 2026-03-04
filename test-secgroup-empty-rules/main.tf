terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

resource "gcore_cloud_security_group" "test" {
  name        = "tf-test-empty-rules"
  description = "Test that no default rules are created"
}

output "sg_id" {
  value = gcore_cloud_security_group.test.id
}

output "sg_rules" {
  value = gcore_cloud_security_group.test.security_group_rules
}
