terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

resource "gcore_cloud_security_group" "sg" {
  project_id = 379987
  region_id  = 76
  name       = "tf-test-desc-removal"
}
