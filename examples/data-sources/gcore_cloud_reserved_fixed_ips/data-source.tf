data "gcore_cloud_reserved_fixed_ips" "example_cloud_reserved_fixed_ips" {
  project_id = 0
  region_id = 0
  available_only = true
  device_id = "device_id"
  external_only = true
  internal_only = true
  ip_address = "ip_address"
  order_by = "order_by"
  vip_only = true
}
