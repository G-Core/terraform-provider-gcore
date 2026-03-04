terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
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
  flavor     = "lb1-2-4"
  name       = "qa-lb-tags-test-fix"

  # ADDING TAGS
  tags = {
    "qa"          = "load-balancer"
    "environment" = "test"
  }
}
