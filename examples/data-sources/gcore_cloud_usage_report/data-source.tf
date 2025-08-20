data "gcore_cloud_usage_report" "example_cloud_usage_report" {
  time_from = "2023-01-01T00:00:00Z"
  time_to = "2023-02-01T00:00:00Z"
  enable_last_day = false
  limit = 10
  offset = 0
  projects = [16, 17, 18, 19, 20]
  regions = [1, 2, 3]
  schema_filter = {
    field = "flavor"
    type = "instance"
    values = ["g1-standard-1-2"]
  }
  sorting = [{
    billing_value = "asc"
    first_seen = "asc"
    last_name = "asc"
    last_seen = "asc"
    project = "asc"
    region = "asc"
    type = "asc"
  }]
  tags = {
    conditions = [{
      key = "os_version"
      strict = true
      value = "22.04"
    }, {
      key = "os_version"
      strict = true
      value = "23.04"
    }]
    condition_type = "OR"
  }
  types = ["egress_traffic", "instance"]
}
