resource "gcore_cloud_secret" "example_cloud_secret" {
  project_id = 1
  region_id  = 1
  name       = "Load balancer certificate #1"
  payload = {
    certificate_wo       = "<certificate>"
    certificate_chain_wo = "<certificate_chain>"
    private_key_wo       = "<private_key>"
  }
  payload_wo_version = 1
  expiration         = "2019-12-27T18:11:19.117Z"
}
