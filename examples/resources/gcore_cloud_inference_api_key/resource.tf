resource "gcore_cloud_inference_api_key" "example_cloud_inference_api_key" {
  project_id = 1
  name = "my-api-key"
  description = "This key is used for accessing the inference service."
  expires_at = "2024-10-01T12:00:00Z"
}
