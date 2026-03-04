terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

variable "api_key" {
  type        = string
  description = "Gcore API key"
  sensitive   = true
}

provider "gcore" {
  api_key = var.api_key
}

data "gcore_cloud_projects" "my_projects" {
  name = "default"
}

locals {
  project_id = [for p in data.gcore_cloud_projects.my_projects.items : p.id]
}

data "gcore_cloud_region" "rg" {
  region_id = 76
}

resource "gcore_cloud_load_balancer" "lb" {
  project_id = local.project_id[0]
  region_id  = data.gcore_cloud_region.rg.id
  flavor = "lb1-2-4"
  name = "qa-tf-lb"
  tags = {
    "qa" = "load-balancer"
  }
}

resource "gcore_cloud_load_balancer_listener" "ls" {
  project_id = local.project_id[0]
  region_id  = data.gcore_cloud_region.rg.id
  load_balancer_id = gcore_cloud_load_balancer.lb.id
  name = "qa-ls-name"
  protocol = "HTTP"
  protocol_port = 80
  connection_limit = 5000
  timeout_client_data = 45000
  timeout_member_connect = 1000
}

resource "gcore_cloud_load_balancer_pool" "lb_pool" {
  project_id = local.project_id[0]
  region_id  = data.gcore_cloud_region.rg.id
  lb_algorithm = "LEAST_CONNECTIONS"
  name = "pool-3"
  protocol = "HTTP"
  listener_id = gcore_cloud_load_balancer_listener.ls.id
  timeout_client_data = 50000
  timeout_member_connect = 50000
  timeout_member_data = 0
  healthmonitor = {
    delay = 10
    max_retries = 3
    timeout = 5
    type = "HTTP"
    expected_codes = "200,301,302"
    http_method = "GET"
    max_retries_down = 3
    url_path = "/"
  }
  # Testing fix: commenting out entire members attribute
  # members = [
  #   {
  #     address = "10.0.0.1"
  #     protocol_port = 8080
  #   }
  # ]
}

output "pool_id" {
  value = gcore_cloud_load_balancer_pool.lb_pool.id
}

output "pool_members" {
  value = gcore_cloud_load_balancer_pool.lb_pool.members
}
