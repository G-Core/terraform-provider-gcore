# Allow an office IP range. Use `allow = {}`: the allow action has no settings.
data "gcore_waap_domain" "example" {
  find_one_by = {
    name = "example.com"
  }
}

resource "gcore_waap_domain_firewall_rule" "allow_office" {
  domain_id = data.gcore_waap_domain.example.id
  name      = "Allow office network"
  enabled   = true

  action = {
    allow = {}
  }

  conditions = [{
    ip_range = {
      lower_bound = "198.51.100.1"
      upper_bound = "198.51.100.254"
    }
  }]
}
