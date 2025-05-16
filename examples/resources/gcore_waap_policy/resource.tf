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
  name   = gcore_cdn_resource.example.cname
}

data "gcore_waap_policy" "policy" {
  domain_id = gcore_waap_domain.domain.id
}

resource "gcore_waap_policy" "policy_resource" {
  domain_id = gcore_waap_domain.domain_resource.id

  policies = [
    {
      policy_name = "protocol_validation"
      rules = {
        invalid_user_agent = true
        unknown_user_agent = true
      }
    },
    {
      policy_name = "core_waf_owasp_top_threats"
      rules = {
        sql_injection = false
        xss           = false
      }
    },
    {
      policy_name = "anti_automation_bot_protection"
      rules = {
        anti_scraping     = true
        automated_clients = true
      }
    }
  ]
}