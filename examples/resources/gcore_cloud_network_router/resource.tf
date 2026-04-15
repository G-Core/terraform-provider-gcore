resource "gcore_cloud_network_router" "example_cloud_network_router" {
  project_id = 1
  region_id = 1
  name = "my_wonderful_router"
  external_gateway_info = {
    enable_snat = true
    type = "default"
  }
  interfaces = [{
    subnet_id = "3ed9e2ce-f906-47fb-ba32-c25a3f63df4f"
    type = "subnet"
  }]
  routes = [{
    destination = "10.0.3.0/24"
    nexthop = "10.0.0.13"
  }]
}
