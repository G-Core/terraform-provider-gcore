resource "gcore_baremetal" "baremetal_with_one_interface" {
  flavor_id     = "bm1-infrastructure-small"
  name          = "my-baremetal"
  keypair_name  = "my-keypair"
  image_id      = data.gcore_image.ubuntu.id

  interface {
    type = "external"
  }

  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
}