resource "gcore_file_share" "vast" {
  name       = "tf-file-share-vast"
  size       = 10
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
  type_name  = "vast"
  protocol   = "NFS"
  share_settings {
    allowed_characters = "LCD"
    path_length = "LCD"
    root_squash = true
  }
}

resource "gcore_network" "network" {
  name       = "my-network"
  type       = "vxlan"
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
}

resource "gcore_subnet" "subnet" {
  name       = "my-subnet"
  cidr       = "192.168.10.0/24"
  network_id = gcore_network.network.id

  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
}

resource "gcore_k8sv2" "cluster" {
  project_id    = data.gcore_project.project.id
  region_id     = data.gcore_region.region.id
  name          = "my-k8s-cluster"
  fixed_network = gcore_network.network.id
  fixed_subnet  = gcore_subnet.subnet.id
  keypair       = gcore_keypair.my_keypair.sshkey_name
  version       = "v1.33.3"

  cni {
    provider = "cilium"
  }

  csi {
    nfs {
      vast_enabled = true
    }
  }

  pool {
    name               = "gpu-1"
    flavor_id          = "bm3-ai-ndp-1xlarge-h100-80-8"
    is_public_ipv4     = false
    min_node_count     = 1
    max_node_count     = 1
  }
}