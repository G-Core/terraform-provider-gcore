data "gcore_cloud_inference_deployment_logs" "example_cloud_inference_deployment_logs" {
  project_id = 1
  deployment_name = "my-instance"
  region_id = 1
}
