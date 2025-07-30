resource "gcore_cloud_instance_image" "example_cloud_instance_image" {
  project_id = 0
  region_id = 0
  image_id = "image_id"
  hw_firmware_type = "bios"
  hw_machine_type = "q35"
  is_baremetal = false
  name = "my-image"
  os_type = "linux"
  ssh_key = "allow"
  tags = {
    foo = "my-tag-value"
  }
}
