# PostgreSQL cluster with high availability
resource "gcore_postgres_cluster" "ha_cluster" {
  name       = "ha-pg-cluster"
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  flavor {
    cpu    = 4
    memory = 8
  }

  database {
    name  = "proddb"
    owner = "produser"
  }

  database {
    name  = "stgdb"
    owner = "stguser"
  }

  network {
    acl = [
      "192.168.1.0/24", # Office network
      "10.0.0.0/8"      # Private network
    ]
    network_type = "public"
  }

  pg_config {
    version     = "15"
    pg_conf     = <<-EOT
      max_connections=300
      shared_buffers=2GB
      effective_cache_size=4GB
      maintenance_work_mem=256MB
      work_mem=8MB
      checkpoint_completion_target=0.9
      wal_buffers=32MB
      min_wal_size=4GB
      max_wal_size=16GB
      random_page_cost=1.1
      effective_io_concurrency=200
      log_statement=all
      log_min_duration_statement=1000
    EOT
    pooler_mode = "transaction"
    pooler_type = "pgbouncer"
  }

  storage {
    size = 100
    type = "ssd-hiiops"
  }

  user {
    name = "produser"
    role_attributes = [
      "LOGIN",
      "CREATEDB"
    ]
  }

  user {
    name = "stguser"
    role_attributes = [
      "LOGIN",
      "CREATEDB"
    ]
  }

  user {
    name = "pg"
    role_attributes = [
      "LOGIN",
      "CREATEDB",
      "CREATEROLE",
    ]
  }

  # Enable synchronous replication for high availability
  ha_replication_mode = "sync"
}

