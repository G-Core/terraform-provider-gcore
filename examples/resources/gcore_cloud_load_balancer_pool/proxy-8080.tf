resource "gcore_cloud_load_balancer_listener" "proxy_8080" {
  project_id = 1
  region_id  = 1

  load_balancer_id = gcore_cloud_load_balancer.lb.id

  name          = "My first proxy listener with pool"
  protocol      = "TCP"
  protocol_port = 8080
}

resource "gcore_cloud_load_balancer_pool" "proxy_8080" {
  project_id = 1
  region_id  = 1

  load_balancer_id = gcore_cloud_load_balancer.lb.id
  listener_id      = gcore_cloud_load_balancer_listener.proxy_8080.id

  name         = "My first proxy pool"
  protocol     = "PROXY"
  lb_algorithm = "LEAST_CONNECTIONS"
}
