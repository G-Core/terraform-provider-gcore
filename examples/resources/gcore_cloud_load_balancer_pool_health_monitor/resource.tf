resource "gcore_cloud_load_balancer_pool_health_monitor" "example_cloud_load_balancer_pool_health_monitor" {
  project_id = 1
  region_id = 1
  pool_id = "00000000-0000-4000-8000-000000000000"
  delay = 10
  max_retries = 2
  timeout = 5
  type = "HTTP"
  expected_codes = "200,301,302"
  http_method = "CONNECT"
  max_retries_down = 2
  url_path = "/"
}
