resource "gcore_cloud_floating_ip" "private_lb_fip" {
  project_id = 1
  region_id  = 1

  fixed_ip_address = gcore_cloud_load_balancer.private_lb.vip_address
  port_id          = gcore_cloud_load_balancer.private_lb.vip_port_id
}

output "private_lb_fip" {
  value = gcore_cloud_floating_ip.private_lb_fip.floating_ip_address
}
