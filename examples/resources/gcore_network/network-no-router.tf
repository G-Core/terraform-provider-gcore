resource "gcore_network" "isolated" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  name          = "isolated-network"
  create_router = false
}
