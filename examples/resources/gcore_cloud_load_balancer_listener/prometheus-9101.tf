resource "random_password" "prometheus_password" {
  length           = 16
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}

resource "gcore_cloud_load_balancer_listener" "prometheus_9101" {
  project_id = 1
  region_id  = 1

  load_balancer_id = gcore_cloud_load_balancer.lb.id

  name          = "prometheus-9101"
  protocol      = "PROMETHEUS"
  protocol_port = 9101
  allowed_cidrs = ["10.0.0.0/8"] # allow access only from private network

  user_list = [{
    username           = "admin1"
    encrypted_password = random_password.prometheus_password.bcrypt_hash
  }]
}

output "prometheus_password" {
  value     = random_password.prometheus_password.result
  sensitive = true
}
