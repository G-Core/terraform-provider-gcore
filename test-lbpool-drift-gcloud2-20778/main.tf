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
  flavor     = "lb1-1-2"
  name       = "qa-lbpool-drift-test"
}

resource "gcore_cloud_load_balancer_listener" "ls" {
  project_id       = local.project_id[0]
  region_id        = data.gcore_cloud_region.rg.id
  load_balancer_id = gcore_cloud_load_balancer.lb.id
  name             = "qa-ls-name"
  protocol         = "HTTP"
  protocol_port    = 80
  connection_limit = 5000
  # Using hardcoded bcrypt hash for stable drift testing
  # bcrypt() function generates new salt on each run causing drift - see:
  # https://developer.hashicorp.com/terraform/language/functions/bcrypt
  user_list = [{
    encrypted_password = "$5$isRr.HJ1IrQP38.m$oViu3DJOpUG2ZsjCBtbITV3mqpxxbZfyWJojLPNSPO5"
    username           = "qauser"
  }]
}

# Adding pool after LB and Listener are already created
resource "gcore_cloud_load_balancer_pool" "lb_pool" {
  project_id   = local.project_id[0]
  region_id    = data.gcore_cloud_region.rg.id
  lb_algorithm = "LEAST_CONNECTIONS"
  name         = "pool-drift-test"
  protocol     = "HTTP"
  listener_id  = gcore_cloud_load_balancer_listener.ls.id
}
