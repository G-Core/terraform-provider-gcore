provider gcore {
  permanent_api_token = "768660$.............a43f91f"
}

resource "gcore_cdn_resource" "cdn_resource" {
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

data "gcore_waap_domain_policy" "invalid_user_agent" {
  domain_id = gcore_waap_domain.domain.id
  name = "Invalid user agent"
  group = "PROTOCOL"
}

data "gcore_waap_domain_policy" "waf_sql_injection" {
  domain_id = gcore_waap_domain.domain.id
  name = "SQL injection"
  group = "WAF"
}

resource "gcore_waap_policy" "domain_policy" {
  domain_id = gcore_waap_domain.domain.id

  policies = {
    # you can use gcore_waap_domain_policy data source to specify the policy id
    "${data.gcore_waap_domain_policy.invalid_user_agent.id}" = false
    "${data.gcore_waap_domain_policy.waf_sql_injection.id}" = true

    # or you can use the policy id directly
    S55075105 = true  # Traffic via TOR network
    S55075106 = true  # Traffic via proxy networks
    S55075107 = true  # Traffic from hosting services
  }
}