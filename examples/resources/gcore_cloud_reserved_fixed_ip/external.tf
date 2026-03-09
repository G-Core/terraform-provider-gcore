# Reserve an external (public) IP address
resource "gcore_cloud_reserved_fixed_ip" "external" {
  project_id = 1
  region_id  = 1

  type   = "external"
  is_vip = false
}
