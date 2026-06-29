resource "gcore_cloud_reserved_fixed_ip" "example_cloud_reserved_fixed_ip" {
  project_id = 1
  region_id = 4
  type = "external"
  ip_family = "ipv4"
  is_vip = false
}
