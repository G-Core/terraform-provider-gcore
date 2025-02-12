resource "gcore_inference_secret" "aws" {
  project_id = data.gcore_project.project.id
  name = "my-aws-iam-secret"
  data_aws_access_key_id = "my-aws-access-key-id"
  data_aws_secret_access_key = "my-aws-access-key"
}

resource "gcore_inference_deployment" "inf" {
  project_id = data.gcore_project.project.id
  name = "my-inference-deployment"
  image = "nginx:latest"
  listening_port = 80
  flavor_name = "inference-4vcpu-16gib"
  timeout = 60
  containers {
    region_id  = data.gcore_region.region.id
    cooldown_period = 60
    polling_interval = 60
    scale_min = 0
    scale_max = 2
    triggers_cpu_threshold = 80

    triggers_sqs_secret_name = gcore_inference_secret.aws.name
    triggers_sqs_aws_region = "us-west-2"
    triggers_sqs_queue_url = "https://sqs.us-west-2.amazonaws.com/1234567890/my-queue"
    triggers_sqs_queue_length = 5
    triggers_sqs_activation_queue_length = 2
  }

  liveness_probe {
    enabled = false
  }

  readiness_probe {
    enabled = false
  }

  startup_probe {
    enabled = false
  }
}