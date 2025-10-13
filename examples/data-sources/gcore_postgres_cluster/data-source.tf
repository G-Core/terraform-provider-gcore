provider gcore {
  permanent_api_token = "251$d3361.............1b35f26d8"
}

data "gcore_project" "project" {
  name = "Default"
}

data "gcore_region" "region" {
  name = "Luxembourg Preprod"
}

data "gcore_postgres_cluster" "postgres_cluster" {
  name       = "poliveira-db"
  region_id  = data.gcore_region.region.id
  project_id = data.gcore_project.project.id
}

output "postgres_cluster" {
  value = data.gcore_postgres_cluster.postgres_cluster
}

