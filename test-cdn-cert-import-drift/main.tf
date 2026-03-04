terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

resource "gcore_cdn_certificate" "test" {
  name      = "tf-test-import-drift-check"
  automated = true
}
