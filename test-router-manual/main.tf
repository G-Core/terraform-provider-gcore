terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {
  # Credentials loaded from environment variables:
  # GCORE_CLOUD_API_KEY
  # Or from .env via set_env.sh
}

# Network
resource "gcore_cloud_network" "nw" {
  project_id    = 379987
  region_id     = 76
  create_router = true
  name          = "qa-terr-nw"
}

# Subnet
resource "gcore_cloud_network_subnet" "sb" {
  project_id = 379987
  region_id  = 76
  name       = "sys"
  cidr       = "192.168.0.0/24"
  network_id = gcore_cloud_network.nw.id
}

# Router with interface attached
resource "gcore_cloud_network_router" "router" {
  project_id = 379987
  region_id  = 76
  name       = "qa-terr-router"

  external_gateway_info = {
    enable_snat = true
    //type        = "default"
  }

  interfaces = [
    {
      subnet_id = gcore_cloud_network_subnet.sb.id
      type      = "subnet"
    }
  ]
}

# Outputs
output "network_id" {
  value       = gcore_cloud_network.nw.id
  description = "Network ID"
}

output "subnet_id" {
  value       = gcore_cloud_network_subnet.sb.id
  description = "Subnet ID"
}

output "router_id" {
  value       = gcore_cloud_network_router.router.id
  description = "Router ID"
}

output "router_interfaces" {
  value       = gcore_cloud_network_router.router.interfaces
  description = "Router interfaces"
}
