# Basic PostgreSQL cluster example
resource "gcore_postgres_cluster" "basic_cluster" {
  name       = "basic-pg-cluster"
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  flavor {
    cpu    = 1
    memory = 2
  }

  database {
    name  = "mydb"
    owner = "pg"
  }

  network {
    acl = [
      "0.0.0.0/0" # Allow access from anywhere (adjust for production)
    ]
    network_type = "public"
  }

  pg_config {
    version = "15"
  }

  storage {
    size = 20
    type = "ssd-hiiops"
  }

  user {
    name = "pg"
    role_attributes = [
      "LOGIN",
      "CREATEDB",
      "CREATEROLE",
    ]
  }
}
