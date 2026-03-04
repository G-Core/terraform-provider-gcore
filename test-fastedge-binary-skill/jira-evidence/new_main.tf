terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

# Test 1: Create binary from minimal WASM file
resource "gcore_fastedge_binary" "test" {
  filename = "${path.module}/minimal.wasm"
}

output "binary_id" {
  value = gcore_fastedge_binary.test.id
}

output "binary_checksum" {
  value = gcore_fastedge_binary.test.checksum
}

output "binary_status" {
  value = gcore_fastedge_binary.test.status
}

output "binary_api_type" {
  value = gcore_fastedge_binary.test.api_type
}

output "binary_source" {
  value = gcore_fastedge_binary.test.source
}
