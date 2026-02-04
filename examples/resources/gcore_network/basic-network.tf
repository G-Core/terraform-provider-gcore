resource "gcore_network" "basic" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  name = "my-network"
}
