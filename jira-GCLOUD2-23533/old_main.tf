terraform {
  required_providers {
    gcore = {
      source  = "G-Core/gcore"
      version = "~> 0.10"
    }
  }
}

provider "gcore" {
  permanent_api_token = var.api_token
}

variable "api_token" {
  type      = string
  sensitive = true
}

variable "resource_id" {
  description = "CDN resource ID to create the rule for"
  type        = number
}

resource "gcore_cdn_rule" "test" {
  resource_id = var.resource_id
  name        = "tf-test-rule-jira-23533"
  rule        = "/images"
  rule_type   = 0
  active      = true
}

output "rule_id" {
  value = gcore_cdn_rule.test.id
}
output "rule_name" {
  value = gcore_cdn_rule.test.name
}
output "rule_weight" {
  value = gcore_cdn_rule.test.weight
}
