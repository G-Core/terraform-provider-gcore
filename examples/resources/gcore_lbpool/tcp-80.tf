resource "gcore_lblistener" "tcp_80" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  loadbalancer_id = gcore_loadbalancerv2.lb.id

  name          = "My first tcp listener with pool"
  protocol      = "TCP"
  protocol_port = 80
}

resource "gcore_lbpool" "tcp_80" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  loadbalancer_id = gcore_loadbalancerv2.lb.id
  listener_id     = gcore_lblistener.tcp_80.id

  name            = "My first tcp pool"
  protocol        = "TCP"
  lb_algorithm    = "ROUND_ROBIN"

  health_monitor {
    type        = "PING"
    delay       = 10
    max_retries = 5
    timeout     = 5
  }

  session_persistence {
    type        = "APP_COOKIE"
    cookie_name = "test_new_cookie"
  }
}
