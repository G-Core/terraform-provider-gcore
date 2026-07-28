data "gcore_cloud_k8s_cluster_pool" "example_cloud_k8s_cluster_pool" {
  project_id = 1
  region_id = 7
  cluster_name = "my-cluster"
  pool_name = "my-pool"
}
