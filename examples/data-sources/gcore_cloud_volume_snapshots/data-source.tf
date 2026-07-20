data "gcore_cloud_volume_snapshots" "example_cloud_volume_snapshots" {
  project_id = 1
  region_id = 1
  instance_id = "550e8400-e29b-41d4-a716-446655440000"
  lifecycle_policy_id = 1
  schedule_id = "67baa7d1-08ea-4fc5-bef2-6b2465b7d227"
  volume_id = "3ed9e2ce-f906-47fb-ba32-c25a3f63df4f"
}
