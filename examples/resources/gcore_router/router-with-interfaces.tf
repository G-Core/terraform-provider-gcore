# Create networks and subnets first
resource "gcore_network" "network_a" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  name          = "network-a"
  create_router = false
}

resource "gcore_subnet" "subnet_a" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  name            = "subnet-a"
  cidr            = "192.168.1.0/24"
  network_id      = gcore_network.network_a.id
  gateway_ip      = "192.168.1.1"
  dns_nameservers = ["8.8.8.8", "8.8.4.4"]
}

resource "gcore_network" "network_b" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  name          = "network-b"
  create_router = false
}

resource "gcore_subnet" "subnet_b" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  name            = "subnet-b"
  cidr            = "192.168.2.0/24"
  network_id      = gcore_network.network_b.id
  gateway_ip      = "192.168.2.1"
  dns_nameservers = ["8.8.8.8", "8.8.4.4"]
}

# Create router connecting both subnets
resource "gcore_router" "multi_subnet" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  name = "multi-subnet-router"

  external_gateway_info {
    type        = "default"
    enable_snat = true
  }

  interfaces {
    type      = "subnet"
    subnet_id = gcore_subnet.subnet_a.id
  }

  interfaces {
    type      = "subnet"
    subnet_id = gcore_subnet.subnet_b.id
  }
}
