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

# Step 1: Create Load Balancer and Listener with user_list
resource "gcore_cloud_load_balancer" "lb" {
  project_id = local.project_id[0]
  region_id  = data.gcore_cloud_region.rg.id
  flavor     = "lb1-1-2"
  name       = "qa-drift-test-lb-tags"
}

resource "gcore_cloud_load_balancer_listener" "ls" {
  project_id       = local.project_id[0]
  region_id        = data.gcore_cloud_region.rg.id
  load_balancer_id = gcore_cloud_load_balancer.lb.id
  name             = "qa-ls-name"
  protocol         = "HTTP"
  protocol_port    = 80

  # This user_list with encrypted_password is key to reproducing the bug
  user_list = [
    {
      username           = "testuser"
      encrypted_password = "$2a$10$vp25Soo.i6aYcyWtrfV7SeBsXMa1GjMRHZRJyMCTAiY/T6j8kXv7u"
    }
  ]
}

# Step 2: This pool will be added in the second apply
# Uncomment this block after first apply to reproduce the drift
# resource "gcore_cloud_load_balancer_pool" "lb_pool" {
#   project_id   = local.project_id[0]
#   region_id    = data.gcore_cloud_region.rg.id
#   lb_algorithm = "LEAST_CONNECTIONS"
#   name         = "pool-drift-test"
#   protocol     = "HTTP"
#   listener_id  = gcore_cloud_load_balancer_listener.ls.id
# }
