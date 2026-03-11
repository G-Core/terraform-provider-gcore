# Create a private network
resource "gcore_cloud_network" "network" {
  project_id    = 1
  region_id     = 1
  name          = "my-network"
  type          = "vxlan"
  tags = {
    environment = "production"
  }
}

# Create a subnet within the network
resource "gcore_cloud_network_subnet" "subnet" {
  project_id      = 1
  region_id       = 1
  name            = "my-subnet"
  cidr            = "192.168.10.0/24"
  network_id      = gcore_cloud_network.network.id
  dns_nameservers = ["8.8.4.4", "1.1.1.1"]
}
