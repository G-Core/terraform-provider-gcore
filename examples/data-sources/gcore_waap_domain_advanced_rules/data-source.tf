data "gcore_waap_domain_advanced_rules" "example_waap_domain_advanced_rules" {
  domain_id = 1
  action = "allow"
  description = "This rule blocks all the requests coming form a specific IP address"
  enabled = false
  name = "Block by specific IP rule"
  ordering = "id"
  phase = "access"
}
