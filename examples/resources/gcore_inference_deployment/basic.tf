resource "gcore_inference_deployment" "inf" {
  project_id = data.gcore_project.project.id
  name = "my-inference-deployment"
  image = "nginx:latest"
  listening_port = 80
  flavor_name = "inference-4vcpu-16gib"
  containers {
    region_id  = data.gcore_region.region.id
    cooldown_period = 30
    scale_min = 2
    scale_max = 2
    triggers_cpu_threshold = 80
  }

  # If you don't specify any probe, the container may be marked as "ready" too soon,
  # meaning it will start accepting requests before your application has fully initialized.
  # This can lead to errors, as the app might not be prepared to handle incoming traffic yet.
  liveness_probe {
    enabled = true
    failure_threshold = 3
    initial_delay_seconds = 10
    period_seconds = 10
    timeout_seconds = 1
    success_threshold = 1
    http_get_port = 80
    http_get_headers = {
      User-Agent = "my user agent"
    }
    http_get_host = "localhost"
    http_get_path = "/"
    http_get_schema = "HTTPS"
  }

  readiness_probe {
    enabled = false
  }

  startup_probe {
    enabled = false
  }
}