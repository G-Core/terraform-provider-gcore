data "gcore_cloud_audit_logs" "example_cloud_audit_logs" {
  action_type = ["activate", "delete"]
  api_group = ["ai_cluster", "image"]
  from_timestamp = "2019-11-14T10:30:32Z"
  project_id = [1, 2, 3]
  region_id = [1, 2, 3]
  resource_id = ["string"]
  search_field = "default"
  to_timestamp = "2019-11-14T10:30:32Z"
}
