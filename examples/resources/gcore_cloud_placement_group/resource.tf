resource "gcore_cloud_placement_group" "example_cloud_placement_group" {
  project_id = 0
  region_id = 0
  name = "my-server-group"
  policy = "anti-affinity"
}
