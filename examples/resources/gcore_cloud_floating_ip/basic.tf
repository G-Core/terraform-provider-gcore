# Create a floating IP address
resource "gcore_cloud_floating_ip" "public_ip" {
  project_id = 1
  region_id  = 1
}
