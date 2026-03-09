# Create a private network with two subnets for reserved fixed IP examples
resource "gcore_cloud_network" "private_network" {
  project_id = 1
  region_id  = 1

  name = "my-private-network"
}

resource "gcore_cloud_network_subnet" "private_subnet_0" {
  project_id = 1
  region_id  = 1

  network_id = gcore_cloud_network.private_network.id
  name       = "my-private-network-subnet-0"
  cidr       = "172.16.0.0/24"
}

resource "gcore_cloud_network_subnet" "private_subnet_1" {
  project_id = 1
  region_id  = 1

  network_id = gcore_cloud_network.private_network.id
  name       = "my-private-network-subnet-1"
  cidr       = "172.16.1.0/24"
}
