resource "gcore_baremetal" "baremetal_with_two_interface" {
  flavor_id     = "bm1-infrastructure-small"
  name          = "my-baremetal"
  keypair_name  = "my-keypair"
  image_id      = data.gcore_image.ubuntu.id

  interface {
    type = "external"
  }

  interface {
    type = "subnet"

    network_id = gcore_network.network.id
    subnet_id = gcore_subnet.subnet.id
  }
  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
}