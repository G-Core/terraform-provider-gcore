terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

resource "gcore_dns_network_mapping" "test" {
  name = "tf-test-nm-get-check"
  mapping = [
    {
      tags  = ["get-test"]
      cidr4 = ["10.88.0.0/24", "10.89.0.0/24"]
    }
  ]
}

output "id" {
  value = gcore_dns_network_mapping.test.id
}

output "name" {
  value = gcore_dns_network_mapping.test.name
}
