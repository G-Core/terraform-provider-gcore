terraform {
  required_version = ">= 1.5"

  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {
  # Credentials should be set via environment variables or .terraformrc
}
