resource "gcore_ai_cluster" "gpu_cluster" {
  flavor = "bm3-ai-1xlarge-h100-80-8"
  image_id = "37c4fa17-1f18-4904-95f2-dbf39d0318fe"
  cluster_name = "my-gpu-cluster"
  keypair_name = "my-keypair"
  instances_count = 1

  interface {
    type = "external"
  }

  interface {
    type = "subnet"

    network_id = gcore_network.network.id
    subnet_id = gcore_subnet.subnet.id
  }

  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
}
