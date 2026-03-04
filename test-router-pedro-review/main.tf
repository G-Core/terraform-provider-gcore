terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {
  # Credentials from environment: GCORE_API_KEY
}

# Dependencies for router testing
resource "gcore_cloud_network" "test" {
  project_id = 379987
  region_id  = 76
  name       = "pedro-test-network"
}

resource "gcore_cloud_network_subnet" "test" {
  project_id = 379987
  region_id  = 76
  name       = "pedro-test-subnet"
  network_id = gcore_cloud_network.test.id
  cidr       = "192.168.1.0/24"

  # Disable auto-connect - we'll manage router interfaces manually
  connect_to_network_router = false

  # Enable gateway for routes testing
  enable_dhcp    = true
  gateway_ip     = "192.168.1.1"
}

# Main resource under test
resource "gcore_cloud_network_router" "test" {
  project_id = 379987
  region_id  = 76
  name       = "pedro-test-routes-serialization"

  # Workaround for external_gateway_info drift
  external_gateway_info = {
    enable_snat = false
  }

  # Attach subnet so we can add routes
  interfaces = [
    {
      subnet_id = gcore_cloud_network_subnet.test.id
      type      = "subnet"
    }
  ]

  # Testing routes deletion without WithJSONSet
  # Step 2: Partial deletion (3 → 1)
  routes = [
    {
      destination = "10.0.1.0/24"
      nexthop     = "192.168.1.10"
    }
  ]
}

# Outputs for tracking
output "router_id" {
  value = gcore_cloud_network_router.test.id
}

output "router_status" {
  value = gcore_cloud_network_router.test.status
}

output "network_id" {
  value = gcore_cloud_network.test.id
}

output "subnet_id" {
  value = gcore_cloud_network_subnet.test.id
}
