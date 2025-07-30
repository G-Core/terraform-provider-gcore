resource "gcore_cloud_load_balancer_l7_policy" "example_cloud_load_balancer_l7_policy" {
  project_id = 0
  region_id = 0
  action = "REDIRECT_PREFIX"
  listener_id = "listener_id"
  name = "name"
  position = 0
  redirect_http_code = 0
  redirect_pool_id = "redirect_pool_id"
  redirect_prefix = "redirect_prefix"
  redirect_url = "redirect_url"
  tags = ["string"]
}
