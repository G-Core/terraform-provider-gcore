# GPU bare metal cluster with two interfaces: one public and one private
resource "gcore_cloud_gpu_baremetal_cluster" "gpu_cluster" {
  project_id    = 1
  region_id     = 1
  flavor        = "bm3-ai-ndp2-1xlarge-h100-80-8"
  image_id      = "234c133c-b37e-4744-8a26-dc32fe407066"
  name          = "my-gpu-cluster"
  servers_count = 1

  servers_settings = {
    interfaces = [
      {
        type = "external"
      },
      {
        type       = "subnet"
        network_id = gcore_cloud_network.network.id
        subnet_id  = gcore_cloud_network_subnet.subnet.id
      },
    ]
    credentials = {
      ssh_key_name = "my-keypair"
    }
  }
}
