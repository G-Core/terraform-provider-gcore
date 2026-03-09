# Create a router with external gateway and subnet interface
resource "gcore_cloud_network_router" "main" {
  project_id = 1
  region_id  = 1
  name       = "main-router"

  external_gateway_info = {
    enable_snat = true
    type        = "default"
  }

  interfaces = [{
    subnet_id = gcore_cloud_network_subnet.private.id
    type      = "subnet"
  }]
}
