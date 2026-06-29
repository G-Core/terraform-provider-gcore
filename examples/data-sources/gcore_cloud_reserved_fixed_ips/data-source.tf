data "gcore_cloud_reserved_fixed_ips" "example_cloud_reserved_fixed_ips" {
  project_id = 1
  region_id = 4
  device_id = "device_id"
  ip_address = "ip_address"
}
