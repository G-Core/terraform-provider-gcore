# TC-05: Multiple rules in a loop (for_each pattern)
# Expected: All rules created successfully, no interference
variable "ports" {
  type = map(object({
    port        = number
    protocol    = string
    description = string
  }))
  default = {
    http = {
      port        = 80
      protocol    = "tcp"
      description = "HTTP traffic"
    }
    mysql = {
      port        = 3306
      protocol    = "tcp"
      description = "MySQL database"
    }
    redis = {
      port        = 6379
      protocol    = "tcp"
      description = "Redis cache"
    }
  }
}

resource "gcore_cloud_security_group_rule" "ports" {
  for_each = var.ports

  group_id   = gcore_cloud_security_group.test.id
  project_id = 379987
  region_id  = 76

  direction        = "ingress"
  ethertype        = "IPv4"
  protocol         = each.value.protocol
  port_range_min   = each.value.port
  port_range_max   = each.value.port
  remote_ip_prefix = "0.0.0.0/0"
  description      = each.value.description
}

output "loop_rule_ids" {
  value = { for k, r in gcore_cloud_security_group_rule.ports : k => r.id }
}
