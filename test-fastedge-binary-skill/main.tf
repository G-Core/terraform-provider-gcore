terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

# Main resource under test
resource "gcore_fastedge_binary" "test" {
  filename = "${path.module}/${var.wasm_file}"
}

# Data source test
data "gcore_fastedge_binary" "test" {
  count = var.enable_datasource ? 1 : 0
  id    = gcore_fastedge_binary.test.id
}

variable "wasm_file" {
  default = "minimal.wasm"
}

variable "enable_datasource" {
  default = false
}

output "binary_id" {
  value = gcore_fastedge_binary.test.id
}

output "binary_checksum" {
  value = gcore_fastedge_binary.test.checksum
}

output "binary_api_type" {
  value = gcore_fastedge_binary.test.api_type
}

output "binary_source" {
  value = gcore_fastedge_binary.test.source
}

output "binary_status" {
  value = gcore_fastedge_binary.test.status
}

output "binary_unref_since" {
  value = gcore_fastedge_binary.test.unref_since
}

output "datasource_checksum" {
  value = var.enable_datasource ? data.gcore_fastedge_binary.test[0].checksum : "N/A"
}

output "datasource_api_type" {
  value = var.enable_datasource ? data.gcore_fastedge_binary.test[0].api_type : "N/A"
}
