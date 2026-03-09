# Create a FastEdge app from a Wasm binary
resource "gcore_fastedge_app" "example" {
  name    = "my-edge-app"
  comment = "My FastEdge application"
  binary  = gcore_fastedge_binary.wasm_module.id
  status  = 1

  env = {
    API_URL = "https://api.example.com"
    DEBUG   = "false"
  }
}
