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

output "policy_ids" {
  value = data.gcore_waap_policy.policy_data.policies
}
