resource "random_password" "prometheus_password" {
  length           = 16
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}

resource "gcore_lblistener" "prometheus_9101" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  loadbalancer_id = gcore_loadbalancerv2.lb.id

  name          = "prometheus-9101"
  protocol      = "PROMETHEUS"
  protocol_port = 9101
  allowed_cidrs = ["10.0.0.0/8"]  # example of how to allow access only from private network

  user_list {
    username = "admin1"
    encrypted_password = random_password.prometheus_password.bcrypt_hash
  }
}

output "prometheus_password" {
  value = random_password.prometheus_password.result
  sensitive = True
}
