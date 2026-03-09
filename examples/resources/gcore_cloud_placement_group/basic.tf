resource "gcore_cloud_placement_group" "servergroup" {
  project_id = 1
  region_id  = 1

  name   = "my-servergroup"
  policy = "affinity"
}
