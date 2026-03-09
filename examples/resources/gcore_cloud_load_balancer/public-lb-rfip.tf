# After destroying the load balancer, you can attach the Reserved Fixed IP to another load balancer or instance

resource "gcore_cloud_reserved_fixed_ip" "public_lb_fixed_ip" {
  project_id = 1
  region_id  = 1

  is_vip = false
  type   = "external"
}

resource "gcore_cloud_load_balancer" "public_lb_with_fixed_ip" {
  project_id = 1
  region_id  = 1

  name        = "My first public load balancer with reserved fixed ip"
  flavor      = "lb1-1-2"
  vip_port_id = gcore_cloud_reserved_fixed_ip.public_lb_fixed_ip.port_id
}

output "public_lb_with_fixed_ip" {
  value = gcore_cloud_load_balancer.public_lb_with_fixed_ip.vip_address
}
