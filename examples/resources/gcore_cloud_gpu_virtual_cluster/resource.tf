resource "gcore_cloud_gpu_virtual_cluster" "example_cloud_gpu_virtual_cluster" {
  project_id = 1
  region_id = 7
  flavor = "g3-ai-32-192-1500-l40s-48-1"
  name = "gpu-cluster-1"
  servers_count = 3
  servers_settings = {
    interfaces = [{
      type = "external"
      ip_family = "ipv4"
      name = "eth0"
    }]
    volumes = [{
      boot_index = 1
      name = "my-data-disk"
      size = 100
      source = "new"
      type = "cold"
      delete_on_termination = true
      tags = {
        key1 = "value1"
      }
    }]
    credentials = {
      password = "securepassword"
      ssh_key_name = "my-ssh-key"
      username = "admin"
    }
    file_shares = [{
      id = "a3f2d1b8-45e6-4f8a-bb5d-19dbf2cd7e9a"
      mount_path = "/mnt/vast"
    }]
    security_groups = [{
      id = "b4849ffa-89f2-45a1-951f-0ae5b7809d98"
    }]
    user_data = "eyJ0ZXN0IjogImRhdGEifQ=="
  }
  tags = {
    my-tag = "my-tag-value"
  }
}
