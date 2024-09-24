resource "gcore_baremetal" "windows_baremetal" {
  flavor_id     = "bm1-infrastructure-small"
  name          = "my-windows-baremetal"
  password      = "my-s3cR3tP@ssw0rd"
  image_id      = data.gcore_image.windows.id

  interface {
    type = "external"
  }

  project_id = data.gcore_project.project.id
  region_id  = data.gcore_region.region.id
}