# Create a GPU virtual cluster with public networking
resource "gcore_cloud_gpu_virtual_cluster" "example" {
  project_id    = 1
  region_id     = 1
  name          = "my-gpu-cluster"
  flavor        = "g3-ai-192-1536-12000-l40s-48-8"
  servers_count = 2

  servers_settings = {
    interfaces = [{
      name      = "pub_net"
      type      = "external"
      ip_family = "ipv4"
    }]
    volumes = [{
      name       = "root-volume"
      size       = 120
      type       = "ssd_hiiops"
      source     = "image"
      image_id   = "4536337d-17c7-48f4-8ac5-01a41dc06f58"
      boot_index = 0
    }]
    credentials = {
      ssh_key_name = "my-ssh-key"
    }
  }

  tags = {
    environment = "ml-training"
  }
}
