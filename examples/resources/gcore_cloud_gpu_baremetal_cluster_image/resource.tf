resource "gcore_cloud_gpu_baremetal_cluster_image" "example_cloud_gpu_baremetal_cluster_image" {
  project_id = 1
  region_id = 7
  name = "ubuntu-23.10-x64"
  url = "http://mirror.noris.net/cirros/0.4.0/cirros-0.4.0-x86_64-disk.img"
  architecture = "x86_64"
  cow_format = true
  hw_firmware_type = "bios"
  os_distro = "os_distro"
  os_type = "linux"
  os_version = "19.04"
  ssh_key = "allow"
  tags = {
    my-tag = "my-tag-value"
  }
}
