resource "gcore_loadbalancerv2" "lb" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  name       = "My first loadbalancer with listener and pool"
  flavor     = "lb1-1-2"
}

