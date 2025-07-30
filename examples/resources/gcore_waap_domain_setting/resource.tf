resource "gcore_waap_domain_setting" "example_waap_domain_setting" {
  domain_id = 1
  api = {
    api_urls = ["api/v1/.*", "v2/.*"]
    is_api = true
  }
  ddos = {
    burst_threshold = 30
    global_threshold = 250
  }
}
