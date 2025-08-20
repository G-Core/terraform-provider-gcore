data "gcore_cloud_volumes" "example_cloud_volumes" {
  project_id = 1
  region_id = 1
  bootable = false
  cluster_id = "t12345"
  has_attachments = true
  id_part = "726ecfcc-7fd0-4e30-a86e-7892524aa483"
  instance_id = "169942e0-9b53-42df-95ef-1a8b6525c2bd"
  name_part = "test"
  tag_key = ["key1", "key2"]
  tag_key_value = "tag_key_value"
}
