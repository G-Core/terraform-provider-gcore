terraform {
  required_providers {
    gcore = {
      source  = "G-Core/gcore"
      version = "~> 0.10"
    }
  }
}

# Old provider uses GCORE_PERMANENT_TOKEN env var for auth
provider "gcore" {
  # permanent_api_token is read from GCORE_PERMANENT_TOKEN env var
}

# Test: Create binary from minimal WASM file
resource "gcore_fastedge_binary" "test" {
  filename = "${path.module}/minimal.wasm"
}

output "binary_id" {
  value = gcore_fastedge_binary.test.id
}

output "binary_checksum" {
  value = gcore_fastedge_binary.test.checksum
}
