# Block a single IP address for 12 hours with a 429 response
data "gcore_waap_domain" "example" {
  find_one_by = {
    name = "example.com"
  }
}

resource "gcore_waap_domain_firewall_rule" "block_ip" {
  domain_id   = data.gcore_waap_domain.example.id
  name        = "Block abusive client"
  description = "Rate-limit offender"
  enabled     = true

  action = {
    block = {
      status_code     = 429
      action_duration = "12h"
    }
  }

  conditions = [{
    ip = {
      ip_address = "203.0.113.10"
    }
  }]
}
