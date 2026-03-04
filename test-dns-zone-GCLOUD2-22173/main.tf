terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

# Step 2: Toggle only dnssec_enabled = true (should send PATCH + GET only, no PUT)
resource "gcore_dns_zone" "test" {
  name           = "dnssec-skip-put-test.qa"
  dnssec_enabled = true
}
