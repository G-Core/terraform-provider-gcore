data "gcore_cloud_tasks" "example_cloud_tasks" {
  from_timestamp = "2019-12-27T18:11:19.117Z"
  is_acknowledged = true
  project_id = [0, 0]
  region_id = [0, 0]
  state = ["ERROR", "FINISHED"]
  task_type = "task_type"
  to_timestamp = "2019-12-27T18:11:19.117Z"
}
