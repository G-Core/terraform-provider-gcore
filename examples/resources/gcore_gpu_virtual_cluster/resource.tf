provider "gcore" {
  permanent_api_token = "251$d3361.............1b35f26d8"
}

data "gcore_project" "project" {
  name = "Default"
}

data "gcore_region" "region" {
  name = "Luxembourg-2"
}

resource "gcore_gpu_virtual_cluster" "example" {
  name          = "example-gpu-cluster"
  flavor        = "g3-ai-192-1536-12000-l40s-48-8"
  servers_count = 2
  project_id    = data.gcore_project.project.id
  region_id     = data.gcore_region.region.id

  servers_settings {
    interface {
      name      = "pub_net"
      type      = "external"
      ip_family = "ipv4"
    }
    volume {
      name       = "root-volume"
      size       = 120
      type       = "ssd_hiiops"
      source     = "image"
      image_id   = "4536337d-17c7-48f4-8ac5-01a41dc06f58"
      boot_index = 0
    }

    security_groups = []

    credentials {
      ssh_key_name = "my-ssh-key"
    }
  }
} 