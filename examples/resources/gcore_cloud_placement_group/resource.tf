resource "gcore_cloud_placement_group" "example_cloud_placement_group" {
  project_id = 1
  region_id = 1
  name = "my-server-group"
  policy = "anti-affinity"
}
