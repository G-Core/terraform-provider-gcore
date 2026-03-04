terraform {
  required_providers {
    gcore = {
      source = "gcore/gcore"
    }
  }
}

provider "gcore" {}

# Use existing subnet (data source)
data "gcore_cloud_network_subnet" "existing" {
  project_id = 379987
  region_id  = 76
  subnet_id  = "59a5f550-7fac-4b02-a834-3385b48cc79b"
}

resource "gcore_cloud_network_router" "router" {
  project_id = 379987
  region_id  = 76
  name       = "mitm-test-router-subnet"

  external_gateway_info = {
    enable_snat = true
    type        = "default"
  }

  # Remove interface - should trigger POST /detach
  interfaces = []
}

output "router_id" {
  value = gcore_cloud_network_router.router.id
}

output "subnet_id" {
  value = data.gcore_cloud_network_subnet.existing.id
}
