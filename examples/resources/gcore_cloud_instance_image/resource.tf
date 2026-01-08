resource "gcore_cloud_instance_image" "example_cloud_instance_image" {
  project_id = 0
  region_id = 0
  name = "my-image"
  url = "http://mirror.noris.net/cirros/0.4.0/cirros-0.4.0-x86_64-disk.img"
  architecture = "x86_64"
  cow_format = false
  hw_firmware_type = "bios"
  hw_machine_type = "q35"
  is_baremetal = false
  os_distro = "ubuntu"
  os_type = "linux"
  os_version = "22.04"
  ssh_key = "allow"
  tags = {
    my-tag = "my-tag-value"
  }
}
