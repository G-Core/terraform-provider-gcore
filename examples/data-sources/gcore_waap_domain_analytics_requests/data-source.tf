data "gcore_waap_domain_analytics_requests" "example_waap_domain_analytics_requests" {
  domain_id = 1
  start = "2019-12-27T18:11:19.117Z"
  actions = ["allow"]
  countries = ["Mv"]
  end = "2019-12-27T18:11:19.117Z"
  ip = ".:"
  ordering = "ordering"
  reference_id = "2c02efDd09B3BA1AEaDd3dCAa7aC7A37"
  security_rule_name = "security_rule_name"
  status_code = 100
  traffic_types = ["policy_allowed"]
}
