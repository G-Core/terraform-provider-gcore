terraform {
  required_providers {
    gcore = { source = "gcore/gcore" }
  }
}

data "gcore_cloud_projects" "my_projects" { name = "default" }
locals { project_id = [for p in data.gcore_cloud_projects.my_projects.items : p.id] }
data "gcore_cloud_region" "rg" { region_id = 76 }

resource "gcore_cloud_network" "test" {
  project_id = local.project_id[0]
  region_id  = data.gcore_cloud_region.rg.id
  name       = "test-net-lb-manual"
}

resource "gcore_cloud_network_subnet" "test" {
  project_id = local.project_id[0]
  region_id  = data.gcore_cloud_region.rg.id
  network_id = gcore_cloud_network.test.id
  name       = "test-subnet-lb-manual"
  cidr       = "10.0.30.0/24"
}

resource "gcore_cloud_load_balancer" "test" {
  project_id     = local.project_id[0]
  region_id      = data.gcore_cloud_region.rg.id
  flavor         = "lb1-2-4"
  name           = var.lb_name
  vip_network_id = gcore_cloud_network.test.id
  vip_subnet_id  = gcore_cloud_network_subnet.test.id
}

variable "lb_name" { default = "test-lb-manual" }

output "lb_id" { value = gcore_cloud_load_balancer.test.id }
output "vrrp_ips" { value = gcore_cloud_load_balancer.test.vrrp_ips }
output "vrrp_ips_count" { value = length(gcore_cloud_load_balancer.test.vrrp_ips) }
