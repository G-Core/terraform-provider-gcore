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
  default     = 949318
}

variable "rule_name" {
  description = "Name of the CDN rule"
  type        = string
  default     = "tf-test-rule-23533"
}

variable "rule_path" {
  description = "Rule pattern"
  type        = string
  default     = "/images"
}

variable "active" {
  description = "Whether the rule is active"
  type        = bool
  default     = true
}

resource "gcore_cdn_cdn_resource_rule" "test" {
  resource_id = var.resource_id
  name        = var.rule_name
  rule        = var.rule_path
  rule_type   = 0
  active      = var.active
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
output "rule_active" {
  value = gcore_cdn_cdn_resource_rule.test.active
}
output "origin_protocol" {
  value = gcore_cdn_cdn_resource_rule.test.origin_protocol
}
output "deleted" {
  value = gcore_cdn_cdn_resource_rule.test.deleted
}
output "preset_applied" {
  value = gcore_cdn_cdn_resource_rule.test.preset_applied
}
