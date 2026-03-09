resource "gcore_cloud_load_balancer" "public_lb" {
  project_id = 1
  region_id  = 1

  name   = "My first public load balancer"
  flavor = "lb1-1-2"

  tags = {
    managed_by = "terraform"
  }
}

output "public_lb_ip" {
  value = gcore_cloud_load_balancer.public_lb.vip_address
}
