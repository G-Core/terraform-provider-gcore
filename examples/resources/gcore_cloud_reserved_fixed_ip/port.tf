# Reserve an existing port (e.g., from a load balancer VIP)
resource "gcore_cloud_load_balancer" "lb" {
  project_id = 1
  region_id  = 1
  name       = "my-load-balancer"
  flavor     = "lb1-1-2"
}

resource "gcore_cloud_reserved_fixed_ip" "from_port" {
  project_id = 1
  region_id  = 1

  type    = "port"
  port_id = gcore_cloud_load_balancer.lb.vip_port_id
}
