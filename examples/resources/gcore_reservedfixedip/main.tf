terraform {
  required_providers {
    gcore = {
      source  = "G-Core/gcore"
      version = ">= 0.3.65"
      # source = "local.gcore.com/repo/gcore"
      # version = ">=0.3.64"
    }
  }
  required_version = ">= 0.13.0"
}

provider gcore {
  gcore_cloud_api = "https://cloud-api-preprod.k8s-ed7-2.cloud.gc.onl/"
  permanent_api_token = "369557$4b9bce05a6857f630c3173f37c34a2ace15e5741cb667f944a4ad8fc72af1a70f2c41a27666c459dc4121a0646bde3a28efb76d6b4ddecfa587c8a4b245a6530"
}

data "gcore_project" "project" {
  name = "Default"
}

data "gcore_region" "region" {
  name = "Luxembourg Preprod"
}
