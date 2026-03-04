terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

resource "gcore_dns_zone" "test" {
  name           = "tf-test-set-records.qa"
  dnssec_enabled = true
  meta = {
    webhook        = "https://example.com/hook-updated"
    webhook_method = "PUT"
  }
}
