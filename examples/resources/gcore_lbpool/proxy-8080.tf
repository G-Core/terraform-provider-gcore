resource "gcore_lblistener" "proxy_8080" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  loadbalancer_id = gcore_loadbalancerv2.lb.id

  name          = "My first proxy listener with pool"
  protocol      = "TCP"
  protocol_port = 8080
}

resource "gcore_lbpool" "proxy_8080" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  loadbalancer_id = gcore_loadbalancerv2.lb.id
  listener_id     = gcore_lblistener.proxy_8080.id

  name            = "My first proxy pool"
  protocol        = "PROXY"
  lb_algorithm    = "LEAST_CONNECTIONS"
}
