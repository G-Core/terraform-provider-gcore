terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

data "gcore_cloud_region" "rg" {
  name = "Luxembourg-2"
}

resource "gcore_cloud_network" "test" {
  project_id = 379987
  region_id  = data.gcore_cloud_region.rg.id
  name       = "test-router-manual-net"
}

output "region_id" {
  value = data.gcore_cloud_region.rg.id
}
