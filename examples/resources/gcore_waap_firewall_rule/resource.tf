provider gcore {
  permanent_api_token = "768660$.............a43f91f"
}

resource "gcore_cdn_resource" "example" {
  cname  = "api.example.com"
  origin = "origin.example.com"
  options {
    waap { value = true }
  }
}

resource "gcore_waap_domain" "domain" {
  name   = gcore_cdn_resource.cdn_resource.cname
  status = "monitor"
}

resource "gcore_waap_firewall_rule" "firewall_rule" {
  domain_id   = gcore_waap_domain.domain.id
  name        = "Block IP Range"
  description = "Block range of IP addresses"
  enabled     = true

  action {
    block {
      status_code = 403
      action_duration = "12h"
    }
  }

  conditions {
    ip_range {
      lower_bound = "192.168.1.1"
      upper_bound = "192.168.1.7"
    }
  }
}