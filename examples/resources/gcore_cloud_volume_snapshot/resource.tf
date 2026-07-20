resource "gcore_cloud_volume_snapshot" "example_cloud_volume_snapshot" {
  project_id = 1
  region_id = 1
  name = "my-snapshot"
  volume_id = "67baa7d1-08ea-4fc5-bef2-6b2465b7d227"
  description = "Snapshot description"
  tags = {
    my-tag = "my-tag-value"
  }
}
