resource "gcore_cloud_volume" "example_cloud_volume" {
  project_id = 1
  region_id = 1
  image_id = "169942e0-9b53-42df-95ef-1a8b6525c2bd"
  name = "volume-1"
  size = 10
  source = "image"
  attachment_tag = "device-tag"
  instance_id_to_attach_to = "88f3e0bd-ca86-4cf7-be8b-dd2988e23c2d"
  lifecycle_policy_ids = [1, 2]
  tags = {

  }
  type_name = "standard"
}
