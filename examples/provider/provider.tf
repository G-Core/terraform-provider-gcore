terraform {
  required_providers {
    gcore = {
      source = "G-Core/gcore"
      # Change the version to the one you want to test
      version = "2.0.0-alpha.1"
    }
  }
}

provider "gcore" {
  api_key = var.gcore_api_key
}
