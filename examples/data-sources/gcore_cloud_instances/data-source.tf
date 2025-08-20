data "gcore_cloud_instances" "example_cloud_instances" {
  project_id = 1
  region_id = 1
  available_floating = true
  changes_before = "2025-10-01T12:00:00Z"
  changes_since = "2025-10-01T12:00:00Z"
  exclude_flavor_prefix = "g1-"
  exclude_secgroup = "secgroup_name"
  flavor_id = "g2-standard-32-64"
  flavor_prefix = "g2-"
  ip = "192.168.0.1"
  name = "name"
  only_with_fixed_external_ip = true
  profile_name = "profile_name"
  protection_status = "Active"
  status = "ACTIVE"
  tag_key_value = "tag_key_value"
  tag_value = ["value1", "value2"]
  type_ddos_profile = "advanced"
  uuid = "b5b4d65d-945f-4b98-ab6f-332319c724ef"
}
