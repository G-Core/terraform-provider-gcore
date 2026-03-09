resource "gcore_cloud_load_balancer_listener" "tcp_80" {
  project_id = 1
  region_id  = 1

  load_balancer_id = gcore_cloud_load_balancer.lb.id

  name          = "My first tcp listener with pool"
  protocol      = "TCP"
  protocol_port = 80
}

resource "gcore_cloud_load_balancer_pool" "tcp_80" {
  project_id = 1
  region_id  = 1

  load_balancer_id = gcore_cloud_load_balancer.lb.id
  listener_id      = gcore_cloud_load_balancer_listener.tcp_80.id

  name         = "My first tcp pool"
  protocol     = "TCP"
  lb_algorithm = "ROUND_ROBIN"

  healthmonitor = {
    type        = "PING"
    delay       = 10
    max_retries = 5
    timeout     = 5
  }

  session_persistence = {
    type        = "APP_COOKIE"
    cookie_name = "test_new_cookie"
  }
}
