resource "gcore_cloud_load_balancer_l7_policy_rule" "example_cloud_load_balancer_l7_policy_rule" {
  project_id = 0
  region_id = 0
  l7policy_id = "l7policy_id"
  compare_type = "CONTAINS"
  type = "COOKIE"
  value = "value"
  invert = true
  key = "key"
  tags = ["string"]
}
