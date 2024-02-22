provider gcore {
  permanent_api_token = "251$d3361.............1b35f26d8"
}

data "gcore_project" "project" {
  name = "Default"
}

data "gcore_region" "region" {
  name = "Luxembourg-2"
}

resource "gcore_network" "network" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  name       = "network_example"
  type       = "vxlan"
}
