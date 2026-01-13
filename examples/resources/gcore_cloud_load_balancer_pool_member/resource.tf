resource "gcore_cloud_load_balancer_pool_member" "example_cloud_load_balancer_pool_member" {
  project_id = 1
  region_id = 1
  pool_id = "00000000-0000-4000-8000-000000000000"
  address = "192.168.40.33"
  protocol_port = 80
  admin_state_up = true
  backup = true
  instance_id = "a7e7e8d6-0bf7-4ac9-8170-831b47ee2ba9"
  monitor_address = "monitor_address"
  monitor_port = 1
  subnet_id = "32283b0b-b560-4690-810c-f672cbb2e28d"
  weight = 1
}
