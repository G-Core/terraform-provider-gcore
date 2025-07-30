data "gcore_cloud_instance_image" "example_cloud_instance_image" {
  project_id = 0
  region_id = 0
  image_id = "image_id"
  include_prices = true
}
