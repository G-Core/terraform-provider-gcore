resource "gcore_loadbalancerv2" "public_lb" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  name       = "My first public load balancer"
  flavor     = "lb1-1-2"

  metadata_map = {
    managed_by = "terraform"
  }
}

output "public_lb_ip" {
  value = gcore_loadbalancerv2.public_lb.vip_address
}
