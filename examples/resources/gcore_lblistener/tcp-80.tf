resource "gcore_lblistener" "tcp_80" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  loadbalancer_id = gcore_loadbalancerv2.lb.id

  name          = "tcp-80"
  protocol      = "TCP"
  protocol_port = 80
}