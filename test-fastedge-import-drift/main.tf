terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

resource "gcore_fastedge_binary" "test" {
  filename = "${path.module}/test.wasm"
}

output "binary_id" {
  value = gcore_fastedge_binary.test.id
}

output "binary_checksum" {
  value = gcore_fastedge_binary.test.checksum
}
