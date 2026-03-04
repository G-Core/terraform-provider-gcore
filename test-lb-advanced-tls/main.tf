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
  name            = "${var.lb_name}-backend"
  cidr            = var.backend_subnet_cidr
  network_id      = gcore_cloud_network.backend.id
  project_id      = local.project_id
  region_id       = data.gcore_cloud_region.target.id
  dns_nameservers = var.backend_dns
}

resource "gcore_cloud_load_balancer" "advanced" {
  project_id            = local.project_id
  region_id             = data.gcore_cloud_region.target.id
  flavor                = var.lb_flavor
  name                  = var.lb_name
  vip_network_id        = gcore_cloud_network.vip.id
  vip_subnet_id         = gcore_cloud_network_subnet.vip.id
  vip_ip_family         = "dual"
  preferred_connectivity = "L2"

  floating_ip {
    source = "new"
  }

  logging {
    enabled               = true
    topic_name            = var.logging_topic_name
    destination_region_id = var.region_id

    retention_policy {
      period = var.logging_retention_days
    }
  }

  tags = var.lb_tags
}

resource "gcore_cloud_load_balancer_listener" "https" {
  project_id       = local.project_id
  region_id        = data.gcore_cloud_region.target.id
  load_balancer_id = gcore_cloud_load_balancer.advanced.id
  name             = "${var.lb_name}-https"
  protocol         = "TERMINATED_HTTPS"
  protocol_port    = 443

  allowed_cidrs        = var.allowed_cidrs
  connection_limit     = var.connection_limit
  insert_x_forwarded   = true
  timeout_client_data  = 30000
  timeout_member_connect = 5000
  timeout_member_data  = 15000

  secret_id     = var.https_secret_id
  sni_secret_id = var.https_sni_secret_ids
}

resource "gcore_cloud_load_balancer_listener" "prometheus" {
  project_id       = local.project_id
  region_id        = data.gcore_cloud_region.target.id
  load_balancer_id = gcore_cloud_load_balancer.advanced.id
  name             = "${var.lb_name}-prometheus"
  protocol         = "PROMETHEUS"
  protocol_port    = var.prometheus_port

  timeout_client_data    = 15000
  timeout_member_connect = 3000
  timeout_member_data    = 5000

  secret_id = var.prometheus_secret_id
  user_list = var.prometheus_users
}

resource "gcore_cloud_load_balancer_pool" "https" {
  project_id       = local.project_id
  region_id        = data.gcore_cloud_region.target.id
  listener_id      = gcore_cloud_load_balancer_listener.https.id
  load_balancer_id = gcore_cloud_load_balancer.advanced.id
  name             = "${var.lb_name}-primary"
  protocol         = "HTTPS"
  lb_algorithm     = "LEAST_CONNECTIONS"

  secret_id     = var.backend_client_secret_id
  ca_secret_id  = var.backend_ca_secret_id

  session_persistence {
    type        = "HTTP_COOKIE"
    cookie_name = "sticky-session"
  }

  healthmonitor {
    type            = "HTTPS"
    delay           = 15
    timeout         = 5
    max_retries     = 3
    max_retries_down = 2
    url_path        = var.healthcheck_path
    http_method     = "GET"
    expected_codes  = "200,302"
  }

  timeout_client_data    = 40000
  timeout_member_connect = 8000
  timeout_member_data    = 20000
}

resource "gcore_cloud_load_balancer_pool_member" "primary_a" {
  pool_id       = gcore_cloud_load_balancer_pool.https.id
  project_id    = local.project_id
  region_id     = data.gcore_cloud_region.target.id
  address       = cidrhost(var.backend_subnet_cidr, 10)
  subnet_id     = gcore_cloud_network_subnet.backend.id
  protocol_port = 8443
  weight        = 10
  monitor_port  = 8443
}

resource "gcore_cloud_load_balancer_pool_member" "primary_b" {
  pool_id         = gcore_cloud_load_balancer_pool.https.id
  project_id      = local.project_id
  region_id       = data.gcore_cloud_region.target.id
  address         = cidrhost(var.backend_subnet_cidr, 20)
  subnet_id       = gcore_cloud_network_subnet.backend.id
  protocol_port   = 8443
  weight          = 5
  backup          = true
  monitor_port    = 9443
  monitor_address = cidrhost(var.backend_subnet_cidr, 21)
}

resource "gcore_cloud_load_balancer_pool" "prometheus" {
  project_id       = local.project_id
  region_id        = data.gcore_cloud_region.target.id
  listener_id      = gcore_cloud_load_balancer_listener.prometheus.id
  load_balancer_id = gcore_cloud_load_balancer.advanced.id
  name             = "${var.lb_name}-metrics"
  protocol         = "HTTP"
  lb_algorithm     = "ROUND_ROBIN"

  timeout_client_data    = 15000
  timeout_member_connect = 3000
  timeout_member_data    = 5000

  healthmonitor {
    type        = "HTTP"
    delay       = 10
    timeout     = 3
    max_retries = 2
    url_path    = "/metrics"
    http_method = "GET"
  }
}

resource "gcore_cloud_load_balancer_pool_member" "metrics" {
  pool_id       = gcore_cloud_load_balancer_pool.prometheus.id
  project_id    = local.project_id
  region_id     = data.gcore_cloud_region.target.id
  address       = cidrhost(var.backend_subnet_cidr, 30)
  subnet_id     = gcore_cloud_network_subnet.backend.id
  protocol_port = var.prometheus_backend_port
  weight        = 1
}

output "lb_id" {
  value = gcore_cloud_load_balancer.advanced.id
}

output "https_listener_id" {
  value = gcore_cloud_load_balancer_listener.https.id
}

output "prometheus_listener_id" {
  value = gcore_cloud_load_balancer_listener.prometheus.id
}

output "https_pool_id" {
  value = gcore_cloud_load_balancer_pool.https.id
}

output "metrics_pool_id" {
  value = gcore_cloud_load_balancer_pool.prometheus.id
}
