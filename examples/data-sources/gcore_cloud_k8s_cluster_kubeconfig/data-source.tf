data "gcore_cloud_k8s_cluster_kubeconfig" "example_cloud_k8s_cluster_kubeconfig" {
  project_id = 1
  region_id = 7
  cluster_name = "my-cluster"
}
