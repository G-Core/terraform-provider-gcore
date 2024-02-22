resource "gcore_loadbalancerv2" "lb" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  name       = "My first complex load balancer"
  flavor     = "lb1-1-2"
}

resource "gcore_lblistener" "http_80" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  loadbalancer_id = gcore_loadbalancerv2.lb.id

  name          = "http-80"
  protocol      = "HTTP"
  protocol_port = 80
}

resource "gcore_lbpool" "http" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  loadbalancer_id = gcore_loadbalancerv2.lb.id
  listener_id     = gcore_lblistener.http_80.id

  name            = "My HTTP pool"
  protocol        = "HTTP"
  lb_algorithm    = "ROUND_ROBIN"

  health_monitor {
    type        = "TCP"
    delay       = 10
    max_retries = 3
    timeout     = 5
  }
}
