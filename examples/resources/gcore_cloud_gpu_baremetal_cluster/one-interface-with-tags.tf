# GPU bare metal cluster with one public interface and tags
resource "gcore_cloud_gpu_baremetal_cluster" "gpu_cluster" {
  project_id    = 1
  region_id     = 1
  flavor        = "bm3-ai-ndp2-1xlarge-h100-80-8"
  image_id      = "234c133c-b37e-4744-8a26-dc32fe407066"
  name          = "my-gpu-cluster"
  servers_count = 1

  servers_settings = {
    interfaces = [{
      type = "external"
    }]
    credentials = {
      ssh_key_name = "my-keypair"
    }
  }

  tags = {
    my-tag-key = "my-tag-value"
  }
}
