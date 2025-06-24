provider gcore {
  permanent_api_token = "768660$.............a43f91f"
}

data "gcore_waap_security_insight_type" "attack_on_disabled_policy" {
  name = "Attack on disabled policy"
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

resource "gcore_waap_security_insight_silence" "insight_silence" {
  domain_id = gcore_waap_domain.domain.id
  insight_type = data.gcore_waap_security_insight_type.attack_on_disabled_policy.id
  comment = "Example Insight Silence"
  author = "Gcore"
  labels = {
    attack_type = "exmple"
  }
  expire_at = "2026-08-24T14:00:00Z"
}
