# Test for GCLOUD2-20778: Tags inconsistency error fix
#
# Test Scenario:
# 1. Create LB without tags
# 2. Add tags and verify no inconsistency error
# 3. Verify no drift on subsequent plan

terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

# Get region data
data "gcore_cloud_region" "rg" {
  region_id = var.region_id
}

# Network for LB
resource "gcore_cloud_network" "test" {
  name = "test-net-lb-tags-20778"
}

resource "gcore_cloud_network_subnet" "test" {
  name       = "test-subnet-lb-tags-20778"
  network_id = gcore_cloud_network.test.id
  cidr       = "10.100.0.0/24"
}

# Load Balancer - testing tags behavior
resource "gcore_cloud_load_balancer" "lb" {
  name           = "qa-lb-tags-test-20778"
  flavor         = "lb1-2-4"
  vip_network_id = gcore_cloud_network.test.id
  vip_subnet_id  = gcore_cloud_network_subnet.test.id

  # Tags controlled by variable for testing
  # Start without tags, then add them
  tags = var.lb_tags
}

variable "region_id" {
  default = 76
}

variable "lb_tags" {
  description = "Tags for load balancer - set to {} initially, then add tags"
  type        = map(string)
  default     = {}
}

output "lb_id" {
  value = gcore_cloud_load_balancer.lb.id
}

output "lb_tags" {
  value = gcore_cloud_load_balancer.lb.tags
}

output "lb_tags_v2" {
  value = gcore_cloud_load_balancer.lb.tags_v2
}
