data "gcore_cloud_file_shares" "example_cloud_file_shares" {
  project_id = 1
  region_id = 1
  name = "test-sfs"
  type_name = "standard"
}
