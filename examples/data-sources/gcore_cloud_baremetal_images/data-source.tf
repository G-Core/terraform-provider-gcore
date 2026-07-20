data "gcore_cloud_baremetal_images" "example_cloud_baremetal_images" {
  project_id = 1
  region_id = 7
  architecture = "x86_64"
  include_prices = true
  name = "ubuntu"
  os_distro = "ubuntu"
  os_version = "22.04"
  private = "private"
  tag_key = ["key1", "key2"]
  tag_key_value = "tag_key_value"
  visibility = "private"
}
