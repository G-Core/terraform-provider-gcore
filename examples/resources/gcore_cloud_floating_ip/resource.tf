resource "gcore_cloud_floating_ip" "example_cloud_floating_ip" {
  project_id = 1
  region_id = 1
  fixed_ip_address = "192.168.10.15"
  port_id = "ee2402d0-f0cd-4503-9b75-69be1d11c5f1"
  tags = {
    my-tag = "my-tag-value"
  }
}
