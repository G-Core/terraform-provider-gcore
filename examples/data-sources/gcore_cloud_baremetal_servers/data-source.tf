data "gcore_cloud_baremetal_servers" "example_cloud_baremetal_servers" {
  project_id = 1
  region_id = 1
  changes_before = "2025-10-01T12:00:00Z"
  changes_since = "2025-10-01T12:00:00Z"
  flavor_id = "bm2-hf-small"
  flavor_prefix = "bm2-"
  ip = "192.168.0.1"
  name = "name"
  only_with_fixed_external_ip = true
  profile_name = "profile_name"
  protection_status = "Active"
  status = "ACTIVE"
  tag_key_value = "tag_key_value"
  tag_value = ["value1", "value2"]
  uuid = "b5b4d65d-945f-4b98-ab6f-332319c724ef"
}
