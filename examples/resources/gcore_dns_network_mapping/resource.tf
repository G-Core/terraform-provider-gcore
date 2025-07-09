provider gcore {
  permanent_api_token = "251$d3361.............1b35f26d8"
}

resource "gcore_dns_network_mapping" "example" {
  name = "DevNetwork"
  mapping {
    tags = ["development", "test"]
    cidr4 = ["10.0.0.0/16", "10.1.0.0/16"]
    cidr6 = ["fd00::/8"]
  }
}
