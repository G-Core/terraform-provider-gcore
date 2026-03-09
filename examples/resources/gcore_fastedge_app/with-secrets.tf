# FastEdge app with secret references and response headers
resource "gcore_fastedge_app" "api_gateway" {
  name   = "api-gateway"
  binary = gcore_fastedge_binary.wasm_module.id
  status = 1

  env = {
    BACKEND_URL = "https://api.example.com"
  }

  # Reference secrets managed by gcore_fastedge_secret
  secret_ids = [101, 102]

  response_headers = {
    X-Powered-By = "Gcore FastEdge"
  }
}
