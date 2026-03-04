terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

# Test DNS Network Mapping Resource
# GCLOUD2-22172 Comprehensive Testing

variable "mapping_name" {
  description = "Name of the network mapping"
  type        = string
  default     = "test-network-mapping-tf-comprehensive"
}

variable "mapping_entries" {
  description = "List of mapping entries"
  type = list(object({
    tags  = list(string)
    cidr4 = optional(list(string), [])
    cidr6 = optional(list(string), [])
  }))
  default = [
    {
      tags  = ["europe", "datacenter-1"]
      cidr4 = ["10.0.0.0/8", "172.16.0.0/12"]
      cidr6 = []
    }
  ]
}

resource "gcore_dns_network_mapping" "test" {
  name = var.mapping_name

  # mapping is a ListNestedAttribute, use attribute syntax
  mapping = [
    for entry in var.mapping_entries : {
      tags  = entry.tags
      cidr4 = length(entry.cidr4) > 0 ? entry.cidr4 : null
      cidr6 = length(entry.cidr6) > 0 ? entry.cidr6 : null
    }
  ]
}

output "mapping_id" {
  value = gcore_dns_network_mapping.test.id
}

output "mapping_name" {
  value = gcore_dns_network_mapping.test.name
}
