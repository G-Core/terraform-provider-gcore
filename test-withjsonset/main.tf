terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {
  # Will use GCORE_API_KEY from environment
}

resource "gcore_cloud_network_router" "test" {
  project_id = 379987
  region_id  = 76
  name       = "withjsonset-test"

  # Start with routes - will modify this in test phases
  routes = [
    {
      destination = "10.0.1.0/24"
      nexthop     = "192.168.1.1"
    },
    {
      destination = "10.0.2.0/24"
      nexthop     = "192.168.1.2"
    }
  ]
}

output "router_id" {
  value = gcore_cloud_network_router.test.id
}

output "routes" {
  value = gcore_cloud_network_router.test.routes
}
