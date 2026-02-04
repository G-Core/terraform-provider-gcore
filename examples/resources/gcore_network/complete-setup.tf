# Complete example: Network with subnet

# 1. Create the network
resource "gcore_network" "main" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  name = "my-application-network"

  metadata_map = {
    environment = "production"
    managed_by  = "terraform"
  }
}

# 2. Create a subnet in the network
resource "gcore_subnet" "main" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  name            = "my-application-subnet"
  cidr            = "192.168.0.0/24"
  network_id      = gcore_network.main.id
  gateway_ip      = "192.168.0.1"
  dns_nameservers = ["8.8.8.8", "8.8.4.4"]
}

# Outputs
output "network_id" {
  value = gcore_network.main.id
}

output "subnet_id" {
  value = gcore_subnet.main.id
}
