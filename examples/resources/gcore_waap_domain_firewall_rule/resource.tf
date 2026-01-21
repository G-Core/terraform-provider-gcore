resource "gcore_waap_domain_firewall_rule" "example_waap_domain_firewall_rule" {
  domain_id = 1
  action = {
    allow = {

    }
    block = {
      action_duration = "12h"
      status_code = 403
    }
  }
  conditions = [{
    ip = {
      ip_address = "ip_address"
      negation = true
    }
    ip_range = {
      lower_bound = "lower_bound"
      upper_bound = "upper_bound"
      negation = true
    }
  }]
  enabled = true
  name = "Block foobar bot"
  description = "description"
}
