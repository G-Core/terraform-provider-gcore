resource "gcore_network" "with_metadata" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  name = "my-tagged-network"

  metadata_map = {
    environment = "production"
    team        = "platform"
    managed_by  = "terraform"
  }
}
