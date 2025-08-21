resource "gcore_k8sv2" "cluster" {
  project_id    = data.gcore_project.project.id
  region_id     = data.gcore_region.region.id
  name          = "my-k8s-cluster"
  keypair       = gcore_keypair.my_keypair.sshkey_name
  version       = "v1.31.9"
  pool {
    name             = "my-k8s-pool"
    flavor_id        = "g1-standard-2-4"
    servergroup_policy = "soft-anti-affinity"
    min_node_count   = 1
    max_node_count   = 1
    boot_volume_size = 10
    boot_volume_type = "standard"
    is_public_ipv4 = true
  }
  security_group_rules {
    direction      = "ingress"
    ethertype      = "IPv4"
    protocol       = "tcp"
    port_range_min = 80
    port_range_max = 80
  }
}

data "gcore_k8sv2_kubeconfig" "config" {
  cluster_name       = gcore_k8sv2.cluster.name
  region_id          = data.gcore_region.region.id
  project_id         = data.gcore_project.project.id
}

// to store kubeconfig in a file pls use
// terraform output -raw kubeconfig > config.yaml
output "kubeconfig" {
  value = data.gcore_k8sv2_kubeconfig.config.kubeconfig
}
