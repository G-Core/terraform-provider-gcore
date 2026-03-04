terraform {
  required_providers {
    gcore = {
      source  = "G-Core/gcore"
      version = "~> 0.10"
    }
  }
}

provider "gcore" {
  # Uses GCORE_PERMANENT_TOKEN env var
}

resource "gcore_fastedge_binary" "test" {
  filename = "${path.module}/../minimal.wasm"
}

output "binary_id" {
  value = gcore_fastedge_binary.test.id
}

output "binary_checksum" {
  value = gcore_fastedge_binary.test.checksum
}
