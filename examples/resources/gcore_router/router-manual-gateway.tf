# First, find the external network you want to use
data "gcore_network" "external" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  name = "external-network"
}

resource "gcore_router" "with_manual_gateway" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  name = "router-manual-gateway"

  external_gateway_info {
    type        = "manual"
    enable_snat = true
    network_id  = data.gcore_network.external.id
  }
}
