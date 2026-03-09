resource "gcore_cloud_file_share" "vast" {
  project_id = 1
  region_id  = 1

  name      = "tf-file-share-vast"
  size      = 10
  type_name = "vast"
  protocol  = "NFS"

  share_settings = {
    allowed_characters = "LCD"
    path_length        = "LCD"
    root_squash        = true
  }
}

resource "gcore_cloud_network" "network" {
  project_id = 1
  region_id  = 1

  name = "my-network"
}

resource "gcore_cloud_network_subnet" "subnet" {
  project_id = 1
  region_id  = 1

  name       = "my-subnet"
  cidr       = "192.168.10.0/24"
  network_id = gcore_cloud_network.network.id
}

resource "gcore_cloud_k8s_cluster" "cluster" {
  project_id    = 1
  region_id     = 1
  name          = "my-k8s-cluster"
  fixed_network = gcore_cloud_network.network.id
  fixed_subnet  = gcore_cloud_network_subnet.subnet.id
  keypair       = gcore_cloud_ssh_key.my_keypair.name
  version       = "v1.33.3"

  cni = {
    cloud_k8s_cluster_provider = "cilium"
  }

  csi = {
    nfs = {
      vast_enabled = true
    }
  }

  pools = [{
    name           = "gpu-1"
    flavor_id      = "bm3-ai-ndp-1xlarge-h100-80-8"
    is_public_ipv4 = false
    min_node_count = 1
    max_node_count = 1
  }]
}
