resource "gcore_cloud_load_balancer_pool_health_monitor" "http" {
  project_id = 1
  region_id  = 1

  pool_id = gcore_cloud_load_balancer_pool.http.id

  type        = "HTTP"
  delay       = 10
  timeout     = 5
  max_retries = 3

  http_method    = "GET"
  url_path       = "/health"
  expected_codes = "200,204"
}
