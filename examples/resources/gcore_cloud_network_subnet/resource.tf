resource "gcore_cloud_network_subnet" "example_cloud_network_subnet" {
  project_id = 1
  region_id = 1
  cidr = "192.168.10.0/24"
  name = "my subnet"
  network_id = "ee2402d0-f0cd-4503-9b75-69be1d11c5f1"
  connect_to_network_router = true
  dns_nameservers = ["8.8.4.4", "1.1.1.1"]
  enable_dhcp = true
  gateway_ip = "192.168.10.1"
  host_routes = [{
    destination = "10.0.3.0/24"
    nexthop = "10.0.0.13"
  }]
  ip_version = 4
  router_id_to_connect = "00000000-0000-4000-8000-000000000000"
  tags = {
    my-tag = "my-tag-value"
  }
}
