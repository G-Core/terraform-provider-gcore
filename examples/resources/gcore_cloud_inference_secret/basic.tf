# Create an AWS IAM secret for use with inference SQS triggers
resource "gcore_cloud_inference_secret" "example" {
  project_id     = 1
  name           = "my-aws-iam-secret"
  type           = "aws-iam"
  data_wo_version = 1
  data = {
    aws_access_key_id_wo     = "my-aws-access-key-id"
    aws_secret_access_key_wo = "my-aws-secret-key"
  }
}
