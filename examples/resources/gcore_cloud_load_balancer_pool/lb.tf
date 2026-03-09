resource "gcore_cloud_load_balancer" "lb" {
  project_id = 1
  region_id  = 1

  name   = "My first loadbalancer with listener and pool"
  flavor = "lb1-1-2"
}
