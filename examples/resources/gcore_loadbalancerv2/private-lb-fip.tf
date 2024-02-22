resource "gcore_floatingip" "private_lb_fip" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  fixed_ip_address = gcore_loadbalancerv2.private_lb.vip_address
  port_id = gcore_loadbalancerv2.private_lb.vip_port_id
}

output "private_lb_fip" {
  value = gcore_floatingip.private_lb_fip.floating_ip_address
}
