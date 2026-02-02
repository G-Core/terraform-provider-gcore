data "gcore_cloud_k8s_cluster" "example_cloud_k8s_cluster" {
  project_id = 1
  region_id = 7
  cluster_name = "my-cluster"
}
