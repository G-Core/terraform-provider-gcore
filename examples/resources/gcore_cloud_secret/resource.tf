resource "gcore_cloud_secret" "example_cloud_secret" {
  project_id = 1
  region_id = 1
  name = "Load balancer certificate #1"
  payload = {
    certificate = "<certificate>"
    certificate_chain = "<certificate_chain>"
    private_key = "<private_key>"
  }
  expiration = "2019-12-27T18:11:19.117Z"
}
