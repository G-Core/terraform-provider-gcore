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

resource "gcore_cloud_load_balancer" "udp" {
  project_id             = local.project_id
  region_id              = data.gcore_cloud_region.target.id
  flavor                 = var.lb_flavor
  name                   = var.lb_name
  vip_network_id         = gcore_cloud_network.vip.id
  vip_subnet_id          = gcore_cloud_network_subnet.vip.id
  preferred_connectivity = "L3"

  tags = var.lb_tags
}

resource "gcore_cloud_load_balancer_listener" "udp" {
  project_id       = local.project_id
  region_id        = data.gcore_cloud_region.target.id
  load_balancer_id = gcore_cloud_load_balancer.udp.id
  name             = "${var.lb_name}-udp"
  protocol         = "UDP"
  protocol_port    = var.udp_listener_port

  timeout_client_data    = var.timeout_client_data
  timeout_member_connect = var.timeout_member_connect
  timeout_member_data    = var.timeout_member_data
}

resource "gcore_cloud_load_balancer_pool" "udp" {
  project_id       = local.project_id
  region_id        = data.gcore_cloud_region.target.id
  listener_id      = gcore_cloud_load_balancer_listener.udp.id
  load_balancer_id = gcore_cloud_load_balancer.udp.id
  name             = "${var.lb_name}-pool"
  protocol         = "UDP"
  lb_algorithm     = "SOURCE_IP"

  session_persistence = {
    type                = "SOURCE_IP"
    persistence_timeout = var.persistence_timeout
  }

  timeout_client_data    = var.timeout_client_data
  timeout_member_connect = var.timeout_member_connect
  timeout_member_data    = var.timeout_member_data

  healthmonitor = {
    type        = "UDP-CONNECT"
    delay       = 10
    timeout     = 5
    max_retries = 3
  }

  members = [
    {
      address        = cidrhost(var.backend_subnet_cidr, 5)
      protocol_port  = var.backend_port
      subnet_id      = gcore_cloud_network_subnet.backend.id
      weight         = 4
      monitor_port   = var.backend_port
      monitor_address = cidrhost(var.backend_subnet_cidr, 6)
    },
    {
      address        = cidrhost(var.backend_subnet_cidr, 15)
      protocol_port  = var.backend_port
      subnet_id      = gcore_cloud_network_subnet.backend.id
      weight         = 2
      backup         = true
      monitor_port   = var.backend_port
      monitor_address = cidrhost(var.backend_subnet_cidr, 16)
    }
  ]
}

output "udp_lb_id" {
  value = gcore_cloud_load_balancer.udp.id
}

output "udp_listener_id" {
  value = gcore_cloud_load_balancer_listener.udp.id
}

output "udp_pool_id" {
  value = gcore_cloud_load_balancer_pool.udp.id
}

output "udp_member_addresses" {
  value = [for m in gcore_cloud_load_balancer_pool.udp.members : m.address]
}
