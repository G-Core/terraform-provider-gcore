resource "gcore_cloud_security_group" "example_cloud_security_group" {
  project_id = 0
  region_id = 0
  security_group = {
    name = "my_security_group"
    description = "Some description"
    security_group_rules = [{
      description = "Some description"
      direction = "ingress"
      ethertype = "IPv4"
      port_range_max = 80
      port_range_min = 80
      protocol = "tcp"
      remote_group_id = "00000000-0000-4000-8000-000000000000"
      remote_ip_prefix = "10.0.0.0/8"
    }]
    tags = {
      my-tag = "bar"
    }
  }
  instances = ["00000000-0000-4000-8000-000000000000"]
}
