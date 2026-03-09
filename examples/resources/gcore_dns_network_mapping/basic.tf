# Create a DNS network mapping for private DNS resolution
resource "gcore_dns_network_mapping" "example" {
  name = "DevNetwork"
  mapping = [{
    tags  = ["development", "test"]
    cidr4 = ["10.0.0.0/16", "10.1.0.0/16"]
    cidr6 = ["fd00::/8"]
  }]
}
