data "gcore_cloud_gpu_virtual_clusters" "example_cloud_gpu_virtual_clusters" {
  project_id = 1
  region_id = 7
  created_at = {
    gt = "2019-12-27T18:11:19.117Z"
    gte = "2019-12-27T18:11:19.117Z"
    lt = "2019-12-27T18:11:19.117Z"
    lte = "2019-12-27T18:11:19.117Z"
  }
  flavor = {
    contains = ["string"]
    exact = ["string"]
    prefix = ["string"]
    suffix = ["string"]
  }
  ids = ["1aaaab48-10d0-46d9-80cc-85209284ceb4"]
  name = {
    contains = ["string"]
    exact = ["string"]
    prefix = ["string"]
    suffix = ["string"]
  }
  servers_count = {
    gt = 0
    gte = 0
    lt = 0
    lte = 0
  }
  tag_key = {
    contains = ["string"]
    exact = ["string"]
    prefix = ["string"]
    suffix = ["string"]
  }
  tag_value = {
    contains = ["string"]
    exact = ["string"]
    prefix = ["string"]
    suffix = ["string"]
  }
  tags = {
    env = "prod"
  }
  updated_at = {
    gt = "2019-12-27T18:11:19.117Z"
    gte = "2019-12-27T18:11:19.117Z"
    lt = "2019-12-27T18:11:19.117Z"
    lte = "2019-12-27T18:11:19.117Z"
  }
}
