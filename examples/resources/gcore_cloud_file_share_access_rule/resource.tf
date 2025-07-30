resource "gcore_cloud_file_share_access_rule" "example_cloud_file_share_access_rule" {
  project_id = 1
  region_id = 1
  file_share_id = "bd8c47ee-e565-4e26-8840-b537e6827b08"
  access_mode = "ro"
  ip_address = "192.168.1.1"
}
