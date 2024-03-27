resource "gcore_network" "private_network" {
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  name = "my-private-network"
}

resource "gcore_subnet" "private_subnet" {
  count = 2

  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id

  network_id = gcore_network.private_network.id
  name       = "${gcore_network.private_network.name}-subnet-${count.index}"

  cidr       = "172.16.${count.index}.0/24"
}
