# Attach a floating IP to an existing instance port
resource "gcore_cloud_floating_ip" "web_ip" {
  project_id       = 1
  region_id        = 1
  port_id          = "ee2402d0-f0cd-4503-9b75-69be1d11c5f1"
  fixed_ip_address = "192.168.10.15"

  tags = {
    environment = "production"
    role        = "web-server"
  }
}
