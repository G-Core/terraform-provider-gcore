terraform {
  required_providers {
    gcore = {
      source  = "gcore/gcore"
      version = ">= 0.1"
    }
  }
}

provider "gcore" {
  api_key = "your-api-key"
}
