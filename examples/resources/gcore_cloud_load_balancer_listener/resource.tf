resource "gcore_cloud_load_balancer_listener" "example_cloud_load_balancer_listener" {
  project_id = 1
  region_id = 1
  load_balancer_id = "30f4f55b-4a7c-48e0-9954-5cddfee216e7"
  name = "my_listener"
  protocol = "HTTP"
  protocol_port = 80
  allowed_cidrs = ["10.0.0.0/8"]
  connection_limit = 100000
  insert_x_forwarded = false
  secret_id = "f2e734d0-fa2b-42c2-ad33-4c6db5101e00"
  sni_secret_id = ["f2e734d0-fa2b-42c2-ad33-4c6db5101e00", "eb121225-7ded-4ff3-ae1f-599e145dd7cb"]
  timeout_client_data = 50000
  timeout_member_connect = 50000
  timeout_member_data = null
  user_list = [{
    encrypted_password = "$5$isRr.HJ1IrQP38.m$oViu3DJOpUG2ZsjCBtbITV3mqpxxbZfyWJojLPNSPO5"
    username = "admin"
  }]
}
