terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

variable "resource_id" {
  description = "CDN resource ID to create the rule for"
  type        = number
}

resource "gcore_cdn_cdn_resource_rule" "test" {
  resource_id = var.resource_id
  name        = "tf-test-rule-jira-23533"
  rule        = "/images"
  rule_type   = 0
  active      = true
}

output "rule_id" {
  value = gcore_cdn_cdn_resource_rule.test.id
}
output "rule_name" {
  value = gcore_cdn_cdn_resource_rule.test.name
}
output "rule_weight" {
  value = gcore_cdn_cdn_resource_rule.test.weight
}
output "origin_protocol" {
  value = gcore_cdn_cdn_resource_rule.test.origin_protocol
}
