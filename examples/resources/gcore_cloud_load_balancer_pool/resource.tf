resource "gcore_cloud_load_balancer_pool" "example_cloud_load_balancer_pool" {
  project_id = 1
  region_id = 1
  lb_algorithm = "LEAST_CONNECTIONS"
  name = "pool_name"
  protocol = "HTTP"
  ca_secret_id = "ca_secret_id"
  crl_secret_id = "crl_secret_id"
  healthmonitor = {
    delay = 10
    max_retries = 3
    timeout = 5
    type = "HTTP"
    expected_codes = "200,301,302"
    http_method = "GET"
    max_retries_down = 3
    url_path = "/"
  }
  listener_id = "listener_id"
  loadbalancer_id = "bbb35f84-35cc-4b2f-84c2-a6a29bba68aa"
  members = [{
    address = "192.168.1.101"
    protocol_port = 8000
    admin_state_up = true
    backup = true
    instance_id = "a7e7e8d6-0bf7-4ac9-8170-831b47ee2ba9"
    monitor_address = "monitor_address"
    monitor_port = 0
    subnet_id = "32283b0b-b560-4690-810c-f672cbb2e28d"
    weight = 2
  }, {
    address = "192.168.1.102"
    protocol_port = 8000
    admin_state_up = true
    backup = true
    instance_id = "169942e0-9b53-42df-95ef-1a8b6525c2bd"
    monitor_address = "monitor_address"
    monitor_port = 0
    subnet_id = "32283b0b-b560-4690-810c-f672cbb2e28d"
    weight = 4
  }]
  secret_id = "secret_id"
  session_persistence = {
    type = "APP_COOKIE"
    cookie_name = "cookie_name"
    persistence_granularity = "persistence_granularity"
    persistence_timeout = 0
  }
  timeout_client_data = 50000
  timeout_member_connect = 50000
  timeout_member_data = 0
}
