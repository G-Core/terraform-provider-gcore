# after destroying load balancer, you can attach the Reserved Fixed IP to another load balancer or instance

resource "gcore_reservedfixedip" "public_lb_fixed_ip" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  is_vip = false
  type   = "external"
}

resource "gcore_loadbalancerv2" "public_lb_with_fixed_ip" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  name       = "My first public load balancer with reserved fixed ip"
  flavor     = "lb1-1-2"
  vip_port_id = gcore_reservedfixedip.public_lb_fixed_ip.port_id
}

output "public_lb_with_fixed_ip" {
  value = gcore_loadbalancerv2.public_lb_with_fixed_ip.vip_address
}
