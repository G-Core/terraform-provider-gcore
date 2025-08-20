resource "gcore_cloud_reserved_fixed_ip" "example_cloud_reserved_fixed_ip" {
  project_id = 0
  region_id = 0
  type = "external"
  ip_family = "dual"
  is_vip = false
}
