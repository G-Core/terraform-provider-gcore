terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

resource "gcore_cdn_client_config" "test" {
  utilization_level = 90
}

data "gcore_cdn_client_config" "info" {
  depends_on = [gcore_cdn_client_config.test]
}

output "resource" {
  value = gcore_cdn_client_config.test
}

output "data_source" {
  value = data.gcore_cdn_client_config.info
}
