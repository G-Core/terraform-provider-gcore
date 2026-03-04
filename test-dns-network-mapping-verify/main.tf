terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

resource "gcore_dns_network_mapping" "test" {
  name = "tf-test-nm-ds"
  mapping = [
    {
      tags  = ["ds-test"]
      cidr4 = ["10.2.0.0/24"]
    }
  ]
}

data "gcore_dns_network_mapping" "test" {
  id = gcore_dns_network_mapping.test.id
}

output "ds_name" {
  value = data.gcore_dns_network_mapping.test.name
}

output "ds_mapping" {
  value = data.gcore_dns_network_mapping.test.mapping
}
