data "gcore_waap_domain_firewall_rules" "example_waap_domain_firewall_rules" {
  domain_id = 1
  action = "allow"
  description = "description"
  enabled = true
  name = "name"
  ordering = "-id"
}
