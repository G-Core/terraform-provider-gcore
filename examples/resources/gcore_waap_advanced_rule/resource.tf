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

resource "gcore_waap_advanced_rule" "advanced_rule" {
  domain_id = gcore_waap_domain.domain.id
  name = "Advanced Rule"
  enabled = true
  action {
    block {
      status_code = 403
      action_duration = "5m"
    }
  }
  source = "request.ip == '117.20.32.55'"
  description = "Description of the advanced rule"
  phase = "access"
}