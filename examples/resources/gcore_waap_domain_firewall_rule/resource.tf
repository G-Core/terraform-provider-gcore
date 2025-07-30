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
      ip_address = "192.168.1.1"
      negation = true
    }
    ip_range = {
      lower_bound = "192.168.1.1"
      upper_bound = "192.168.1.1"
      negation = true
    }
  }]
  enabled = true
  name = "name"
  description = "description"
}
