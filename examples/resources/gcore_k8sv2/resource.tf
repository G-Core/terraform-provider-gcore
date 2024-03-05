provider gcore {
  permanent_api_token = "251$d3361.............1b35f26d8"
}

resource "gcore_k8sv2" "cl" {
  project_id    = 1
  region_id     = 1
  name          = "cluster1"
  fixed_network = "6bf878c1-1ce4-47c3-a39b-6b5f1d79bf25"
  fixed_subnet  = "dc3a3ea9-86ae-47ad-a8e8-79df0ce04839"
  keypair       = "test_key"
  version       = "v1.26.7"
  pool {
    name             = "pool1"
    flavor_id        = "g1-standard-1-2"
    min_node_count   = 1
    max_node_count   = 1
    boot_volume_size = 10
    boot_volume_type = "standard"
  }
}

data "gcore_k8sv2_kubeconfig" "config" {
  cluster_name       = gcore_k8sv2.cl.name
  region_id          = data.gcore_region.rg.id
  project_id         = data.gcore_project.pr.id
}

// to store kubeconfig in a file pls use
// terraform output -raw kubeconfig > config.yaml
output "kubeconfig" {
  value = data.gcore_k8sv2_kubeconfig.config.kubeconfig
}
