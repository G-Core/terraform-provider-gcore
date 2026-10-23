resource "gcore_cloud_load_balancer" "lb" {
  project_id = 1
  region_id  = 1

  name   = "My first complex load balancer"
  flavor = "lb1-1-2"
}

resource "gcore_cloud_load_balancer_listener" "http_80" {
  project_id = 1
  region_id  = 1

  load_balancer_id = gcore_cloud_load_balancer.lb.id

  name          = "http-80"
  protocol      = "HTTP"
  protocol_port = 80
}

resource "gcore_cloud_load_balancer_pool" "http" {
  project_id = 1
  region_id  = 1

  load_balancer_id = gcore_cloud_load_balancer.lb.id
  listener_id      = gcore_cloud_load_balancer_listener.http_80.id

  name         = "My HTTP pool"
  protocol     = "HTTP"
  lb_algorithm = "ROUND_ROBIN"
}

resource "gcore_cloud_load_balancer_pool" "tcp" {
  project_id = 1
  region_id  = 1

  load_balancer_id = gcore_cloud_load_balancer.lb.id

  name         = "My TCP pool"
  protocol     = "TCP"
  lb_algorithm = "LEAST_CONNECTIONS"
}
