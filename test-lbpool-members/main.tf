terraform {
  required_providers {
    gcore = {
      source  = "gcore/gcore"
      version = "~> 1.0"
    }
  }
}

provider "gcore" {
  # Uses environment variables from .env file
  # Will load GCORE_API_KEY automatically
}

# Create a load balancer
resource "gcore_cloud_load_balancer" "test" {
  name       = "test-lb-member-simple"
  flavor     = "lb1-1-2"
  project_id = 379987
  region_id  = 76
}

# Create a pool WITHOUT members
resource "gcore_cloud_load_balancer_pool" "test" {
  name             = "test-pool-simple"
  lb_algorithm     = "ROUND_ROBIN"
  protocol         = "HTTP"
  load_balancer_id = gcore_cloud_load_balancer.test.id
  project_id       = 379987
  region_id        = 76

  # Simple health monitor
  healthmonitor = {
    delay       = 10
    max_retries = 3
    timeout     = 5
    type        = "TCP"
  }
}

# Create a private network and subnet for members
resource "gcore_cloud_network" "test" {
  name       = "test-network-members"
  project_id = 379987
  region_id  = 76
}

resource "gcore_cloud_network_subnet" "test" {
  name       = "test-subnet-members"
  cidr       = "192.168.100.0/24"
  network_id = gcore_cloud_network.test.id
  project_id = 379987
  region_id  = 76

  dns_nameservers = ["8.8.8.8", "8.8.4.4"]
}

# Add first member to the pool using the new member resource
# Using IP addresses from the subnet range
resource "gcore_cloud_load_balancer_pool_member" "member1" {
  pool_id       = gcore_cloud_load_balancer_pool.test.id
  project_id    = 379987
  region_id     = 76
  address       = "192.168.100.10"
  protocol_port = 80
  weight        = 1
  subnet_id     = gcore_cloud_network_subnet.test.id
}

# Add a second member with different weight
resource "gcore_cloud_load_balancer_pool_member" "member2" {
  pool_id       = gcore_cloud_load_balancer_pool.test.id
  project_id    = 379987
  region_id     = 76
  address       = "192.168.100.11"
  protocol_port = 80
  weight        = 2
  subnet_id     = gcore_cloud_network_subnet.test.id
}

output "load_balancer_id" {
  value = gcore_cloud_load_balancer.test.id
}

output "pool_id" {
  value = gcore_cloud_load_balancer_pool.test.id
}

output "network_id" {
  value = gcore_cloud_network.test.id
}

output "subnet_id" {
  value = gcore_cloud_network_subnet.test.id
}

output "member1_id" {
  value = gcore_cloud_load_balancer_pool_member.member1.id
}

output "member1_address" {
  value = gcore_cloud_load_balancer_pool_member.member1.address
}

output "member2_id" {
  value = gcore_cloud_load_balancer_pool_member.member2.id
}

output "member2_weight" {
  value = gcore_cloud_load_balancer_pool_member.member2.weight
}
