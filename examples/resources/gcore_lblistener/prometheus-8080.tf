resource "gcore_lblistener" "prometheus_80" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  loadbalancer_id = gcore_loadbalancerv2.lb.id

  name          = "prometheus-80"
  protocol      = "PROMETHEUS"
  protocol_port = 8080
  allowed_cidrs = ["10.0.0.0/8"]  # example of how to allow access only from private network
}
