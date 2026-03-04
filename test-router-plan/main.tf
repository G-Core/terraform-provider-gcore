terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {
  # Configuration will be read from .env
}

resource "gcore_cloud_network_router" "router" {
  project_id = 1
  region_id  = 1
  name       = "qa-terr-router"
  external_gateway_info = {
    enable_snat = true
    type        = "default"
  }
}
