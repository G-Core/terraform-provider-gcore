provider gcore {
  permanent_api_token = "251$d3361.............1b35f26d8"
}

data "gcore_project" "project" {
  name = "Default"
}

resource "gcore_inference_secret" "aws" {
  project_id = data.gcore_project.project.id
  name = "my-aws-iam-secret"
  data_aws_access_key_id = "my-aws-access-key-id"
  data_aws_secret_access_key = "my-aws-access-key"
}