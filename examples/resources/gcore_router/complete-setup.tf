# Complete example: Private network with internet access via router

# 1. Create a private network
resource "gcore_network" "private" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  name          = "my-private-network"
  create_router = false
}

# 2. Create a subnet in the private network
resource "gcore_subnet" "private" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  name            = "my-private-subnet"
  cidr            = "192.168.100.0/24"
  network_id      = gcore_network.private.id
  gateway_ip      = "192.168.100.1"
  dns_nameservers = ["8.8.8.8", "8.8.4.4"]
}

# 3. Create a router with external gateway for internet access
resource "gcore_router" "internet_gateway" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  name = "internet-gateway-router"

  # Connect to external network for internet access
  external_gateway_info {
    type        = "default"
    enable_snat = true
  }

  # Connect the private subnet to the router
  interfaces {
    type      = "subnet"
    subnet_id = gcore_subnet.private.id
  }
}

# Output the router ID
output "router_id" {
  value = gcore_router.internet_gateway.id
}
