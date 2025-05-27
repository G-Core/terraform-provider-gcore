provider gcore {
  permanent_api_token = "251$d3361.............1b35f26d8"
}

data "gcore_waap_domain_policy" "waf_sql_injection" {
  domain_id = 12345
  name = "SQL injection"
  group = "WAF"
}

output "waf_sql_injection_id" {
  value = data.gcore_waap_domain_policy.waf_sql_injection.id
}
