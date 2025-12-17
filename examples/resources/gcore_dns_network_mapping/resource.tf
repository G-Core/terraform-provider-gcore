resource "gcore_dns_network_mapping" "example_dns_network_mapping" {
  mapping = [{
    cidr4 = ["string"]
    cidr6 = ["string"]
    tags = ["string"]
  }]
  name = "name"
}
