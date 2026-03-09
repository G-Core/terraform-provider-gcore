# Create a standalone data volume
resource "gcore_cloud_volume" "data_volume" {
  project_id = 1
  region_id  = 1
  name       = "my-data-volume"
  source     = "new-volume"
  size       = 50
  type_name  = "standard"
}
