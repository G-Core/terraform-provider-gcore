resource "gcore_cloud_load_balancer_pool_health_monitor" "tcp" {
  project_id = 1
  region_id  = 1

  pool_id = gcore_cloud_load_balancer_pool.tcp.id

  type             = "TCP"
  delay            = 10
  timeout          = 5
  max_retries      = 3
  max_retries_down = 2
}
