resource "gcore_cloud_load_balancer_listener" "tcp_80" {
  project_id = 1
  region_id  = 1

  load_balancer_id = gcore_cloud_load_balancer.lb.id

  name          = "tcp-80"
  protocol      = "TCP"
  protocol_port = 80
}
