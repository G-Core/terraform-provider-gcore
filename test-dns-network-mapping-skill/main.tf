terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

resource "gcore_dns_network_mapping" "test" {
  name    = var.name
  mapping = var.mapping
}

data "gcore_dns_network_mapping" "test" {
  count = var.test_data_source ? 1 : 0
  id    = gcore_dns_network_mapping.test.id
}

variable "name" {
  default = "tf-test-nm-skill"
}

variable "mapping" {
  default = [
    {
      tags  = ["dev", "dc1"]
      cidr4 = ["10.0.0.0/24"]
      cidr6 = ["2001:db8:1::/48"]
    },
    {
      tags  = ["prod", "dc2"]
      cidr4 = ["10.1.0.0/24"]
    }
  ]
}

variable "test_data_source" {
  default = false
}

output "resource_id" {
  value = gcore_dns_network_mapping.test.id
}

output "resource_name" {
  value = gcore_dns_network_mapping.test.name
}

output "resource_mapping" {
  value = gcore_dns_network_mapping.test.mapping
}

output "ds_name" {
  value = var.test_data_source ? data.gcore_dns_network_mapping.test[0].name : null
}

output "ds_mapping" {
  value = var.test_data_source ? data.gcore_dns_network_mapping.test[0].mapping : null
}
