data "gcore_cloud_k8s_cluster_pools" "example_cloud_k8s_cluster_pools" {
  project_id = 1
  region_id = 7
  cluster_name = "my-cluster"
}
