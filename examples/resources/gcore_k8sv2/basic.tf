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
  version       = "v1.31.9"
  pool {
    name               = "my-k8s-pool"
    flavor_id          = "g1-standard-2-4"
    servergroup_policy = "soft-anti-affinity"
    min_node_count     = 1
    max_node_count     = 1
    boot_volume_size   = 10
    boot_volume_type   = "standard"
  }
}

data "gcore_k8sv2_kubeconfig" "config" {
  cluster_name = gcore_k8sv2.cluster.name
  region_id    = data.gcore_region.region.id
  project_id   = data.gcore_project.project.id
}

// to store kubeconfig in a file pls use
// terraform output -raw kubeconfig > config.yaml
output "kubeconfig" {
  value = data.gcore_k8sv2_kubeconfig.config.kubeconfig
}
