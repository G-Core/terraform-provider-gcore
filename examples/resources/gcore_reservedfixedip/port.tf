resource "gcore_loadbalancerv2" "lb" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
  name       = "My first public load balancer"
  flavor     = "lb1-1-2"
}

resource "gcore_reservedfixedip" "fixed_ip_by_port" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  type       = "port"
  port_id = gcore_loadbalancerv2.lb.vip_port_id
}
