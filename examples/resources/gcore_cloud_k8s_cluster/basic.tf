resource "gcore_cloud_network" "network" {
  project_id = 1
  region_id  = 1

  name = "my-network"
}

resource "gcore_cloud_network_subnet" "subnet" {
  project_id = 1
  region_id  = 1

  name       = "my-subnet"
  cidr       = "192.168.10.0/24"
  network_id = gcore_cloud_network.network.id
}

resource "gcore_cloud_k8s_cluster" "cluster" {
  project_id    = 1
  region_id     = 1
  name          = "my-k8s-cluster"
  fixed_network = gcore_cloud_network.network.id
  fixed_subnet  = gcore_cloud_network_subnet.subnet.id
  keypair       = gcore_cloud_ssh_key.my_keypair.name
  version       = "v1.31.9"

  pools = [{
    name               = "my-k8s-pool"
    flavor_id          = "g1-standard-2-4"
    servergroup_policy = "soft-anti-affinity"
    min_node_count     = 1
    max_node_count     = 1
    boot_volume_size   = 10
    boot_volume_type   = "standard"
  }]
}
