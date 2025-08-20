data "gcore_waap_domain_custom_rules" "example_waap_domain_custom_rules" {
  domain_id = 1
  action = "block"
  description = "This rule blocks all the requests coming form a specific IP address."
  enabled = false
  name = "Block by specific IP rule."
  ordering = "-id"
}
