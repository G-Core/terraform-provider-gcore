resource "gcore_cloud_inference_deployment" "example_cloud_inference_deployment" {
  project_id = 1
  containers = [{
    region_id = 1
    scale = {
      max = 3
      min = 1
      cooldown_period = 60
      polling_interval = 30
      triggers = {
        cpu = {
          threshold = 80
        }
        gpu_memory = {
          threshold = 80
        }
        gpu_utilization = {
          threshold = 80
        }
        http = {
          rate = 1
          window = 60
        }
        memory = {
          threshold = 70
        }
        sqs = {
          activation_queue_length = 1
          aws_region = "us-east-1"
          queue_length = 10
          queue_url = "https://sqs.us-east-1.amazonaws.com/123456789012/MyQueue"
          secret_name = "x"
          aws_endpoint = "aws_endpoint"
          scale_on_delayed = true
          scale_on_flight = true
        }
      }
    }
  }]
  flavor_name = "inference-16vcpu-232gib-1xh100-80gb"
  image = "nginx:latest"
  listening_port = 80
  name = "my-instance"
  api_keys = ["key1", "key2"]
  auth_enabled = false
  command = ["nginx", "-g", "daemon off;"]
  credentials_name = "dockerhub"
  description = "My first instance"
  envs = {
    DEBUG_MODE = "False"
    KEY = "12345"
  }
  ingress_opts = {
    disable_response_buffering = true
  }
  logging = {
    destination_region_id = 1
    enabled = true
    retention_policy = {
      period = 42
    }
    topic_name = "my-log-name"
  }
  probes = {
    liveness_probe = {
      enabled = true
      probe = {
        exec = {
          command = ["ls", "-l"]
        }
        failure_threshold = 3
        http_get = {
          port = 80
          headers = {
            Authorization = "Bearer token 123"
          }
          host = "127.0.0.1"
          path = "/healthz"
          schema = "HTTP"
        }
        initial_delay_seconds = 0
        period_seconds = 5
        success_threshold = 1
        tcp_socket = {
          port = 80
        }
        timeout_seconds = 1
      }
    }
    readiness_probe = {
      enabled = true
      probe = {
        exec = {
          command = ["ls", "-l"]
        }
        failure_threshold = 3
        http_get = {
          port = 80
          headers = {
            Authorization = "Bearer token 123"
          }
          host = "127.0.0.1"
          path = "/healthz"
          schema = "HTTP"
        }
        initial_delay_seconds = 0
        period_seconds = 5
        success_threshold = 1
        tcp_socket = {
          port = 80
        }
        timeout_seconds = 1
      }
    }
    startup_probe = {
      enabled = true
      probe = {
        exec = {
          command = ["ls", "-l"]
        }
        failure_threshold = 3
        http_get = {
          port = 80
          headers = {
            Authorization = "Bearer token 123"
          }
          host = "127.0.0.1"
          path = "/healthz"
          schema = "HTTP"
        }
        initial_delay_seconds = 0
        period_seconds = 5
        success_threshold = 1
        tcp_socket = {
          port = 80
        }
        timeout_seconds = 1
      }
    }
  }
  timeout = 120
}
