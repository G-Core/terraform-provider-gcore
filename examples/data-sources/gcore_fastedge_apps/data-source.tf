data "gcore_fastedge_apps" "example_fastedge_apps" {
  api_type = "wasi-http"
  binary = 1
  limit = 1
  name = "x"
  ordering = "name"
  plan = 1
  status = 0
  template = 1
}
