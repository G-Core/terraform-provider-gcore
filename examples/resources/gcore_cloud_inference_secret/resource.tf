resource "gcore_cloud_inference_secret" "example_cloud_inference_secret" {
  project_id = 1
  data = {
    aws_access_key_id = "fake-key-id"
    aws_secret_access_key = "fake-secret"
  }
  name = "aws-dev"
  type = "aws-iam"
}
