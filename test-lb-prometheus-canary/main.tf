terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

data "gcore_cloud_projects" "selected" {
  name = var.project_name
}

locals {
  project_id = data.gcore_cloud_projects.selected.items[0].id
}

data "gcore_cloud_region" "target" {
  region_id = var.region_id
}

resource "gcore_cloud_network" "vip" {
  name       = "${var.lb_name}-vip"
  project_id = local.project_id
  region_id  = data.gcore_cloud_region.target.id
}

resource "gcore_cloud_network_subnet" "vip" {
  name       = "${var.lb_name}-vip"
  cidr       = var.vip_network_cidr
  network_id = gcore_cloud_network.vip.id
  project_id = local.project_id
  region_id  = data.gcore_cloud_region.target.id
}

resource "gcore_cloud_network" "backend" {
  name       = "${var.lb_name}-backend"
  project_id = local.project_id
  region_id  = data.gcore_cloud_region.target.id
}

resource "gcore_cloud_network_subnet" "backend" {
  name       = "${var.lb_name}-backend"
  cidr       = var.backend_subnet_cidr
  network_id = gcore_cloud_network.backend.id
  project_id = local.project_id
  region_id  = data.gcore_cloud_region.target.id
}

resource "gcore_cloud_load_balancer" "canary" {
  project_id = local.project_id
  region_id  = data.gcore_cloud_region.target.id
  flavor     = var.lb_flavor
  name       = var.lb_name

  vip_network_id = gcore_cloud_network.vip.id
  vip_subnet_id  = gcore_cloud_network_subnet.vip.id

  tags = var.lb_tags
}

resource "gcore_cloud_load_balancer_listener" "http_primary" {
  project_id       = local.project_id
  region_id        = data.gcore_cloud_region.target.id
  load_balancer_id = gcore_cloud_load_balancer.canary.id
  name             = "${var.lb_name}-prod"
  protocol         = "HTTP"
  protocol_port    = 80

  allowed_cidrs        = var.allowed_cidrs
  connection_limit     = 80000
  timeout_client_data  = 60000
  timeout_member_connect = 8000
  timeout_member_data  = 45000
}

resource "gcore_cloud_load_balancer_listener" "http_canary" {
  project_id       = local.project_id
  region_id        = data.gcore_cloud_region.target.id
  load_balancer_id = gcore_cloud_load_balancer.canary.id
  name             = "${var.lb_name}-canary"
  protocol         = "HTTP"
  protocol_port    = var.canary_port

  connection_limit     = 5000
  timeout_client_data  = 15000
  timeout_member_connect = 2000
  timeout_member_data  = 10000
}

resource "gcore_cloud_load_balancer_listener" "prometheus" {
  project_id       = local.project_id
  region_id        = data.gcore_cloud_region.target.id
  load_balancer_id = gcore_cloud_load_balancer.canary.id
  name             = "${var.lb_name}-prometheus"
  protocol         = "PROMETHEUS"
  protocol_port    = var.prometheus_port

  timeout_client_data    = 15000
  timeout_member_connect = 3000
  timeout_member_data    = 5000

  secret_id = var.prometheus_secret_id
  user_list = var.prometheus_users
}

resource "gcore_cloud_load_balancer_pool" "prod" {
  project_id       = local.project_id
  region_id        = data.gcore_cloud_region.target.id
  listener_id      = gcore_cloud_load_balancer_listener.http_primary.id
  load_balancer_id = gcore_cloud_load_balancer.canary.id
  name             = "${var.lb_name}-prod"
  protocol         = "HTTP"
  lb_algorithm     = "ROUND_ROBIN"

  healthmonitor = {
    type           = "HTTP"
    delay          = 10
    timeout        = 5
    max_retries    = 3
    max_retries_down = 2
    url_path       = "/ready"
    http_method    = "GET"
    expected_codes = "200-299"
  }

  timeout_client_data    = 60000
  timeout_member_connect = 8000
  timeout_member_data    = 45000
}

resource "gcore_cloud_load_balancer_pool" "canary" {
  project_id       = local.project_id
  region_id        = data.gcore_cloud_region.target.id
  listener_id      = gcore_cloud_load_balancer_listener.http_canary.id
  load_balancer_id = gcore_cloud_load_balancer.canary.id
  name             = "${var.lb_name}-canary"
  protocol         = "HTTP"
  lb_algorithm     = "LEAST_CONNECTIONS"

  session_persistence = {
    type        = "APP_COOKIE"
    cookie_name = "canary-cookie"
  }

  healthmonitor = {
    type       = "HTTP"
    delay      = 5
    timeout    = 3
    max_retries = 2
    url_path   = var.canary_health_path
  }

  timeout_client_data    = 20000
  timeout_member_connect = 4000
  timeout_member_data    = 15000
}

resource "gcore_cloud_load_balancer_pool" "metrics" {
  project_id       = local.project_id
  region_id        = data.gcore_cloud_region.target.id
  listener_id      = gcore_cloud_load_balancer_listener.prometheus.id
  load_balancer_id = gcore_cloud_load_balancer.canary.id
  name             = "${var.lb_name}-metrics"
  protocol         = "HTTP"
  lb_algorithm     = "ROUND_ROBIN"

  timeout_client_data    = 15000
  timeout_member_connect = 3000
  timeout_member_data    = 5000

  healthmonitor = {
    type        = "HTTP"
    delay       = 10
    timeout     = 3
    max_retries = 2
    url_path    = "/metrics"
  }
}

resource "gcore_cloud_load_balancer_pool_member" "prod_a" {
  pool_id       = gcore_cloud_load_balancer_pool.prod.id
  project_id    = local.project_id
  region_id     = data.gcore_cloud_region.target.id
  subnet_id     = gcore_cloud_network_subnet.backend.id
  address       = cidrhost(var.backend_subnet_cidr, 10)
  protocol_port = 8080
  weight        = 5
}

resource "gcore_cloud_load_balancer_pool_member" "prod_b" {
  pool_id         = gcore_cloud_load_balancer_pool.prod.id
  project_id      = local.project_id
  region_id       = data.gcore_cloud_region.target.id
  subnet_id       = gcore_cloud_network_subnet.backend.id
  address         = cidrhost(var.backend_subnet_cidr, 11)
  protocol_port   = 8080
  weight          = 3
  monitor_port    = 18080
  monitor_address = cidrhost(var.backend_subnet_cidr, 12)
}

resource "gcore_cloud_load_balancer_pool_member" "canary" {
  pool_id         = gcore_cloud_load_balancer_pool.canary.id
  project_id      = local.project_id
  region_id       = data.gcore_cloud_region.target.id
  subnet_id       = gcore_cloud_network_subnet.backend.id
  address         = cidrhost(var.backend_subnet_cidr, 21)
  protocol_port   = var.canary_backend_port
  weight          = 1
  admin_state_up  = false
}

resource "gcore_cloud_load_balancer_pool_member" "metrics" {
  pool_id         = gcore_cloud_load_balancer_pool.metrics.id
  project_id      = local.project_id
  region_id       = data.gcore_cloud_region.target.id
  subnet_id       = gcore_cloud_network_subnet.backend.id
  address         = cidrhost(var.backend_subnet_cidr, 30)
  protocol_port   = var.prometheus_backend_port
  monitor_port    = var.prometheus_backend_port
  monitor_address = cidrhost(var.backend_subnet_cidr, 31)
}

output "lb_id" {
  value = gcore_cloud_load_balancer.canary.id
}

output "listener_ids" {
  value = {
    prod      = gcore_cloud_load_balancer_listener.http_primary.id
    canary    = gcore_cloud_load_balancer_listener.http_canary.id
    prometheus = gcore_cloud_load_balancer_listener.prometheus.id
  }
}

output "pool_ids" {
  value = {
    prod      = gcore_cloud_load_balancer_pool.prod.id
    canary    = gcore_cloud_load_balancer_pool.canary.id
    prometheus = gcore_cloud_load_balancer_pool.metrics.id
  }
}
