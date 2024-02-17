provider gcore {
  permanent_api_token = "251$d3361.............1b35f26d8"
}

data "gcore_project" "project" {
  name = "Default"
}

data "gcore_region" "region" {
  name = "Luxembourg-2"
}

data "gcore_lbpool" "pool" {
  region_id  = data.gcore_region.region.id
  project_id = data.gcore_project.project.id

  name       = "test-pool"
}

output "view" {
  value = data.gcore_lbpool.pool
}
