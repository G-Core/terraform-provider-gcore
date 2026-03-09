resource "gcore_cloud_load_balancer_pool_member" "public_member" {
  project_id = 1
  region_id  = 1

  pool_id = gcore_cloud_load_balancer_pool.http.id

  address       = "8.8.8.8"
  protocol_port = 80
  weight        = 1
}
