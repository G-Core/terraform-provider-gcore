provider gcore {
  permanent_api_token = "251$d3361.............1b35f26d8"
}

data "gcore_waap_security_insight_type" "attack_on_disabled_policy" {
  name = "Attack on disabled policy"
}

output "insight_type" {
  value = data.gcore_waap_security_insight_type.attack_on_disabled_policy.id
}
