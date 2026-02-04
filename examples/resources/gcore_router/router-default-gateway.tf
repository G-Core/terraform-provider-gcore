resource "gcore_router" "with_default_gateway" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  name = "router-with-internet"

  external_gateway_info {
    type        = "default"
    enable_snat = true
  }
}
