# Restore a volume from a snapshot. The size comes from the snapshot, so it does
# not have to be declared.
resource "gcore_cloud_volume" "restored_volume" {
  project_id  = 1
  region_id   = 1
  name        = "my-restored-volume"
  source      = "snapshot"
  snapshot_id = "88f3e0bd-ca86-4cf7-be8b-dd2988e23c2d"
  type_name   = "standard"
}
