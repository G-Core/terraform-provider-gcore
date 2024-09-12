resource "gcore_servergroup" "servergroup" {
  name       = "my-servergroup"
  policy     = "affinity"
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
}
